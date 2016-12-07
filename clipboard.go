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
		c, _ := clip.ReadAll()
		if len(currentBoard) == len(c) {
			for i := 0; i < len(currentBoard); i++ {
				if currentBoard[i] != c[i] {
					currentBoard = []byte(c)
					buffer := new(bytes.Buffer)
					for _, val := range activeDevices {
						fmt.Println(currentBoard)
						buffer.Write(currentBoard)
						_, err := http.Post("http://"+val.IP+":"+fmt.Sprint(val.Port)+"/send", "application/octet-stream", buffer)
						if err != nil {
							log.Println("Failed to send Clipboard to client:", val.IP)
						}
					}
					break
				}
			}
		} else {
			currentBoard = []byte(c)
			buffer := new(bytes.Buffer)
			for _, val := range activeDevices {
				fmt.Println(currentBoard)
				buffer.Write(currentBoard)
				_, err := http.Post("http://"+val.IP+":"+fmt.Sprint(val.Port)+"/send", "application/octet-stream", buffer)
				if err != nil {
					log.Println("Failed to send Clipboard to client:", val.IP)
				}
			}
		}
		if len(recievedBoard) == len(currentBoard) {
			for i := 0; i < len(recievedBoard); i++ {
				if recievedBoard[i] != currentBoard[i] {
					currentBoard = recievedBoard
					clip.WriteAll(string(recievedBoard))
				}
			}
		} else {
			currentBoard = recievedBoard
			clip.WriteAll(string(recievedBoard))
		}
		time.Sleep(time.Millisecond * 500)
	}
}
