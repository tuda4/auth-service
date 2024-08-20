package db

import (
	"context"

	"github.com/google/uuid"
)

type VerifyEmailTxParams struct {
	AccountID  uuid.UUID `json:"account_id"`
	SecretCode string    `json:"secret_code"`
}

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) error {
	err := store.execTx(ctx, func(q *Queries) error {
		err := q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			AccountID:  arg.AccountID,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}

		err = q.UpdateAccountEmail(ctx, arg.AccountID)
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
