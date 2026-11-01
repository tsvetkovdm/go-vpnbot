package payments_service

import (
	"context"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

type PaymentsService struct {
	paymentsRepository PaymentsRepository
}

type PaymentsRepository interface {
	GetUserStatusByTelegramID(
		ctx context.Context,
		telegramID int64,
	) (int64, domain.Status, error)

	GetPlanPriceForPayment(
		ctx context.Context,
		planPriceID int64,
	) (domain.Plan, domain.PlanPrice, error)

	CreatePayment(
		ctx context.Context,
		payment domain.Payment,
	) (domain.Payment, error)

	GetUnresolvedPaymentByTelegramID(
		ctx context.Context,
		telegramID int64,
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
}

func NewPaymentsService(
	paymentsRepository PaymentsRepository,
) *PaymentsService {
	return &PaymentsService{
		paymentsRepository: paymentsRepository,
	}
}
