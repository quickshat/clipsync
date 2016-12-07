package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var settings = struct {
	Interface string
	DisPort   int
	WebPort   int
	Group     string
}{"en0", 1024, 8081, ""}

var disService *discoveryService

var activeDevices map[string]*activeDevice

func main() {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(f)
		defer f.Close()
	}

	log.Println("== CLIP SYNC INSTANCE STARTED ==")

	b, err := ioutil.ReadFile("settings.json")
	if err != nil {
		saveSettings()
	} else {
		json.Unmarshal(b, &settings)
	}

	activeDevices = make(map[string]*activeDevice)

	disService = createDiscoveryService(settings.DisPort, settings.Interface, int64(settings.WebPort))
	disService.start()

	activeManager()
	go detectNewClipboard()
	initWeb()
}

func saveSettings() {
	b, _ := json.MarshalIndent(settings, "", "	")
	ioutil.WriteFile("settings.json", b, 0666)
}

func clearManager() {
	for key, val := range activeDevices {
		if val.LastPing.Sub(time.Now()).Minutes() > 3 {
			delete(activeDevices, key)
			emitLog("DISSERVICE", key, "inactive")
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
					activeDevices[ip] = a
				}
			}
		}
	}()
}
