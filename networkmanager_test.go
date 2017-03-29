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
