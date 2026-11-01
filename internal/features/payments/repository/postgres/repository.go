package payments_postgres_repository

import core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"

type PaymentsRepository struct {
	pool core_postgres_pool.Pool
}

func NewPaymentsRepository(
	pool core_postgres_pool.Pool,
) *PaymentsRepository {
	return &PaymentsRepository{
		pool: pool,
	}
}
