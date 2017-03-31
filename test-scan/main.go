package main

import (
	"fmt"
	"time"

	"github.com/homesound/go-networkmanager"
	log "github.com/sirupsen/logrus"
)

func main() {
	nm := network_manager.New()

	wifiIfaces, err := nm.GetWifiInterfaces()
	if err != nil {
		log.Fatalf("Failed to get wifi interfaces: %v\n", err)
	}

	for {
		for _, iface := range wifiIfaces {
			scanResults, err := nm.WifiScan(iface)
			if err != nil {
				log.Errorf("Failed to get scan results from '%v': %v\n", iface, err)
			}
			fmt.Printf("%v: %v\n", iface, scanResults)
		}
		time.Sleep(3 * time.Second)
	}
}
