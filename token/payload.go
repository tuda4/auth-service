package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	AccountID uuid.UUID `json:"account_id"`
	ExpiredAt time.Time `json:"expired_at"`
	CreatedAt time.Time `json:"created_at"`
}

func NewPayload(accountID uuid.UUID, duration time.Duration) (payload *Payload, err error) {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return
	}

	payload = &Payload{
		ID:        uuid,
		AccountID: accountID,
		ExpiredAt: time.Now().Add(duration),
		CreatedAt: time.Now(),
	}

	return
}

func (payload Payload) Valid() error {
	if time.Now().After(payload.ExpiredAt) {
		return errors.New("expired token")
	}

	return nil
}
