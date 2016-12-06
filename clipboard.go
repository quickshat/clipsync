package main

import (
	"bytes"
	"net/http"
	"time"

	clip "github.com/atotto/clipboard"
)

var recievedBoard []byte
var currentBoard []byte

func detectNewClipboard() {
	for {
		c, _ := clip.ReadAll()
		if len(currentBoard) == len(c) {
			for i := 0; i < len(recievedBoard); i++ {
				if currentBoard[i] != c[i] {
					currentBoard = []byte(c)
					resp, err := http.Post("", "application/octet-stream", bytes.Buffer())
					break
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
		}
		time.Sleep(time.Millisecond * 500)
	}
}
