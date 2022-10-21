package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTransaction(ctx context.Context, args CreateTransferParams) (TransferTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

// var txKey struct{}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db)}
}

func (s *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rError := tx.Rollback(); rError != nil {
			return fmt.Errorf("tx error: %v, rollback error: %v", err, rError)
		}
		return err
	}
	return tx.Commit()
}

type TransferTxResult struct {
	Transfer    Transfer `json:"amount"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"fromEntry"`
	ToEntry     Entry    `json:"toEntry"`
}

func (s *SQLStore) TransferTransaction(ctx context.Context, args CreateTransferParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := s.execTx(ctx, func(q *Queries) error {
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: args.FromAccountID,
			ToAccountID:   args.ToAccountID,
			Amount:        args.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.FromAccountID,
			Amount:    -args.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: args.ToAccountID,
			Amount:    args.Amount,
		})
		if err != nil {
			return err
		}
		if args.FromAccountID < args.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, args.FromAccountID, args.ToAccountID, -args.Amount, args.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, args.ToAccountID, args.FromAccountID, args.Amount, -args.Amount)
		}
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return result, err
	}
	return result, nil
}

func addMoney(ctx context.Context, q *Queries, id1, id2 int64, amount1, amount2 int64) (account1, account2 Account, err error) {
	account1, err = q.AddBalanceToAccount(ctx, AddBalanceToAccountParams{
		ID:     id1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddBalanceToAccount(ctx, AddBalanceToAccountParams{
		ID:     id2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return
}
