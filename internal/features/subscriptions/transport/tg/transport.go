package subscriptions_transport_tg

import (
	"context"
	"time"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

type SubscriptionTGHandler struct {
}

type SubscriptionsService interface {
	Extend(
		ctx context.Context,
		userID int64,
		d domain.Duration,
		expectedVersion int64,
	) (domain.Subscription, error)

	ExtendDays(
		ctx context.Context,
		userID int64,
		days int,
		expectedVersion int64,
	) (domain.Subscription, error)

	SetPaidUntil(
		ctx context.Context,
		userID int64,
		date time.Time,
		expectedVersion int64,
	) (domain.Subscription, error)

	ResolveUser(
		ctx context.Context,
		telegramID int64,
	) (int64, error)
}
