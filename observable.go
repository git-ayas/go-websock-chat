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

// Its suppoosed to spell broker,but I like the typo more than being accurate.
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
			}
			m.mu.Unlock()
		case <-m.terminateObs:
			fmt.Println("Terminating messaging observable")
			return
		case subscriber := <-m.SubscriberChannel:
			m.mu.Lock()
			m.observers = append(m.observers, subscriber)
			m.logger.Printf("🐶: New subscriber connected 🐕<wags!>")
			m.mu.Unlock()
		case subscriber := <-m.UnsubscribeChannel:
			m.mu.Lock()
			for i, observer := range m.observers {
				if observer.id == subscriber.id {
					m.observers = append(m.observers[:i], m.observers[i+1:]...)
					m.logger.Printf("🐶: Subscriber disconnected 🐕<whimpers>")
				}
			}
			m.mu.Unlock()
		case <-ticker:
			allowDebouncedStatusLog = true

		default:
			if allowDebouncedStatusLog {
				m.mu.Lock()
				if len(m.observers) > observerCount {
					m.logger.Printf("\n🐶: Woof! %d new subscriber/s 🐕 \n", len(m.observers)-observerCount)
					observerCount = len(m.observers)
				}
				if observerCount < 1 {
					m.logger.Printf("🐶: No subscribers. Doggo sad :( %v users", len(m.observers))

				} else {
					m.logger.Printf("🐶: %d subscribers connected 🐕<wags>", observerCount)
				}
				allowDebouncedStatusLog = false
				m.mu.Unlock()
			}
		}
	}

}

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
