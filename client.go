package main

import (
	"fmt"
	"net"

	"golang.org/x/net/ipv4"
)

var globalConn *ipv4.PacketConn

func initServer() {
	is, _ := net.Interfaces()
	group := net.IPv4(224, 0, 0, 250)
	c, _ := net.ListenPacket("udp4", "0.0.0.0:1024")
	globalConn = ipv4.NewPacketConn(c)

	for i := 0; i < len(is); i++ {
		globalConn.JoinGroup(&is[i], &net.UDPAddr{IP: group})
	}

	globalConn.SetControlMessage(ipv4.FlagDst, true)

	b := make([]byte, 1500)
	for {
		n, cm, src, _ := globalConn.ReadFrom(b)
		if cm.Dst.IsMulticast() {
			if cm.Dst.Equal(group) {
				fmt.Println(n, src)
			} else {
				continue
			}
		}
	}
}
