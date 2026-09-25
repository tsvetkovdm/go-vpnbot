package plans_postgres_repository

import core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"

type PlansRepository struct {
	pool core_postgres_pool.Pool
}

func NewPlansRepository(
	pool core_postgres_pool.Pool,
) *PlansRepository {
	return &PlansRepository{
		pool: pool,
	}
}
