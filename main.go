package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
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
	resp, err := http.Get("http://" + a.IP + ":" + fmt.Sprint(a.Port) + "/metrics")
	if err != nil {
		return err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	json.Unmarshal(b, &a.Metrics)
	fmt.Println(a.Metrics)
	return nil
}

func (a *activeDevice) register() error {
	log.Println("[MAIN] Manual register")
	r, err := http.PostForm(a.IP+":"+fmt.Sprint(a.Port)+"/register", url.Values{
		"port": {"8081"},
	})
	if err != nil || r.StatusCode != http.StatusOK {
		return err
	}
	return nil
}

var activeDevices map[string]*activeDevice

func main() {
	activeDevices = make(map[string]*activeDevice)
	disService = createDiscoveryService(1024, "en0", 8081)
	disService.start()
	activeManager()
	initWeb()
	go detectNewClipboard()
}

func clearManager() {
	for key, val := range activeDevices {
		if val.LastPing.Sub(time.Now()).Minutes() > 3 {
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
			split := strings.Split(fmt.Sprint(d.IP), ":")
			ip := split[0]

			val, found := activeDevices[ip]
			if found {
				val.LastPing = time.Now()
				val.Port = d.Port
			} else {
				a := &activeDevice{ip, d.Port, time.Now(), metrics{}}
				err := a.loadMetrics()
				if err == nil {
					a.register()
					activeDevices[ip] = a
				}
			}
		}
	}()
}
