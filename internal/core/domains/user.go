package domain

import (
	"crypto/rand"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusActive  Status = "active"
	StatusBlocked Status = "blocked"
)

const UserEmailPrefix = "user-"

type User struct {
	ID         int64
	Version    int64
	TelegramID int64
	Username   *string
	FirstName  string
	Email      string
	XrayUUID   uuid.UUID
	Status     Status
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func NewUser(
	telegramID int64,
	username *string,
	firstName string,
) User {
	xrayUUID := uuid.New()
	email := generateEmail()
	now := time.Now()

	return User{
		ID:         UninitializedID,
		Version:    UninitializedVersion,
		TelegramID: telegramID,
		Username:   username,
		FirstName:  firstName,
		Email:      email,
		XrayUUID:   xrayUUID,
		Status:     StatusActive,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

func generateEmail() string {
	b := make([]byte, 5)
	_, _ = rand.Read(b)
	return fmt.Sprintf("%s%x", UserEmailPrefix, b)
}

func IsBotUserEmail(email string) bool {
	return strings.HasPrefix(email, UserEmailPrefix)
}
