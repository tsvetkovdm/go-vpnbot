package domain

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/shopspring/decimal"
	core_errors "github.com/tsvetkovdm/go-vpnbot/internal/core/errors"
)

type PaymentStatus string

const (
	PaymentStatusCreated    PaymentStatus = "created"
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusSuccessful PaymentStatus = "successful"
	PaymentStatusRejected   PaymentStatus = "rejected"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
	PaymentStatusExpired    PaymentStatus = "expired"
)

const (
	PaymentCommentCodeMin = 1000
	PaymentCommentCodeMax = 9999
)

var (
	ErrPaymentPlanUnavailable  = fmt.Errorf("plan or plan price is unavailable: %w", core_errors.ErrNotFound)
	ErrPaymentUnresolvedExists = fmt.Errorf("user already has unresolved payment: %w", core_errors.ErrConflict)
	ErrPaymentUserBlocked      = fmt.Errorf("user is blocked: %w", core_errors.ErrConflict)
	ErrPaymentAlreadyResolved  = fmt.Errorf("payment is already resolved: %w", core_errors.ErrConflict)
)

type Payment struct {
	ID                int64
	UserID            int64
	PlanPriceID       int64
	PlanTitle         string
	Price             decimal.Decimal
	Duration          Duration
	CommentCode       int
	Status            PaymentStatus
	ResolvedByAdminID *int64
	ResolvedAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ResolvedPayment struct {
	Payment        Payment
	UserTelegramID int64
	Subscription   *Subscription
}

func NewPayment(
	userID int64,
	plan Plan,
	price PlanPrice,
) (Payment, error) {
	if !plan.IsActive || !price.IsActive || price.PlanID != plan.ID {
		return Payment{}, fmt.Errorf(
			"%w: plan_id=%d plan_price_id=%d",
			ErrPaymentPlanUnavailable, plan.ID, price.ID,
		)
	}

	return Payment{
		ID:          UninitializedID,
		UserID:      userID,
		PlanPriceID: price.ID,
		PlanTitle:   plan.Title,
		Price:       price.Price,
		Duration:    price.Duration,
		CommentCode: generateCommentCode(),
		Status:      PaymentStatusCreated,
	}, nil
}

func generateCommentCode() int {
	return PaymentCommentCodeMin + rand.IntN(PaymentCommentCodeMax-PaymentCommentCodeMin+1)
}
