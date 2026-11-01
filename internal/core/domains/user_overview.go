package domain

type UserOverview struct {
	User         User
	Subscription *Subscription
	IsAdmin      bool

	UnresolvedPaymentStatus *PaymentStatus
}

type UserFilter string

const (
	UserFilterAll      UserFilter = "all"
	UserFilterActive   UserFilter = "active"
	UserFilterExpiring UserFilter = "expiring"
	UserFilterInactive UserFilter = "inactive"
	UserFilterBlocked  UserFilter = "blocked"
)
