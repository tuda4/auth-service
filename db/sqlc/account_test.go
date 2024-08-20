package db

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tuda4/mb-backend/util"
)

func createTestAccount(t *testing.T) (account Account) {
	hashPassword, err := util.CreateHashPassword(util.RandomString(8))
	require.NoError(t, err)
	uuid, err := uuid.NewUUID()
	require.NoError(t, err)
	email := util.CreateRandomEmail()
	account, err = testQueries.CreateAccount(context.Background(), CreateAccountParams{
		AccountID:    uuid,
		Email:        email,
		HashPassword: hashPassword,
	})

	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, uuid, account.AccountID)
	require.Equal(t, email, account.Email)
	require.NotNil(t, account.CreatedAt)
	require.NotNil(t, account.UpdatedAt)
	require.False(t, account.DeletedAt.Valid)
	require.Equal(t, hashPassword, account.HashPassword)

	return
}

func TestCreateAccount(t *testing.T) {
	createTestAccount(t)
}

func TestChangePassword(t *testing.T) {
	account := createTestAccount(t)
	hashPassword, err := util.CreateHashPassword(util.RandomString(10))
	require.NoError(t, err)
	err = testQueries.ChangePassword(context.Background(), ChangePasswordParams{
		AccountID:    account.AccountID,
		HashPassword: hashPassword,
	})
	require.NoError(t, err)
}

func TestGetProfile(t *testing.T) {
	account := createTestAccount(t)
	profile, err := testQueries.GetProfileAccount(context.Background(), account.Email)
	require.NoError(t, err)
	require.NotEmpty(t, profile)
	require.Equal(t, account.Email, profile.Email)
}
