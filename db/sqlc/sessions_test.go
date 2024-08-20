package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tuda4/mb-backend/util"
)

const (
	durationTestRefreshToken = 24
)

func createTestSession(t *testing.T) (session Session) {
	account := createTestAccount(t)
	refreshToken := util.RandomString(32)
	userAgent := util.RandomString(16)
	clientID := util.RandomString(8)
	arg := CreateSessionParams{
		AccountID:    account.AccountID,
		RefreshToken: refreshToken,
		UserAgent:    userAgent,
		ClientID:     clientID,
		IsBlocked:    false,
		ExpiredAt:    time.Now().Add(time.Hour * durationTestRefreshToken),
	}
	session, err := testQueries.CreateSession(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, session)
	require.Equal(t, session.AccountID, account.AccountID)
	require.Equal(t, session.RefreshToken, refreshToken)
	require.Equal(t, session.UserAgent, userAgent)
	require.Equal(t, session.ClientID, clientID)
	require.Equal(t, session.IsBlocked, false)
	require.WithinDuration(t, session.ExpiredAt, arg.ExpiredAt, time.Microsecond)
	require.NotEmpty(t, session.CreatedAt)

	return
}

func TestCreateSession(t *testing.T) {
	createTestSession(t)
}

func TestGetOneSession(t *testing.T) {
	session := createTestSession(t)
	currentSession, err := testQueries.GetOneSession(context.Background(), session.RefreshToken)
	require.NoError(t, err)
	require.NotEmpty(t, currentSession)
	require.Equal(t, session.AccountID, currentSession.AccountID)
	require.Equal(t, session.RefreshToken, currentSession.RefreshToken)
	require.Equal(t, session.UserAgent, currentSession.UserAgent)
	require.Equal(t, session.ClientID, currentSession.ClientID)
	require.Equal(t, session.IsBlocked, currentSession.IsBlocked)
	require.Equal(t, session.ExpiredAt, currentSession.ExpiredAt)
	require.Equal(t, session.CreatedAt, currentSession.CreatedAt)
}
