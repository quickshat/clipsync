package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var disService *discoveryService

type activeDevice struct {
	IP       string
	Port     int64
	LastPing time.Time
	Metrics  metrics
}

func (a *activeDevice) loadMetrics() error {
	resp, err := http.Get(fmt.Sprint(a.IP) + ":" + fmt.Sprint(a.Port) + "/metrics")
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(b, &a.Metrics)
	return nil
}

var activeDevices map[string]*activeDevice

func main() {
	disService = createDiscoveryService(1024, "Ethernet", 8081)
	disService.start()
	activeManager()
	initWeb()
	go detectNewClipboard()
}

func clearManager() {
	for key, val := range activeDevices {
		if time.Now().Sub(val.LastPing).Minutes() > 3 {
			delete(activeDevices, key)
			log.Println("[DISSERVICE] device", key, "inactive")
		}
	}
}

func activeManager() {
	go func() {
		for {
			clearManager()
			time.Sleep(time.Minute)
		}
	}()
	go func() {
		for {
			d := <-disService.DiscoveredChannel
			val, found := activeDevices[fmt.Sprint(d.IP)]
			if found {
				val.LastPing = time.Now()
				val.Port = d.Port
			} else {
				a := &activeDevice{fmt.Sprint(d.IP), d.Port, time.Now(), metrics{}}
				err := a.loadMetrics()
				if err == nil {
					activeDevices[fmt.Sprint(d.IP)] = a
				}
			}
		}
	}()
}
