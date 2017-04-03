package network_manager

import (
	"fmt"
	"net"
	"strings"

	"github.com/google/shlex"
	"github.com/gurupras/go-wireless/iwlib"
	"github.com/gurupras/gocommons"
)

type WifiScanResult struct {
	*iwlib.WirelessScanResult
}

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
		return false, fmt.Errorf("Did not find a wlan interface")
	}
	return false, nil
}

func (nm *NetworkManager) WifiScan(iface string) ([]*WifiScanResult, error) {
	scanResults, err := iwlib.GetWirelessNetworks(iface)
	if err != nil {
		return nil, err
	}
	results := make([]*WifiScanResult, len(scanResults))
	for idx, res := range scanResults {
		results[idx] = &WifiScanResult{res}
	}
	return results, nil
}

func (nm *NetworkManager) IfUp(iface string) error {
	return nm.ifconfig(iface, "up")
}

func (nm *NetworkManager) IfDown(iface string) error {
	return nm.ifconfig(iface, "down")
}

func (nm *NetworkManager) ifconfig(iface string, state string) error {
	commandStr := fmt.Sprintf("ifconfig %v %v", iface, state)
	cmdline, err := shlex.Split(commandStr)
	if err != nil {
		return fmt.Errorf("Failed to split commandline: '%v': %v", commandStr, err)
	}
	ret, stdout, stderr := gocommons.Execv(cmdline[0], cmdline[1:], true)
	_ = stdout
	if ret != 0 {
		return fmt.Errorf(stderr)
	}
	return nil
}
