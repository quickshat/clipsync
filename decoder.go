package main

import (
	"bytes"
	"encoding/gob"
)

var decoderStream bytes.Buffer
var recievedPacketChannel = make(chan packet)

func decoder() {
	dec := gob.NewDecoder(&decoderStream)

	for {
		var p packet
		dec.Decode(&p)
		if p.Payload == nil {
			continue
		}
		recievedPacketChannel <- p
	}
}
