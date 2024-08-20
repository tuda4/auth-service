package db

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"github.com/tuda4/mb-backend/util"
)

func createTestProfile(t *testing.T) (profile Profile) {
	account := createTestAccount(t)
	firstName := util.RandomString(6)
	lastName := util.RandomString(8)
	address := util.RandomString(12)
	phone := util.CreateRandomPhone()
	birthday := util.CreateRandomBirthday()
	arg := CreateProfileParams{
		AccountID:   account.AccountID,
		PhoneNumber: phone,
		Birthday:    birthday,
		FirstName:   firstName,
		LastName:    lastName,
	}
	if address != "" {
		arg.Address.String = address
		arg.Address.Valid = true
	}
	profile, err := testQueries.CreateProfile(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, profile)
	require.Equal(t, account.AccountID, profile.AccountID)
	require.Equal(t, firstName, profile.FirstName)
	require.Equal(t, lastName, profile.LastName)
	require.Equal(t, phone, profile.PhoneNumber)
	require.Equal(t, birthday, profile.Birthday)
	require.Equal(t, address, profile.Address.String)
	require.NotEmpty(t, profile.CreatedAt)
	require.NotEmpty(t, profile.UpdatedAt)
	return
}

func TestCreateProfile(t *testing.T) {
	createTestProfile(t)
}

func TestUpdateProfile(t *testing.T) {
	profile := createTestProfile(t)
	firstName := util.RandomString(6)
	lastName := util.RandomString(8)
	address := util.RandomString(12)
	phone := util.CreateRandomPhone()
	birthday := util.CreateRandomBirthday()
	arg := UpdateProfileParams{
		AccountID:   profile.AccountID,
		PhoneNumber: phone,
		Birthday:    birthday,
		FirstName:   firstName,
		LastName:    lastName,
		Address: pgtype.Text{
			String: address,
			Valid:  true,
		},
	}
	newProfile, err := testQueries.UpdateProfile(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, newProfile)
	require.NotEqual(t, profile.FirstName, newProfile.FirstName)
	require.NotEqual(t, profile.LastName, newProfile.LastName)
	require.NotEqual(t, profile.PhoneNumber, newProfile.PhoneNumber)
	require.NotEqual(t, profile.Address.String, newProfile.Address.String)
	require.Equal(t, profile.CreatedAt, newProfile.CreatedAt)
	require.True(t, newProfile.UpdatedAt.After(profile.UpdatedAt))
}

func TestDeleteProfile(t *testing.T) {
	profile := createTestProfile(t)
	err := testQueries.DeleteProfile(context.Background(), profile.AccountID)
	require.NoError(t, err)
	oldProfile, err := testQueries.GetOneProfile(context.Background(), profile.AccountID)
	require.ErrorIs(t, err, pgx.ErrNoRows)
	require.Empty(t, oldProfile)
}
