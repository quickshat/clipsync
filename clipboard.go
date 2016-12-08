package main

import (
	"net/http"
	"time"

	"fmt"

	"bytes"

	"sync"

	clip "github.com/atotto/clipboard"
)

var recievedBoard []byte
var currentBoard []byte

var recievedBoardLock sync.Mutex

func detectNewClipboard() {
	currentBoard = getOsClipboard()
	for {
		if bytes.Compare(recievedBoard, getOsClipboard()) != 0 && recievedBoard != nil {
			currentBoard = recievedBoard
			setOsClipboard(currentBoard)

			recievedBoardLock.Lock()
			recievedBoard = nil
			recievedBoardLock.Unlock()

		} else if bytes.Compare(getOsClipboard(), currentBoard) != 0 {
			currentBoard = getOsClipboard()
			commit(currentBoard)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func commit(c []byte) {
	buffer := new(bytes.Buffer)
	for _, val := range activeDevices {
		buffer.Write(c)
		http.Post("http://"+val.IP+":"+fmt.Sprint(val.Port)+"/send", "application/octet-stream", buffer)
	}
}

func getOsClipboard() []byte {
	c, _ := clip.ReadAll()
	return []byte(c)
}

func setOsClipboard(b []byte) {
	clip.WriteAll(string(b))
}
