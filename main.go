package main

import (
	"fmt"
)

func main() {
	go encoder()
	go decoder()
	go initServer(1024, "Ethernet")
	for {
		p := <-recievedPacketChannel
		fmt.Println(p)
	}
}
