package domain

import "time"

type Role string

const (
	RoleSuper   Role = "super"
	RoleAdmin   Role = "admin"
	RoleSupport Role = "support"
)

type Admin struct {
	ID         int64
	TelegramID int64
	Role       Role
	AddedBy    *int64
	CreatedAt  time.Time
}

func NewAdmin(
	telegramID int64,
	role Role,
	addedBy *int64,
) Admin {
	return Admin{
		ID:         UninitializedID,
		TelegramID: telegramID,
		Role:       role,
		AddedBy:    addedBy,
		CreatedAt:  time.Now(),
	}
}
