package token

import (
	"testing"
	"time"

	"github.com/meomeocoj/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func TestPasetoToken(t *testing.T) {
	secret := utils.RandomString(32)
	username := utils.RandomString(6)

	pasetoMaker, err := NewPasetoMaker(secret)
	require.NoError(t, err)

	issuedAt := time.Now()

	duration := time.Minute
	expiredAt := time.Now().Add(duration)

	token, err := pasetoMaker.CreateToken(username, duration)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := pasetoMaker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload.ID)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

}

func TestPasetoExpiredToken(t *testing.T) {
	pasetoMaker, err := NewPasetoMaker(utils.RandomString(32))
	require.NoError(t, err)

	token, err := pasetoMaker.CreateToken(utils.RandomString(6), -time.Minute)

	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := pasetoMaker.VerifyToken(token)
	require.Error(t, err)
	require.EqualError(t, err, ErrExpiredDate.Error())
	require.Nil(t, payload)

}
