package admins_postgres_repository

import (
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

type AdminsRepository struct {
	pool core_postgres_pool.Pool
}

func NewAdminsRepository(
	pool core_postgres_pool.Pool,
) *AdminsRepository {
	return &AdminsRepository{
		pool: pool,
	}
}
