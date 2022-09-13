package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all functions to excute db queries and transaction
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transaction
// Embedding it in side the Store struct. It's preferred way to extend struct functionality in golang instead of inheritance
type SQLStore struct {
	*Queries
	db *sql.DB
}

// NewStore creates a new store
func NewStore(db *sql.DB) *SQLStore {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

/*
  - Add function to store: a generic database transaction execution
  - execTx: excutes a functionality within a database transaction
  - execTx: Noted that the functionality is unxported (it starts with a lowercase letter e) because we do not want external package to call it directly.
    Instead, we will provide an exported function for each specific transaction
  - func(*Queries) error: get back error

_, blank identifier
*/
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	// input function to get back error
	err = fn(q)

	//if get back an error and if the error is  not nil we need to rollback the transaction
	if err != nil {
		// If the rollback error is also not nill and then we report 2 error
		if rbErr := tx.Rollback(); rbErr != nil {
			//return the rollback error
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		// return original transaction error
		return err
	}

	return tx.Commit()
}

// TransferParams contains the input parameters of the transfer transactions
type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

var txKey = struct{}{}

// TransferTxResult is the result of the transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to the other.
// It creates a trasfer record, adn account entries, and update accounts balance within a single database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		txName := ctx.Value(txKey)
		fmt.Println(txName, "Create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		fmt.Println(txName, "Create entry 1")
		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})

		if err != nil {
			return err
		}

		fmt.Println(txName, "Create entry 2")

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    +arg.Amount,
		})

		if err != nil {
			return err
		}

		// fmt.Println(txName, "get update account 2")
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})

			if err != nil {
				return err
			}

			fmt.Println(txName, "update account 2")
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})
			if err != nil {
				return err
			}
		} else {
			fmt.Println(txName, "update account 2")
			result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.ToAccountID,
				Amount: arg.Amount,
			})

			if err != nil {
				return err
			}

			result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
				ID:     arg.FromAccountID,
				Amount: -arg.Amount,
			})

			if err != nil {
				return err
			}
		}
		return nil
	})
	return result, err
}
