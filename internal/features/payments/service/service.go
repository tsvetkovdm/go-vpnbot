package payments_service

import (
	"context"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

type PaymentsService struct {
	paymentsRepository   PaymentsRepository
	subscriptionsService SubscriptionsService
	transactor           core_postgres_pool.Transactor
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

	ListPendingPayments(
		ctx context.Context,
		limit int,
	) ([]domain.PaymentOverview, int, error)

	GetPaymentOverview(
		ctx context.Context,
		paymentID int64,
	) (domain.PaymentOverview, error)

	ResolvePayment(
		ctx context.Context,
		paymentID int64,
		adminTelegramID int64,
		status domain.PaymentStatus,
	) (domain.Payment, int64, error)
}

type SubscriptionsService interface {
	Extend(
		ctx context.Context,
		userID int64,
		d domain.Duration,
		expectedVersion int64,
	) (domain.Subscription, error)
}

func NewPaymentsService(
	paymentsRepository PaymentsRepository,
	subscriptionsService SubscriptionsService,
	transactor core_postgres_pool.Transactor,
) *PaymentsService {
	return &PaymentsService{
		paymentsRepository:   paymentsRepository,
		subscriptionsService: subscriptionsService,
		transactor:           transactor,
	}
}
