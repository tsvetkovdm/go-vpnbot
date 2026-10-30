package plans_postgres_repository

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (r *PlansRepository) CreatePlan(
	ctx context.Context,
	plan domain.Plan,
) (domain.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	INSERT INTO vpnbot.plans (title, is_active)
	VALUES ($1, $2)
	RETURNING ` + planColumns + `;
	`

	row := r.pool.QueryRow(ctx, query, plan.Title, plan.IsActive)

	m, err := scanPlan(row)
	if err != nil {
		return domain.Plan{}, fmt.Errorf("scan error: %w", err)
	}

	return m.toDomain(), nil
}
