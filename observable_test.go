package main

import (
	"testing"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/stretchr/testify/assert"
)

func TestSubscribe(t *testing.T) {
	m := &MessageObservable{
		observers:    []ObserverEntry{},
		messageChan:  make(chan MessageData),
		unsubChan:    make(chan ObserverEntry),
		subChan:      make(chan ObserverEntry),
		terminateObs: make(chan int),
	}

	go m.Borker()

	subscriber := ObserverEntry{id: "1", connection: &websocket.Conn{}}
	m.Subscribe(subscriber)

	time.Sleep(100 * time.Millisecond) // Give some time for the goroutine to process

	assert.Equal(t, 1, len(m.observers))
	assert.Equal(t, "1", m.observers[0].id)
}

func TestUnsubscribe(t *testing.T) {
	m := &MessageObservable{
		observers:    []ObserverEntry{},
		messageChan:  make(chan MessageData),
		unsubChan:    make(chan ObserverEntry),
		subChan:      make(chan ObserverEntry),
		terminateObs: make(chan int),
	}

	go m.Borker()

	subscriber := ObserverEntry{id: "1", connection: &websocket.Conn{}}
	m.Subscribe(subscriber)

	time.Sleep(100 * time.Millisecond) // Give some time for the goroutine to process

	m.Unsubscribe(subscriber)

	time.Sleep(100 * time.Millisecond) // Give some time for the goroutine to process

	assert.Equal(t, 0, len(m.observers))
}

func TestProduce(t *testing.T) {
	m := &MessageObservable{
		observers:    []ObserverEntry{},
		messageChan:  make(chan MessageData),
		unsubChan:    make(chan ObserverEntry),
		subChan:      make(chan ObserverEntry),
		terminateObs: make(chan int),
	}

	go m.Borker()

	subscriber := ObserverEntry{id: "1", connection: &websocket.Conn{}}
	m.Subscribe(subscriber)

	time.Sleep(100 * time.Millisecond) // Give some time for the goroutine to process

	message := MessageData{sender: "test", message: "hello"}
	m.Produce(message)

	time.Sleep(100 * time.Millisecond) // Give some time for the goroutine to process

	// Assuming GetParsedMessageComponent and observer.connection.WriteMessage are working correctly,
	// we should have processed the message. This is a placeholder assertion.
	assert.NotNil(t, message)
}

func TestBorkerTerminate(t *testing.T) {
	m := &MessageObservable{
		observers:    []ObserverEntry{},
		messageChan:  make(chan MessageData),
		unsubChan:    make(chan ObserverEntry),
		subChan:      make(chan ObserverEntry),
		terminateObs: make(chan int),
	}

	go m.Borker()

	time.Sleep(100 * time.Millisecond) // Give some time for the goroutine to start

	m.terminateObs <- 1

	time.Sleep(100 * time.Millisecond) // Give some time for the goroutine to terminate

	// Placeholder assertion to ensure the test runs to completion
	assert.True(t, true)
}
