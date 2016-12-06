package main

import (
	"fmt"
	"net"

	"golang.org/x/net/ipv4"
)

var globalPort int
var globalConn *ipv4.PacketConn
var globalGroup net.IP
var networkInterfaces []net.Interface

func initServer(port int) {
	globalPort = port
	globalGroup = net.IPv4(224, 0, 0, 250)

	is, _ := net.Interfaces()
	networkInterfaces = is

	c, _ := net.ListenPacket("udp4", "0.0.0.0:"+fmt.Sprint(globalPort))
	globalConn = ipv4.NewPacketConn(c)

	for i := 0; i < len(is); i++ {
		globalConn.JoinGroup(&is[i], &net.UDPAddr{IP: globalGroup})
	}

	b := make([]byte, 1500)
	for {
		n, _, src, _ := globalConn.ReadFrom(b)
		fmt.Println(b, n, src)
	}
}

func sendData(data []byte) error {
	dst := &net.UDPAddr{IP: globalGroup, Port: globalPort}
	for i := 0; i < len(networkInterfaces); i++ {
		if err := globalConn.SetMulticastInterface(&networkInterfaces[i]); err != nil {
			return err
		}
		globalConn.SetMulticastTTL(2)
		if _, err := globalConn.WriteTo(data, nil, dst); err != nil {
			return err
		}
	}
	return nil
}
