package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_config "github.com/tsvetkovdm/go-vpnbot/internal/core/config"
	core_logger "github.com/tsvetkovdm/go-vpnbot/internal/core/logger"
	core_pgx_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool/pgx"
	core_tg_client "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/client"
	core_tg_middleware "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/middleware"
	core_tg_poller "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/poller"
	admins_postgres_repository "github.com/tsvetkovdm/go-vpnbot/internal/features/admins/repository/postgres"
	admins_service "github.com/tsvetkovdm/go-vpnbot/internal/features/admins/service"
	admins_transport_tg "github.com/tsvetkovdm/go-vpnbot/internal/features/admins/transport/tg"
	users_postgres_repository "github.com/tsvetkovdm/go-vpnbot/internal/features/users/repository/postgres"
	users_service "github.com/tsvetkovdm/go-vpnbot/internal/features/users/service"
	users_transport_tg "github.com/tsvetkovdm/go-vpnbot/internal/features/users/transport/tg"
	"go.uber.org/zap"
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

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "admins"))
	adminsRepository := admins_postgres_repository.NewAdminsRepository(pool)
	adminsService := admins_service.NewAdminsService(adminsRepository)
	adminsTransportTG := admins_transport_tg.NewAdminsTGHandler(bot, adminsService)
	adminOnly := core_tg_middleware.AdminOnly(bot, adminsService)

	if err := adminsService.EnsureSuperAdmin(ctx, cfg.SuperAdminTelegramID); err != nil {
		logger.Fatal("failed to ensure super admin", zap.Error(err))
	}

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportTG := users_transport_tg.NewUsersTGHandler(bot, usersService, adminOnly)

	logger.Debug("initializing TG poller")

	router := core_tg_poller.NewRouter()
	router.RegisterRoutes(usersTransportTG.Routes()...)
	router.RegisterRoutes(adminsTransportTG.Routes()...)

	pollerConfig := core_tg_poller.NewConfigMust()
	poller := core_tg_poller.NewTGPoller(
		router,
		bot,
		pollerConfig,
		logger,
		core_tg_middleware.Logger(logger),
		core_tg_middleware.Trace(),
		core_tg_middleware.Panic(bot),
	)

	if err := poller.Run(ctx); err != nil {
		logger.Error("TG poller run error", zap.Error(err))
	}
}
