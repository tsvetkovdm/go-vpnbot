package domain

import "github.com/google/uuid"

type AccessEntry struct {
	InboundTag string
	Email      string
	XrayUUID   uuid.UUID
}
