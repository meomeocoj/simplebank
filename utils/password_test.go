package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPassword(t *testing.T) {
	password := RandomString(10)

	hashPassword, err := HashPassword(password)

	require.NoError(t, err)
	require.NotEmpty(t, hashPassword)

	err = CheckIsHashPassword(hashPassword, password)

	require.NoError(t, err)

	wrongPassword := RandomString(10)

	err = CheckIsHashPassword(wrongPassword, hashPassword)
	require.Error(t, err)

}
