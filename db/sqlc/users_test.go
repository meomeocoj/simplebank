package db

import (
	"context"
	"testing"

	"github.com/meomeocoj/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createTestUser(t *testing.T) User {
	args := CreateUserParams{
		Username:     utils.RandomOwner(),
		HashPassword: utils.RandomString(20),
		Fullname:     utils.RandomString(15),
		Email:        utils.RandomEmail(),
	}

	user, err := testingQueries.CreateUser(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, user.Username, args.Username)
	require.Equal(t, user.HashPassword, args.HashPassword)
	require.Equal(t, user.Fullname, args.Fullname)
	require.Equal(t, user.Email, args.Email)
	require.NotZero(t, user.CreatedAt)
	return user
}
func TestCreateUser(t *testing.T) {
	createTestUser(t)
}

func TestGetUser(t *testing.T) {
	args := createTestUser(t)
	user, err := testingQueries.GetUser(context.Background(), args.Username)
	require.NoError(t, err)
	require.Equal(t, user.Username, args.Username)
	require.Equal(t, user.HashPassword, args.HashPassword)
	require.Equal(t, user.Email, args.Email)
	require.Equal(t, user.Fullname, args.Fullname)

}
