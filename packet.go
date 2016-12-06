package main

type packet struct {
	ID      byte
	Group   []string
	Payload interface{}
}

type clipboardPayload []byte
