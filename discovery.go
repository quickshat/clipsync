package main

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"

	"encoding/binary"

	"log"

	"golang.org/x/net/ipv4"
)

type discoveryServiceEntry struct {
	IP   net.Addr
	Port int64
}

type discoveryService struct {
	Port              int
	Conn              *ipv4.PacketConn
	Group             net.IP
	IgnoreAddress     []string
	AnnouncmentPort   int64
	DiscoveredChannel chan discoveryServiceEntry
}

func createDiscoveryService(port int, ifc string, announcmentPort int64) *discoveryService {
	d := discoveryService{}

	d.DiscoveredChannel = make(chan discoveryServiceEntry, 100)
	d.AnnouncmentPort = announcmentPort
	d.Port = port
	d.Group = net.IPv4(224, 0, 0, 250)

	is, _ := net.Interfaces()

	c, _ := net.ListenPacket("udp4", "0.0.0.0:"+fmt.Sprint(d.Port))
	d.Conn = ipv4.NewPacketConn(c)

	for i := 0; i < len(is); i++ {
		if ifc == is[i].Name {
			log.Println("[DISSERVICE] bind to", ifc)
			d.Conn.SetMulticastInterface(&is[i])
			d.Conn.JoinGroup(&is[i], &net.UDPAddr{IP: d.Group})
		} else {
			continue
		}

		addr, _ := is[i].Addrs()
		for j := 0; j < len(addr); j++ {
			if strings.Contains(addr[j].String(), ":") {
				continue
			}
			strs := strings.Split(addr[j].String(), "/")
			if len(strs) == 2 {
				d.IgnoreAddress = append(d.IgnoreAddress, strs[0]+":"+fmt.Sprint(d.Port))
			}
		}
	}

	d.Conn.SetMulticastTTL(2)

	return &d
}

func (d *discoveryService) reciever() {
	b := make([]byte, 8)
	buffer := new(bytes.Buffer)
	port := int64(0)

	for {
		_, _, src, _ := d.Conn.ReadFrom(b)

		own := false
		for i := 0; i < len(d.IgnoreAddress); i++ {
			if fmt.Sprint(src) == d.IgnoreAddress[i] {
				own = true
				break
			}
		}

		if own {
			continue
		}

		buffer.Reset()
		buffer.Write(b)
		err := binary.Read(buffer, binary.LittleEndian, &port)
		if err != nil {
			continue
		}

		log.Println("[DISSERVICE]", src, port)
		d.DiscoveredChannel <- discoveryServiceEntry{src, port}
	}
}

func (d *discoveryService) pinger() {
	aBytes := new(bytes.Buffer)
	binary.Write(aBytes, binary.LittleEndian, d.AnnouncmentPort)
	for {
		log.Println("[DISSERVICE] PING!")
		d.sendData(aBytes.Bytes())
		time.Sleep(time.Minute)
	}
}

func (d *discoveryService) sendData(data []byte) error {
	dst := &net.UDPAddr{IP: d.Group, Port: d.Port}
	if _, err := d.Conn.WriteTo(data, nil, dst); err != nil {
		return err
	}
	return nil
}

func (d *discoveryService) start() {
	go d.reciever()
	go d.pinger()
}
