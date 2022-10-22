package token

import (
	"fmt"
	"time"

	"github.com/o1egl/paseto"
	"golang.org/x/crypto/chacha20poly1305"
)

var ErrKeySize = fmt.Errorf("invalid key size: must be exactly %d characters", chacha20poly1305.KeySize)

type PasetoMaker struct {
	pasto       *paseto.V2
	symetricKey []byte
}

func NewPasetoMaker(symetricKey string) (Maker, error) {
	if len(symetricKey) != chacha20poly1305.KeySize {
		return nil, ErrKeySize
	}

	return &PasetoMaker{
		pasto:       paseto.NewV2(),
		symetricKey: []byte(symetricKey),
	}, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {

	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	token, err := maker.pasto.Encrypt(maker.symetricKey, payload, nil)
	return token, err

}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	err := maker.pasto.Decrypt(token, maker.symetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()

	if err != nil {
		return nil, err
	}

	return payload, nil
}
