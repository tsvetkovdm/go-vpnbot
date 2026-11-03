package access_postgres_repository

import (
	"time"

	domain "github.com/tsvetkovdm/go-vpnbot/internal/core/domains"
	core_postgres_pool "github.com/tsvetkovdm/go-vpnbot/internal/core/repository/postgres/pool"
)

const accessRule = `
	FROM vpnbot.connections c
	CROSS JOIN vpnbot.users u
	LEFT JOIN vpnbot.subscriptions s ON s.user_id = u.id
	LEFT JOIN vpnbot.access_grants g ON g.connection_id = c.id AND g.user_id = u.id
	WHERE c.is_active
	  AND u.status = 'active'
	  AND ((c.is_public AND s.expires_at > now()) OR g.expires_at > now())`

const connectionColumns = `
	c.id, c.title, c.inbound_tag, c.link_template,
	c.is_public, c.is_active, c.sort_order,
	c.created_at, c.updated_at`

type ConnectionModel struct {
	ID           int64
	Title        string
	InboundTag   string
	LinkTemplate string
	IsPublic     bool
	IsActive     bool
	SortOrder    int32
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (m ConnectionModel) toDomain() domain.Connection {
	return domain.Connection{
		ID:           m.ID,
		Title:        m.Title,
		InboundTag:   m.InboundTag,
		LinkTemplate: m.LinkTemplate,
		IsPublic:     m.IsPublic,
		IsActive:     m.IsActive,
		SortOrder:    int(m.SortOrder),
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func scanConnection(row core_postgres_pool.Row, extra ...any) (ConnectionModel, error) {
	var m ConnectionModel
	dest := []any{
		&m.ID,
		&m.Title,
		&m.InboundTag,
		&m.LinkTemplate,
		&m.IsPublic,
		&m.IsActive,
		&m.SortOrder,
		&m.CreatedAt,
		&m.UpdatedAt,
	}
	err := row.Scan(append(dest, extra...)...)
	return m, err
}
