package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	s := NewStore(db)
	acc1 := createTestAccount(t)
	acc2 := createTestAccount(t)
	n := 8
	amount := int64(10)
	errors := make(chan error)
	results := make(chan TransferTxResult)

	existed := make(map[int]bool)

	for i := 0; i < n; i++ {
		ctx := context.Background()
		go func() {
			result, err := s.TransferTransaction(ctx, CreateTransferParams{
				FromAccountID: acc1.ID,
				ToAccountID:   acc2.ID,
				Amount:        amount,
			})
			errors <- err
			results <- result
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)
		result := <-results
		// check transfer
		transfer := result.Transfer
		require.Equal(t, transfer.FromAccountID, acc1.ID)
		require.Equal(t, transfer.ToAccountID, acc2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotEmpty(t, transfer)
		require.NotZero(t, transfer.CreatedAt)
		require.NotZero(t, transfer.ID)

		_, err = s.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		// check from entry
		from := result.FromEntry
		require.NotEmpty(t, from)
		require.Equal(t, -amount, from.Amount)
		require.Equal(t, acc1.ID, from.AccountID)
		require.NotZero(t, from.CreatedAt)
		require.NotZero(t, from.ID)

		_, err = s.GetEntry(context.Background(), from.ID)
		require.NoError(t, err)
		// check to entry
		to := result.ToEntry
		require.NotEmpty(t, to)
		require.Equal(t, amount, to.Amount)
		require.Equal(t, acc2.ID, to.AccountID)
		require.NotZero(t, to.CreatedAt)
		require.NotZero(t, to.ID)
		_, err = s.GetEntry(context.Background(), to.ID)
		require.NoError(t, err)

		// TODO: check update balance
		fromAccount := result.FromAccount
		toAccount := result.ToAccount
		diff1 := acc1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - acc2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff2 > 0)
		require.True(t, diff1%amount == 0)
		k := int(diff1 / amount)

		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	updatedAcc1, _ := s.GetAccount(context.Background(), acc1.ID)
	updatedAcc2, _ := s.GetAccount(context.Background(), acc2.ID)
	require.Equal(t, acc1.Balance-int64(n)*amount, updatedAcc1.Balance)
	require.Equal(t, acc2.Balance+int64(n)*amount, updatedAcc2.Balance)

}

func TestDeadlockTransferTx(t *testing.T) {
	s := NewStore(db)
	acc1 := createTestAccount(t)
	acc2 := createTestAccount(t)
	// run n concurrent transactions
	n := 6
	amount := int64(10)
	errors := make(chan error)

	for i := 0; i < n; i++ {
		fromID := acc1.ID
		toID := acc2.ID
		if i%2 == 0 {
			fromID = acc2.ID
			toID = acc1.ID
		}
		ctx := context.Background()
		go func(fromId, toID int64) {
			_, err := s.TransferTransaction(ctx, CreateTransferParams{
				FromAccountID: fromID,
				ToAccountID:   toID,
				Amount:        amount,
			})
			errors <- err
		}(fromID, toID)
	}

	for i := 0; i < n; i++ {
		err := <-errors
		require.NoError(t, err)
	}

	updatedAcc1, _ := s.GetAccount(context.Background(), acc1.ID)
	updatedAcc2, _ := s.GetAccount(context.Background(), acc2.ID)
	require.Equal(t, acc1.Balance, updatedAcc1.Balance)
	require.Equal(t, acc2.Balance, updatedAcc2.Balance)

}
