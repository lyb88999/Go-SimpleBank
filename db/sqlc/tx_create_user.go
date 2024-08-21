package db

import "context"

// CreateUserTxParams contains the input parameters of the CreateUser transaction.
type CreateUserTxParams struct {
	CreateUserParams
	AfterCreate func(user User) error
}

// CreateUserTxResult is the result of the CreateUser transaction.
type CreateUserTxResult struct {
	User User
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult
	var err error
	err = store.execTx(context.Background(), func(queries *Queries) error {
		result.User, err = queries.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}
		if err != nil {
			return err
		}
		return arg.AfterCreate(result.User)
	})
	return result, err
}
