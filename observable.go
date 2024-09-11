package main

import (
	"fmt"
	"log"

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
	observers    []ObserverEntry
	messageChan  chan MessageData
	unsubChan    chan ObserverEntry
	terminateObs chan int
}

// Its suppoosed to spell broker,but I like the typo more than being accurate.
func (m *MessageObservable) Borker() {
	observerCount := len(m.observers)
	for {
		select {
		case message := <-m.messageChan:
			processedOut, err := GetParsedMessageComponent(MesageDataType{
				Sender:  message.sender,
				Content: message.message,
			})
			if err != nil {
				log.Println("Error parsing temp for broadcast.")
			}
			// Write out messages to all observers
			for _, observer := range m.observers {
				observer.connection.WriteMessage(1, processedOut)
			}
		case <-m.terminateObs:
			fmt.Println("Terminating messaging observable")
			return
		default:
			if len(m.observers) > observerCount {
				log.Printf("\n🐶: Woof! %d new subscriber/s 🐕 \n", len(m.observers)-observerCount)
				observerCount = len(m.observers)
			}
			if observerCount < 1 {
				log.Println("🐶: No subscribers. Doggo sad :(")

			} else {
				log.Printf("🐶: %d subscribers connected 🐕<wags>", observerCount)
			}

		}
	}

}

func (m *MessageObservable) Produce(message MessageData) {

}

func (m *MessageObservable) Subscribe(subscriber ObserverEntry) {

}
func (m *MessageObservable) Unsubscribe(subscriber ObserverEntry) {

}
