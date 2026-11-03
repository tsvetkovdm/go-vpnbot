package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/tsvetkovdm/go-vpnbot/internal/core/config"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_pgx_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool/pgx"
	core_http_server "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/http/server"
	core_tg_client "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/client"
	core_tg_dialog "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/dialog"
	core_tg_middleware "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/middleware"
	core_tg_poller "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/poller"
	core_tg_screen "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/screen"
	access_postgres_repository "github.com/tsvetkovdm/go-vpnbot/internal/features/access/repository/postgres"
	access_service "github.com/tsvetkovdm/go-vpnbot/internal/features/access/service"
	access_transport_http "github.com/tsvetkovdm/go-vpnbot/internal/features/access/transport/http"
	admins_postgres_repository "github.com/tsvetkovdm/go-vpnbot/internal/features/admins/repository/postgres"
	admins_service "github.com/tsvetkovdm/go-vpnbot/internal/features/admins/service"
	admins_transport_tg "github.com/tsvetkovdm/go-vpnbot/internal/features/admins/transport/tg"
	payments_postgres_repository "github.com/tsvetkovdm/go-vpnbot/internal/features/payments/repository/postgres"
	payments_service "github.com/tsvetkovdm/go-vpnbot/internal/features/payments/service"
	payments_transport_tg "github.com/tsvetkovdm/go-vpnbot/internal/features/payments/transport/tg"
	plans_postgres_repository "github.com/tsvetkovdm/go-vpnbot/internal/features/plans/repository/postgres"
	plans_service "github.com/tsvetkovdm/go-vpnbot/internal/features/plans/service"
	plans_transport_tg "github.com/tsvetkovdm/go-vpnbot/internal/features/plans/transport/tg"
	subscriptions_postgres_repository "github.com/tsvetkovdm/go-vpnbot/internal/features/subscriptions/repository/postgres"
	subscriptions_service "github.com/tsvetkovdm/go-vpnbot/internal/features/subscriptions/service"
	subscriptions_transport_tg "github.com/tsvetkovdm/go-vpnbot/internal/features/subscriptions/transport/tg"
	users_postgres_repository "github.com/tsvetkovdm/go-vpnbot/internal/features/users/repository/postgres"
	users_service "github.com/tsvetkovdm/go-vpnbot/internal/features/users/service"
	users_transport_tg "github.com/tsvetkovdm/go-vpnbot/internal/features/users/transport/tg"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone", zap.Any("zone", time.Local))

	bot, err := core_tg_client.NewBotAPI(core_tg_client.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init bot API:", zap.Error(err))
	}

	screens := core_tg_screen.NewScreens(bot)

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	dialogStore := core_tg_dialog.NewStore(15 * time.Minute)

	logger.Debug("initializing feature", zap.String("feature", "admins"))
	adminsRepository := admins_postgres_repository.NewAdminsRepository(pool)
	adminsService := admins_service.NewAdminsService(adminsRepository)
	adminOnly := core_tg_middleware.AdminOnly(bot, adminsService)
	superAdminOnly := core_tg_middleware.SuperAdminOnly(bot, adminsService)
	adminsTransportTG := admins_transport_tg.NewAdminsTGHandler(
		bot,
		adminsService,
		adminOnly,
		superAdminOnly,
		screens,
	)
	if err := adminsService.EnsureSuperAdmin(ctx, cfg.SuperAdminTelegramID); err != nil {
		logger.Fatal("failed to ensure super admin", zap.Error(err))
	}

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportTG := users_transport_tg.NewUsersTGHandler(
		bot,
		usersService,
		adminOnly,
		dialogStore,
		screens,
	)

	logger.Debug("initializing feature", zap.String("feature", "subscriptions"))
	subscriptionsRepository := subscriptions_postgres_repository.NewSubscriptionsRepository(pool)
	subscriptionsService := subscriptions_service.NewSubscriptionsService(subscriptionsRepository)
	subscriptionsTransportTG := subscriptions_transport_tg.NewSubscriptionsTGHandler(
		bot,
		subscriptionsService,
		adminOnly,
		dialogStore,
	)

	logger.Debug("initializing feature", zap.String("feature", "plans"))
	plansRepository := plans_postgres_repository.NewPlansRepository(pool)
	plansService := plans_service.NewPlansService(plansRepository)
	plansTransportTG := plans_transport_tg.NewPlansTGHandler(
		bot,
		plansService,
		adminOnly,
		dialogStore,
		screens,
	)

	logger.Debug("initializing feature", zap.String("feature", "payments"))
	paymentsRepository := payments_postgres_repository.NewPaymentsRepository(pool)
	paymentsService := payments_service.NewPaymentsService(
		paymentsRepository,
		subscriptionsService,
		pool,
	)
	paymentsTransportTG := payments_transport_tg.NewPaymentsTGHandler(
		bot,
		paymentsService,
		superAdminOnly,
		screens,
		payments_transport_tg.NewConfigMust(),
		cfg.AdminContact,
		cfg.SuperAdminTelegramID,
	)

	logger.Debug("initializing feature", zap.String("feature", "access"))
	accessRepository := access_postgres_repository.NewAccessRepository(pool)
	accessService := access_service.NewAccessService(accessRepository)
	accessTransportHTTP := access_transport_http.NewAccessHTTPHandler(
		accessService,
		access_transport_http.NewConfigMust(),
		cfg.AdminContact,
	)

	logger.Debug("initializing HTTP server")
	mux := http.NewServeMux()
	accessTransportHTTP.Register(mux)
	httpServer := core_http_server.NewServer(
		core_http_server.NewConfigMust(),
		logger,
		mux,
	)

	logger.Debug("initializing TG poller")

	router := core_tg_poller.NewRouter()
	router.RegisterRoutes(usersTransportTG.Routes()...)
	router.RegisterRoutes(adminsTransportTG.Routes()...)
	router.RegisterRoutes(plansTransportTG.Routes()...)
	router.RegisterRoutes(subscriptionsTransportTG.Routes()...)
	router.RegisterRoutes(paymentsTransportTG.Routes()...)

	flows := make(map[string]core_tg_dialog.Flow)
	for _, flow := range usersTransportTG.Flows() {
		flows[flow.ID] = flow
	}
	for _, flow := range plansTransportTG.Flows() {
		flows[flow.ID] = flow
	}
	for _, flow := range subscriptionsTransportTG.Flows() {
		flows[flow.ID] = flow
	}

	pollerConfig := core_tg_poller.NewConfigMust()
	poller := core_tg_poller.NewTGPoller(
		router,
		bot,
		pollerConfig,
		logger,
		core_tg_middleware.Panic(bot),
		core_tg_middleware.Logger(logger),
		core_tg_middleware.ChatLock(),
		core_tg_middleware.Trace(),
		core_tg_middleware.CallbackAnswer(bot),
		core_tg_middleware.Dialog(bot, dialogStore, flows),
	)

	group, gctx := errgroup.WithContext(ctx)
	group.Go(func() error {
		if err := poller.Run(gctx); err != nil {
			return fmt.Errorf("TG poller: %w", err)
		}
		return nil
	})
	group.Go(func() error {
		if err := httpServer.Run(gctx); err != nil {
			return fmt.Errorf("HTTP server: %w", err)
		}
		return nil
	})

	if err := group.Wait(); err != nil {
		logger.Error("application run error", zap.Error(err))
	}
}
