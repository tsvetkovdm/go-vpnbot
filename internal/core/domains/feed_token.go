package domain

import (
	"crypto/rand"
	"encoding/base64"
)

const feedTokenBytes = 24

func NewFeedToken() string {
	b := make([]byte, feedTokenBytes)
	_, _ = rand.Read(b)

	return base64.RawURLEncoding.EncodeToString(b)
}
