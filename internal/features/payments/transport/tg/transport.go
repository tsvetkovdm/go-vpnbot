package payments_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_middleware "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/middleware"
	core_tg_poller "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/poller"
	core_tg_screen "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/screen"
)

type PaymentsTGHandler struct {
	bot             *tgbotapi.BotAPI
	paymentsService PaymentsService
	superAdminOnly  core_tg_middleware.Middleware
	screens         *core_tg_screen.Screens
	config          Config
	adminContact    string
	adminChatID     int64
}

type PaymentsService interface {
	CreatePayment(
		ctx context.Context,
		telegramID int64,
		planPriceID int64,
	) (domain.Payment, error)

	ClaimFreePayment(
		ctx context.Context,
		telegramID int64,
		planPriceID int64,
	) (domain.ResolvedPayment, error)

	GetUnresolvedPayment(
		ctx context.Context,
		telegramID int64,
	) (domain.Payment, error)

	GetPaymentByID(
		ctx context.Context,
		paymentID int64,
	) (domain.Payment, error)

	CancelPayment(
		ctx context.Context,
		telegramID int64,
		paymentID int64,
	) (domain.Payment, error)

	SubmitPayment(
		ctx context.Context,
		telegramID int64,
		paymentID int64,
	) (domain.Payment, error)

	ConfirmPayment(
		ctx context.Context,
		paymentID int64,
		adminTelegramID int64,
	) (domain.ResolvedPayment, error)

	RejectPayment(
		ctx context.Context,
		paymentID int64,
		adminTelegramID int64,
	) (domain.ResolvedPayment, error)

	ListPendingPayments(
		ctx context.Context,
		limit int,
	) ([]domain.PaymentOverview, int, error)

	GetPaymentOverview(
		ctx context.Context,
		paymentID int64,
	) (domain.PaymentOverview, error)
}

func NewPaymentsTGHandler(
	bot *tgbotapi.BotAPI,
	paymentsService PaymentsService,
	superAdminOnly core_tg_middleware.Middleware,
	screens *core_tg_screen.Screens,
	config Config,
	adminContact string,
	adminChatID int64,
) *PaymentsTGHandler {
	return &PaymentsTGHandler{
		bot:             bot,
		paymentsService: paymentsService,
		superAdminOnly:  superAdminOnly,
		screens:         screens,
		config:          config,
		adminContact:    adminContact,
		adminChatID:     adminChatID,
	}
}

func (h *PaymentsTGHandler) Routes() []core_tg_poller.Route {
	return []core_tg_poller.Route{
		core_tg_poller.NewRoute(
			"pay:create:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuCreatePayment,
		),
		core_tg_poller.NewRoute(
			"pay:free:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuClaimFreePayment,
		),
		core_tg_poller.NewRoute(
			"pay:current",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuCurrentPayment,
		),
		core_tg_poller.NewRoute(
			"pay:cancel:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuCancelPayment,
		),
		core_tg_poller.NewRoute(
			"pay:submit:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuSubmitPayment,
		),
		core_tg_poller.NewRoute(
			"pay:confirm:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuConfirmPayment,
			h.superAdminOnly,
		),
		core_tg_poller.NewRoute(
			"pay:reject:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuRejectPayment,
			h.superAdminOnly,
		),
		core_tg_poller.NewRoute(
			"pay:dismiss",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuDismissNotification,
			h.superAdminOnly,
		),
		core_tg_poller.NewRoute(
			"pay:admin",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuAdminPayments,
			h.superAdminOnly,
		),
		core_tg_poller.NewRoute(
			"pay:aview:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuAdminPayment,
			h.superAdminOnly,
		),
		core_tg_poller.NewRoute(
			"pay:aconfirm:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuAdminResolvePayment,
			h.superAdminOnly,
		),
		core_tg_poller.NewRoute(
			"pay:areject:",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuAdminResolvePayment,
			h.superAdminOnly,
		),
	}
}
