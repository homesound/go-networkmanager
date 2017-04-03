package network_manager

import (
	"testing"
	"time"

	"github.com/google/shlex"
	"github.com/gurupras/gocommons"
	"github.com/stretchr/testify/require"
)

func TestIfUp(t *testing.T) {
	require := require.New(t)

	nm := New()
	err := nm.IfUp("wlan0")
	require.Nil(err)

	// Now, check if it is up
	cmd, _ := shlex.Split("iw dev wlan0 link")
	ret, stdout, stderr := gocommons.Execv(cmd[0], cmd[1:], true)
	_ = stderr
	require.Zero(ret)
	require.NotEqual("Not connected.\n", stdout)
}

func TestIfDown(t *testing.T) {
	require := require.New(t)

	nm := New()
	err := nm.IfDown("wlan0")
	require.Nil(err)

	// Now, check if it is up
	cmd, _ := shlex.Split("iw dev wlan0 link")
	ret, stdout, stderr := gocommons.Execv(cmd[0], cmd[1:], true)
	_ = stderr
	require.Zero(ret)
	require.Equal("Not connected.\n", stdout)

	// Bring it back up
	nm.IfUp("wlan0")
	time.Sleep(3 * time.Second)
}
func TestListInterfaces(t *testing.T) {
	require := require.New(t)

	nm := New()
	_, err := nm.ListInterfaces()
	require.Nil(err)
}

func TestIsWifiConnected(t *testing.T) {
	require := require.New(t)

	nm := New()

	err := nm.IfUp("wlan0")
	require.Nil(err)

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

	err = nm.IfUp("wlan0")
	require.Nil(err)

	v, err = nm.WifiScan("wlan0")
	require.Nil(err)
	require.NotZero(len(v))
}
