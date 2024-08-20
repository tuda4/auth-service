package token

import (
	"errors"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/google/uuid"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (maker Maker, err error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		err = errors.New("invalid key size")
		return
	}

	maker = &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return
}

func (maker *PasetoMaker) CreateToken(accountID uuid.UUID, duration time.Duration) (token string, payload *Payload, err error) {
	payload, err = NewPayload(accountID, duration)
	if err != nil {
		err = errors.New("cannot created payload")
		return
	}
	token, err = maker.paseto.Encrypt(maker.symmetricKey, payload, "")
	return
}

func (maker *PasetoMaker) VerifyToken(token string) (payload *Payload, err error) {
	err = maker.paseto.Decrypt(token, maker.symmetricKey, payload, "")
	if err != nil {
		return nil, errors.New("invalid token")
	}

	if err = payload.Valid(); err != nil {
		return
	}

	return
}
