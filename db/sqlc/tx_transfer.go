package db

import "context"

// TransferTxParams contains the input parameters of the transfer transaction.
type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result of the transfer transaction.
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// TransferTx performs a money transfer from one account to the other.
// It creates a transfer record, add account entries, and update accounts' balance within a single database transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	var err error
	// txName := ctx.Value(txKey)

	err = store.execTx(context.Background(), func(queries *Queries) error {
		// 转账表：增加一条转账记录
		// fmt.Println(txName, " create transfer")
		result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}
		// fmt.Println(txName, " create entry 1")
		// 条目表：增加一条从from_account_id的条目记录 -amount
		result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		// 条目表：增加一条从to_account_id的条目记录 +amount
		// fmt.Println(txName, " create entry2")
		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return nil
		}

		// TODO: update accounts' balance
		// fmt.Println(txName, " get account1")
		//account1, err := queries.GetAccountForUpdate(ctx, arg.FromAccountId)
		//if err != nil {
		//	return err
		//}
		// fmt.Println(txName, " update account1")
		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, queries, arg.FromAccountId, -arg.Amount, arg.ToAccountId, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, queries, arg.ToAccountId, arg.Amount, arg.FromAccountId, -arg.Amount)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return result, nil
}

// addMoney add money to two accounts
func addMoney(ctx context.Context, q *Queries, accountID1, amount1, accountID2, amount2 int64) (account1, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return
}
