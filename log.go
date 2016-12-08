package main

import (
	"io"
	"log"
)

type logSplitter struct {
	Output io.Writer
	Logs   []string
}

func (l *logSplitter) Write(p []byte) (int, error) {
	s := string(p)
	l.Logs = append(l.Logs, s[:len(s)-1])
	if l.Output == nil {
		return len(p), nil
	}
	return l.Output.Write(p)
}

func emitLog(t string, message ...interface{}) {
	log.Print("["+t+"]", message)
}
