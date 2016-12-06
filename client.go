package main

import (
	"fmt"
	"net"
	"strings"

	"golang.org/x/net/ipv4"
)

var globalPort int
var globalConn *ipv4.PacketConn
var globalGroup net.IP
var ignoreAddress []string

func initServer(port int, ifc string) {
	globalPort = port
	globalGroup = net.IPv4(224, 0, 0, 250)

	is, _ := net.Interfaces()

	c, _ := net.ListenPacket("udp4", "0.0.0.0:"+fmt.Sprint(globalPort))
	globalConn = ipv4.NewPacketConn(c)

	for i := 0; i < len(is); i++ {
		addr, _ := is[i].Addrs()
		for i := 0; i < len(addr); i++ {
			if strings.Contains(addr[i].String(), ":") {
				continue
			}
			strs := strings.Split(addr[i].String(), "/")
			if len(strs) == 2 {
				if ifc == is[i].Name {
					globalConn.SetMulticastInterface(&is[i])
					ignoreAddress = append(ignoreAddress, strs[0]+":"+fmt.Sprint(globalPort))
					globalConn.JoinGroup(&is[i], &net.UDPAddr{IP: globalGroup})
					fmt.Println("Interface Found and binded!")
				}
			}
		}
	}

	globalConn.SetMulticastTTL(2)

	b := make([]byte, 100)
	for {
		n, _, src, _ := globalConn.ReadFrom(b)

		own := false
		for i := 0; i < len(ignoreAddress); i++ {
			if fmt.Sprint(src) == ignoreAddress[i] {
				own = true
				break
			}
		}

		if own {
			continue
		}

		fmt.Println(b, n, src)
	}
}

func sendData(data []byte) error {
	dst := &net.UDPAddr{IP: globalGroup, Port: globalPort}
	if _, err := globalConn.WriteTo(data, nil, dst); err != nil {
		return err
	}
	return nil
}
