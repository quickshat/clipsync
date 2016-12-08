package main

import (
	"net/http"
	"time"

	"fmt"

	"bytes"

	"log"

	"sync"

	clip "github.com/atotto/clipboard"
)

var recievedBoard []byte
var currentBoard []byte

var recievedBoardLock sync.Mutex

func detectNewClipboard() {
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

func setOsClipboard(b []byte) {
	clip.WriteAll(string(b))
}
