package network_manager

import (
	"errors"
	"net"
	"strings"

	"github.com/theojulienne/go-wireless/iwlib"
)

type NetworkManager struct {
}

func New() *NetworkManager {
	return &NetworkManager{}
}

func (nm *NetworkManager) ListInterfaces() ([]string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	ret := make([]string, len(ifaces))
	for idx, iface := range ifaces {
		ret[idx] = iface.Name
	}
	return ret, nil
}

func (nm *NetworkManager) GetWifiInterfaces() ([]string, error) {
	ret := make([]string, 0)
	ifaces, err := nm.ListInterfaces()
	if err != nil {
		return nil, err
	}
	for _, name := range ifaces {
		if strings.Index(name, "wlan") == 0 {
			ret = append(ret, name)
		}
	}
	return ret, nil
}

func (nm *NetworkManager) IsWifiConnected() (bool, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return false, err
	}
	found := false
	for _, iface := range ifaces {
		if strings.Contains(iface.Name, "wlan") {
			found = true
			addrs, err := iface.Addrs()
			if err != nil {
				return false, err
			}
			if len(addrs) > 0 {
				return true, nil
			}
		}
	}
	if !found {
		return false, errors.New("Did not find a wlan interface")
	}
	return false, nil
}

func (nm *NetworkManager) WifiScan(iface string) ([]string, error) {

	networks, err := iwlib.GetWirelessNetworks(iface)
	if err != nil {
		return nil, err
	}
	ret := make([]string, 0)
	for _, network := range networks {
		ret = append(ret, network.SSID)
	}
	return ret, nil
}
