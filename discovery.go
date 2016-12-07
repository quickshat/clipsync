package main

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"

	"encoding/binary"

	"golang.org/x/net/ipv4"
)

type discoveryServiceEntry struct {
	IP   net.Addr
	Port int64
	New  bool
}

type discoveryService struct {
	Port              int
	Conn              *ipv4.PacketConn
	Group             net.IP
	IgnoreAddress     []string
	AnnouncmentPort   int64
	DiscoveredChannel chan discoveryServiceEntry
	PingPacket        []byte
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
			emitLog("DISSERVICE", "Bind to", ifc)
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

	pp := new(bytes.Buffer)
	binary.Write(pp, binary.LittleEndian, d.AnnouncmentPort)
	pp.WriteByte(0)

	d.PingPacket = pp.Bytes()

	return &d
}

func (d *discoveryService) reciever() {
	b := make([]byte, 9)
	buffer := new(bytes.Buffer)
	port := int64(0)
	new := byte(0)

	for {
		n, _, src, _ := d.Conn.ReadFrom(b)

		own := false
		for i := 0; i < len(d.IgnoreAddress); i++ {
			if fmt.Sprint(src) == d.IgnoreAddress[i] {
				own = true
				break
			}
		}

		if own || n != 9 {
			continue
		}

		buffer.Reset()
		buffer.Write(b)

		err := binary.Read(buffer, binary.LittleEndian, &port)
		if err != nil {
			continue
		}

		new, err = buffer.ReadByte()
		if err != nil {
			continue
		}

		emitLog("DISSERVICE", "Device Ping recieved from", src, "on", port, "new:", new == 1)
		if new == 1 {
			d.ping()
		}
		d.DiscoveredChannel <- discoveryServiceEntry{src, port, new == 1}
	}
}

func (d *discoveryService) ping() {
	emitLog("DISSERVICE", "PING!")
	d.sendData(d.PingPacket)
}

func (d *discoveryService) pinger() {
	for {
		time.Sleep(time.Minute)
		d.ping()
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
	emitLog("DISSERVICE", "ANNOUNCMENT PING!")

	aBytes := new(bytes.Buffer)
	binary.Write(aBytes, binary.LittleEndian, d.AnnouncmentPort)
	aBytes.WriteByte(1)

	go d.reciever()
	go d.pinger()
}
