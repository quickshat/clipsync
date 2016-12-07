package main

import (
	"net/http"
	"time"

	"fmt"

	"bytes"

	"log"

	clip "github.com/atotto/clipboard"
)

var recievedBoard []byte
var currentBoard []byte

func detectNewClipboard() {
	for {

		time.Sleep(time.Millisecond * 500)
	}
}

func commit(c string) {
	currentBoard = []byte(c)
	buffer := new(bytes.Buffer)
	for _, val := range activeDevices {
		buffer.Write(currentBoard)
		_, err := http.Post("http://"+val.IP+":"+fmt.Sprint(val.Port)+"/send", "application/octet-stream", buffer)
		if err != nil {
			log.Println("Failed to send Clipboard to client:", val.IP)
		}
	}
}

func getOsClipboard() []byte {
	c, _ := clip.ReadAll()
	return []byte(c)
}
