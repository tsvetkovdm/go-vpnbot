package domain

import (
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

const connectionUUIDPlaceholder = "{uuid}"

type Connection struct {
	ID           int64
	Title        string
	InboundTag   string
	LinkTemplate string
	IsPublic     bool
	IsActive     bool
	SortOrder    int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (c Connection) Link(xrayUUID uuid.UUID) string {
	link := strings.Replace(c.LinkTemplate, connectionUUIDPlaceholder, xrayUUID.String(), 1)

	return link + "#" + url.PathEscape(c.Title)
}
