package token

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/meomeocoj/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func TestJWTToken(t *testing.T) {
	secret := utils.RandomString(36)
	username := utils.RandomString(6)

	jwtMaker, err := NewJWTMaker(secret)
	require.NoError(t, err)

	issuedAt := time.Now()

	duration := time.Minute
	expiredAt := time.Now().Add(duration)

	token, err := jwtMaker.CreateToken(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := jwtMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload.ID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}

func TestJWTExpiredToken(t *testing.T) {
	jwtMaker, err := NewJWTMaker(utils.RandomString(36))
	require.NoError(t, err)

	token, err := jwtMaker.CreateToken(utils.RandomString(6), -time.Minute)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := jwtMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredDate.Error())
	require.Nil(t, payload)

}

func TestInvalidJWTAlgoNone(t *testing.T) {
	payload, err := NewPayload(utils.RandomOwner(), time.Minute)
	require.NoError(t, err)

	token := jwt.NewWithClaims(jwt.SigningMethodNone, payload)
	signedToken, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	require.NoError(t, err)
	jwtMaker, err := NewJWTMaker(utils.RandomString(36))
	require.NoError(t, err)
	payload, err = jwtMaker.VerifyToken(signedToken)
	require.Error(t, err)
	require.EqualError(t, err, ErrInvalidToken.Error())
	require.Nil(t, payload)

}
