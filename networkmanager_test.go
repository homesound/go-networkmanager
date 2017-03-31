package network_manager

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestListInterfaces(t *testing.T) {
	require := require.New(t)

	nm := New()
	_, err := nm.ListInterfaces()
	require.Nil(err)
}

func TestIsWifiConnected(t *testing.T) {
	require := require.New(t)

	nm := New()
	v, err := nm.IsWifiConnected()
	require.Nil(err)
	require.True(v)
}

func TestGetWifiInterfaces(t *testing.T) {
	require := require.New(t)

	nm := New()
	v, err := nm.GetWifiInterfaces()
	require.Nil(err)
	require.NotZero(len(v))
}

func TestWifiScan(t *testing.T) {
	require := require.New(t)

	// Should fail
	nm := New()
	v, err := nm.WifiScan("wlan8")
	require.NotNil(err)
	require.Zero(len(v))

	v, err = nm.WifiScan("wlan0")
	require.Nil(err)
	require.NotZero(len(v))
}
