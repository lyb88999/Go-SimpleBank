package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
)

type VerifyEmailTxParams struct {
	EmailID    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult
	var err error
	err = store.execTx(context.Background(), func(queries *Queries) error {
		result.VerifyEmail, err = queries.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailID,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}
		result.User, err = queries.UpdateUser(ctx, UpdateUserParams{
			Username:        result.VerifyEmail.Username,
			IsEmailVerified: pgtype.Bool{Bool: true, Valid: true},
		})
		return err
	})
	return result, err
}
