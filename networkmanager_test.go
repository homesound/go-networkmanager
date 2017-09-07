package networkmanager

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	simpleexec "github.com/gurupras/go-simpleexec"
	"github.com/stretchr/testify/require"
)

func getWifiInterface(nm *NetworkManager) (string, error) {
	wifiInterfaces, err := nm.GetWifiInterfaces()
	if err != nil {
		return "", err
	}
	iface := wifiInterfaces[0]
	return iface, err
}

func TestHostname(t *testing.T) {
	require := require.New(t)

	nm := New()
	hostname, err := nm.Hostname()
	require.Nil(err)
	require.NotEmpty(hostname)
}

func TestIfUp(t *testing.T) {
	require := require.New(t)

	nm := New()
	iface, err := getWifiInterface(nm)
	require.Nil(err)
	err = nm.IfUp(iface)
	require.Nil(err, fmt.Sprintf("Failed: %v", err))

	// Now, check if it is up
	cmd := simpleexec.ParseCmd(fmt.Sprintf("iw dev %v link", iface))
	buf := bytes.NewBuffer(nil)
	cmd.Stdout = buf
	err = cmd.Run()
	require.Nil(err)
	//require.NotEqual("Not connected.\n", buf.String())
}

func TestIfDown(t *testing.T) {
	require := require.New(t)

	nm := New()
	iface, err := getWifiInterface(nm)
	require.Nil(err)
	err = nm.IfDown(iface)
	require.Nil(err)

	// Now, check if it is up
	cmd := simpleexec.ParseCmd(fmt.Sprintf("iw dev %v link", iface))
	buf := bytes.NewBuffer(nil)
	cmd.Stdout = buf
	err = cmd.Run()
	require.Nil(err)
	require.Equal("Not connected.\n", buf.String())

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

	iface, err := getWifiInterface(nm)
	require.Nil(err)
	err = nm.IfUp(iface)
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
	v, err := nm.WifiScan("lo")
	require.NotNil(err)
	require.Zero(len(v))

	iface, err := getWifiInterface(nm)
	require.Nil(err)
	err = nm.IfUp(iface)
	require.Nil(err)

	v, err = nm.WifiScan(iface)
	require.Nil(err)
	require.NotZero(len(v))
}

func TestIPAddress(t *testing.T) {
	require := require.New(t)

	nm := New()
	ip, err := nm.IPAddress("lo")
	require.Nil(err)
	require.NotEqual("", ip)
}
