package networkmanager

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"strings"

	"github.com/google/shlex"
	"github.com/gurupras/go-simpleexec"
	"github.com/gurupras/go-wireless/iwlib"
	"github.com/gurupras/gocommons"
	"github.com/sirupsen/logrus"
)

type WifiScanResult struct {
	*iwlib.WirelessScanResult
}

type NetworkManager struct {
}

func New() *NetworkManager {
	return &NetworkManager{}
}

func (nm *NetworkManager) Hostname() (string, error) {
	data, err := ioutil.ReadFile("/etc/hostname")
	return strings.TrimSpace(string(data)), err
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

// Get a list of Wi-Fi interfaces
// This function works in a very complicated manner
// We first get a list of available interfaces via ListInterfaces
// and then ask iwgetid to list out the channel for each of them
// Non-wireless interfaces throw errors whereas wireless interfaces
// return data
func (nm *NetworkManager) GetWifiInterfaces() ([]string, error) {
	ifaces, err := nm.ListInterfaces()
	if err != nil {
		return nil, err
	}
	wifiInterfaces := make([]string, 0)
	for _, iface := range ifaces {
		cmd := simpleexec.ParseCmd(fmt.Sprintf("/sbin/iwgetid -c %v", iface))
		buf := bytes.NewBuffer(nil)
		cmd.Stdout = buf
		if err := cmd.Run(); err != nil {
			// XXX: This is assumed to be due to 'iface' not being a Wi-Fi interface
			//return nil, fmt.Errorf("Failed to query interface '%v': %v", iface, err)
			logrus.Debugf("Failed to query interface '%v': %v", iface, err)
		} else {
			if buf.Len() != 0 && strings.Contains(buf.String(), "Channel") {
				wifiInterfaces = append(wifiInterfaces, iface)
			}
		}
	}
	return wifiInterfaces, nil
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

func (nm *NetworkManager) IPAddress(iface string) (string, error) {
	// /sbin/ifconfig eth0 | grep 'inet addr:' | cut -d: -f2 | awk '{ print $1}'
	cmd := simpleexec.ParseCmd(fmt.Sprintf("/sbin/ifconfig %v", iface)).Pipe("grep 'inet addr:'").Pipe("cut -d: -f2").Pipe("awk '{ print $1}'")

	buf := bytes.NewBuffer(nil)
	cmd.Stdout = buf
	cmd.Start()
	cmd.Wait()

	return strings.TrimSpace(buf.String()), nil
}
