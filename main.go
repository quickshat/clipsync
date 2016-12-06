package main

import (
	"time"
)

var disService *discoveryService

func main() {
	disService = createDiscoveryService(1024, "en0", 6666)
	disService.start()
	for {
		time.Sleep(time.Hour)
	}
}
