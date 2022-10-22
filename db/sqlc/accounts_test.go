package db

import (
	"context"
	"testing"

	"github.com/meomeocoj/simplebank/utils"
	"github.com/stretchr/testify/require"
)

func createTestAccount(t *testing.T) Account {
	user := createTestUser(t)

	args := CreateAccountParams{
		Owner:    user.Username,
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	acc, err := testingQueries.CreateAccount(context.Background(), args)

	require.NoError(t, err)
	require.Equal(t, acc.Balance, args.Balance)
	require.Equal(t, acc.Owner, args.Owner)
	require.Equal(t, acc.Currency, args.Currency)
	require.NotZero(t, acc.CreatedAt)
	require.NotZero(t, acc.ID)
	return acc
}
func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}

func TestGetAccount(t *testing.T) {
	args := createTestAccount(t)
	acc, err := testingQueries.GetAccount(context.Background(), args.ID)
	require.NoError(t, err)
	require.Equal(t, acc.Balance, args.Balance)
	require.Equal(t, acc.Owner, args.Owner)
}

func TestListAccounts(t *testing.T) {
	var lastAccount Account
	for i := 0; i < 10; i++ {
		lastAccount = createTestAccount(t)
	}

	args := ListAccountsParams{
		Owner:  lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}
	accs, err := testingQueries.ListAccounts(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, accs)
	for _, acc := range accs {

		require.NotEmpty(t, acc.ID)
		require.Equal(t, lastAccount.Owner, acc.Owner)
	}

}

func TestDeleteAccount(t *testing.T) {
	args := createTestAccount(t)
	testingQueries.DeleteAccount(context.Background(), args.ID)
	acc, err := testingQueries.GetAccount(context.Background(), args.ID)
	require.Error(t, err)
	require.Equal(t, acc, *new(Account))
}

func TestUpdateAccount(t *testing.T) {
	acc := createTestAccount(t)
	args := UpdateAccountParams{
		ID:      acc.ID,
		Balance: 456,
	}
	updatedAcc, err := testingQueries.UpdateAccount(context.Background(), args)
	require.NoError(t, err)
	require.Equal(t, updatedAcc.Balance, int64(456))
}
