package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofiber/contrib/websocket"
)

type ObserverEntry struct {
	id         string
	connection *websocket.Conn
}

type MessageData struct {
	sender  string
	message string
}
type MessageObservable struct {
	observers          []ObserverEntry
	MessagesChannel    chan MessageData
	UnsubscribeChannel chan ObserverEntry
	SubscriberChannel  chan ObserverEntry
	terminateObs       chan int
	logger             *log.Logger
	mu                 sync.Mutex
}

/*
Its suppoosed to spell broker,but I like the typo
more than being accurate. So we now have a friendly
Borker who is fetching us messages.
*/
func (m *MessageObservable) Borker() {
	observerCount := len(m.observers)
	ticker := time.Tick(8 * time.Second)
	allowDebouncedStatusLog := true
	for {
		select {
		case message := <-m.MessagesChannel:
			processedOut, err := GetParsedMessageComponent(MesageDataType{
				Sender:  message.sender,
				Content: message.message,
			})
			if err != nil {
				m.logger.Println("Error parsing temp for broadcast.")
			}
			m.mu.Lock()
			// Write out messages to all observers
			for _, observer := range m.observers {
				observer.connection.WriteMessage(1, processedOut)
				m.logger.Printf("Woof! Sending data(%v) to %s  ", message, observer.id)
			}
			m.mu.Unlock()
		case <-m.terminateObs:
			fmt.Println("Terminating messaging observable")
			return
		case subscriber := <-m.SubscriberChannel:
			m.mu.Lock()
			m.observers = append(m.observers, subscriber)
			m.logger.Printf("New subscriber connected 🐕<Perking Ear Floofs>")
			m.mu.Unlock()
		case subscriber := <-m.UnsubscribeChannel:
			m.mu.Lock()
			for i, observer := range m.observers {
				if observer.id == subscriber.id {
					m.observers = append(m.observers[:i], m.observers[i+1:]...)
					m.logger.Printf("Subscriber left <sad borker noises>")
				}
			}
			m.mu.Unlock()
		case <-ticker:
			allowDebouncedStatusLog = true

		default:
			if allowDebouncedStatusLog {
				m.mu.Lock()
				if len(m.observers) > observerCount {
					m.logger.Printf("\n Woof! %d new subscriber/s 🐕 \n", len(m.observers)-observerCount)
					observerCount = len(m.observers)
				}
				if observerCount < 1 {
					m.logger.Printf("No subscribers. Lonely Borker :(")

				} else {
					m.logger.Printf("%d subscribers connected 🐕<tail wags>", observerCount)
				}
				allowDebouncedStatusLog = false
				m.mu.Unlock()
			}
		}
	}

}

/*
Some of the below are not used in this example, but it is a good example
of how to use the observer pattern and are meant to be usage examples
for the channels as well.

They can be used to parse inputs before processing.
It is also deferable allowing you to use it in a cleanup cycle.
*/
func (m *MessageObservable) Produce(message MessageData) {
	m.MessagesChannel <- message

}

func (m *MessageObservable) Subscribe(subscriber ObserverEntry) {
	m.logger.Printf("Subscribing %s", subscriber.id)
	m.SubscriberChannel <- subscriber

}
func (m *MessageObservable) Unsubscribe(subscriber ObserverEntry) {

	m.UnsubscribeChannel <- subscriber

}
func (m *MessageObservable) Terminate() {
	m.terminateObs <- 1
}

func NewMessageObservable(logger *log.Logger) *MessageObservable {
	return &MessageObservable{
		observers:          []ObserverEntry{},
		MessagesChannel:    make(chan MessageData),
		UnsubscribeChannel: make(chan ObserverEntry),
		SubscriberChannel:  make(chan ObserverEntry),
		terminateObs:       make(chan int),
		logger:             logger,
	}
}
