package domain

import (
	"time"

	"github.com/google/uuid"
)

type FeedOwner struct {
	UserID    int64
	XrayUUID  uuid.UUID
	ExpiresAt *time.Time
}

type Feed struct {
	Owner       FeedOwner
	Connections []Connection
}

func (f Feed) Links() []string {
	links := make([]string, 0, len(f.Connections))
	for _, c := range f.Connections {
		links = append(links, c.Link(f.Owner.XrayUUID))
	}

	return links
}
