package access_postgres_repository

import core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"

type AccessRepository struct {
	pool core_postgres_pool.Pool
}

func NewAccessRepository(
	pool core_postgres_pool.Pool,
) *AccessRepository {
	return &AccessRepository{
		pool: pool,
	}
}
