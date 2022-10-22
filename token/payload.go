package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrExpiredDate = errors.New("expired token")

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issuedAt"`
	ExpiredAt time.Time `json:"expiredAt"`
}

func (payload *Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredDate
	}
	return nil
}

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return nil, ErrFailToGenerateUUID
	}
	return &Payload{
		ID:        uuid,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}, nil
}
