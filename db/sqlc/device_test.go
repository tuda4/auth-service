package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tuda4/mb-backend/util"
)

func createTestDevice(t *testing.T) (device Device) {
	account := createTestAccount(t)
	deviceID := util.RandomString(10)
	durationTestRefreshToken := time.Hour * 24 * 180 // Assuming durationTestRefreshToken is the number of days
	device, err := testQueries.CreateDevice(context.Background(), CreateDeviceParams{
		AccountID:      account.AccountID,
		DeviceID:       deviceID,
		ExpTokenDevice: time.Now().Add(durationTestRefreshToken),
	})
	require.NoError(t, err)
	require.NotEmpty(t, device)
	require.Equal(t, device.AccountID, account.AccountID)
	require.Equal(t, device.DeviceID, deviceID)
	require.WithinDuration(t, device.ExpTokenDevice, time.Now().Add(durationTestRefreshToken), time.Second)
	require.NotEmpty(t, device.CreatedAt)

	return
}

func TestCreateDevice(t *testing.T) {
	createTestDevice(t)
}

func TestGetOneDevice(t *testing.T) {
	device := createTestDevice(t)
	currentDevice, err := testQueries.GetOneDevice(context.Background(), GetOneDeviceParams{
		AccountID: device.AccountID,
		DeviceID:  device.DeviceID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, currentDevice)
	require.Equal(t, device.AccountID, currentDevice.AccountID)
	require.Equal(t, device.DeviceID, currentDevice.DeviceID)
	require.Equal(t, device.ExpTokenDevice, currentDevice.ExpTokenDevice)
	require.Equal(t, device.CreatedAt, currentDevice.CreatedAt)
}
