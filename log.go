package main

import "log"

func emitLog(t string, message ...interface{}) {
	log.Println("["+t+"]", message)
}
