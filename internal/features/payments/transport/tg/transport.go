package payments_transport_tg

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_tg_poller "github.com/tsvetkovdm/go-vpnbot/internal/core/transport/tg/poller"
)

type PaymentsTGHandler struct {
	bot             *tgbotapi.BotAPI
	paymentsService PaymentsService
	config          Config
	adminContact    string
}

type PaymentsService interface {
	CreatePayment(
		ctx context.Context,
		telegramID int64,
		planPriceID int64,
	) (domain.Payment, error)

	GetUnresolvedPayment(
		ctx context.Context,
		telegramID int64,
	) (domain.Payment, error)
}

func NewPaymentsTGHandler(
	bot *tgbotapi.BotAPI,
	paymentsService PaymentsService,
	config Config,
	adminContact string,
) *PaymentsTGHandler {
	return &PaymentsTGHandler{
		bot:             bot,
		paymentsService: paymentsService,
		config:          config,
		adminContact:    adminContact,
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
			"pay:current",
			core_tg_poller.TriggerCallbackPrefix,
			h.MenuCurrentPayment,
		),
	}
}
