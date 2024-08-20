package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/tuda4/mb-backend/util"
)

type AccountTxParams struct {
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	PhoneNumber  string    `json:"phone_number"`
	Address      string    `json:"address"`
	BirthDay     time.Time `json:"birthday"`
	AfterAccount func(account Account) error
}

func (store *SQLStore) CreateAccountTx(ctx context.Context, arg AccountTxParams) (bool, error) {
	err := store.execTx(ctx, func(q *Queries) error {
		uuid, err := uuid.NewUUID()
		if err != nil {
			return err
		}

		hashPassword, err := util.CreateHashPassword(arg.Password)
		if err != nil {
			return err
		}
		accountParam := CreateAccountParams{
			AccountID:    uuid,
			Email:        arg.Email,
			HashPassword: hashPassword,
		}
		account, err := q.CreateAccount(ctx, accountParam)
		if err != nil {
			return err
		}

		_, err = q.CreateProfile(ctx, CreateProfileParams{
			AccountID:   account.AccountID,
			PhoneNumber: arg.PhoneNumber,
			Birthday:    arg.BirthDay,
			FirstName:   arg.FirstName,
			LastName:    arg.LastName,
			Address:     pgtype.Text{String: arg.Address},
		})
		if err != nil {
			return err
		}

		err = arg.AfterAccount(account)
		fmt.Print(err)
		return err
	})

	return true, err
}
