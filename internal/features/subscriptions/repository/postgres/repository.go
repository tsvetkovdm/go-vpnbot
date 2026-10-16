package subscriptions_postgres_repository

import core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"

type SubscriptionsRepository struct {
	pool core_postgres_pool.Pool
}

func NewSubscriptionsRepository(
	pool core_postgres_pool.Pool,
) *SubscriptionsRepository {
	return &SubscriptionsRepository{
		pool: pool,
	}
}
