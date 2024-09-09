package main

import "github.com/gofiber/contrib/websocket"

type MessageObservable struct {
	observers []*websocket.Conn
}

func (m *MessageObservable) Consume() {

}

func (m *MessageObservable) Produce() {

}

func (m *MessageObservable) Subscribe() {

}
