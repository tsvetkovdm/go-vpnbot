package payments_postgres_repository

import (
	"context"
	"fmt"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
)

func (r *PaymentsRepository) ListPendingPayments(
	ctx context.Context,
	limit int,
) ([]domain.PaymentOverview, int, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
	SELECT ` + paymentOverviewColumns + `,
		   COUNT(*) OVER () AS total
	FROM vpnbot.payments p
	JOIN vpnbot.users u ON u.id = p.user_id
	WHERE p.status = 'pending'
	ORDER BY p.updated_at ASC, p.id ASC
	LIMIT $1;
	`

	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var (
		overviews []domain.PaymentOverview
		total     int
	)
	for rows.Next() {
		overview, err := scanPaymentOverview(rows, &total)
		if err != nil {
			return nil, 0, fmt.Errorf("scan error: %w", err)
		}
		overviews = append(overviews, overview)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	return overviews, total, nil
}
