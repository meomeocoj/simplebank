package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const MIN_SECRET_LENGTH = 32

var (
	ErrMinSecretLength    = errors.New("too short secret")
	ErrFailToGenerateUUID = errors.New("fail to generate uuid")
	ErrInvalidToken       = errors.New("invalid token")
)

type JWTMaker struct {
	secret string `json:"secret"`
}

func NewJWTMaker(secret string) (Maker, error) {
	if len(secret) < MIN_SECRET_LENGTH {
		return nil, ErrMinSecretLength
	}
	return &JWTMaker{
		secret: secret,
	}, nil

}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {

	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signedToken, err := token.SignedString([]byte(maker.secret))
	return signedToken, err

}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secret), nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredDate) {
			return nil, ErrExpiredDate
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrInvalidToken
	}
	return payload, nil
}
