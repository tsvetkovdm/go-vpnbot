package domain

import (
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestNewUser(t *testing.T) {
	username := "username"

	testCases := []struct {
		name       string
		telegramID int64
		username   *string
	}{
		{
			name:       "with username",
			telegramID: 123456789,
			username:   &username,
		},
		{
			name:       "without username",
			telegramID: 987654321,
			username:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			user := NewUser(tc.telegramID, tc.username)

			if user.TelegramID != tc.telegramID {
				t.Errorf("TelegramID = %d, want %d", user.TelegramID, tc.telegramID)
			}
			switch {
			case tc.username == nil:
				if user.Username != nil {
					t.Errorf("Username = %v, want nil", *user.Username)
				}
			case user.Username == nil:
				t.Errorf("Username = nil, want %v", *tc.username)
			case *user.Username != *tc.username:
				t.Errorf("Username = %v, want %v", *user.Username, *tc.username)
			}
			if user.ID != UninitializedID {
				t.Errorf("ID = %d, want %d", user.ID, UninitializedID)
			}
			if user.Version != UninitializedVersion {
				t.Errorf("Version = %d, want %d", user.Version, UninitializedVersion)
			}
			if user.Status != StatusActive {
				t.Errorf("Status = %v, want %v", user.Status, StatusActive)
			}
			if !strings.HasPrefix(user.Email, "user-") {
				t.Errorf("Email = %q, want prefix \"user-\"", user.Email)
			}
			if user.XrayUUID == uuid.Nil {
				t.Errorf("XrayUUID is empty")
			}
			if !user.CreatedAt.Equal(user.UpdatedAt) {
				t.Errorf("CreatedAt and UpdatedAt should be equal")
			}
		})
	}
}
