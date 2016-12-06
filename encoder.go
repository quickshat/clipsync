package main

import (
	"bytes"
	"encoding/gob"
)

var encoderChannel = make(chan packet)

func encoder() {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)
	for {
		p := <-encoderChannel
		enc.Encode(p)
		sendData(buffer.Bytes())
		buffer.Reset()
	}
}
