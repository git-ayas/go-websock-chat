package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/goombaio/namegenerator"
)

var host, port = flag.String("host", "127.0.0.1", "Host name to listen on"), flag.Int("port", 3174, "Port to listen on")

func main() {
	flag.Parse()
	var logger *log.Logger
	var logfile, logfileerr = os.OpenFile("borker.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if logfileerr != nil {
		log.Println("Error opening observable.log file:", logfileerr)
		logger = log.New(os.Stdout, "[🐶: stdout]", log.LstdFlags)
	} else {
		defer logfile.Close()
		logger = log.New(logfile, "[🐶:]", log.LstdFlags)
	}
	var ChatroomPool = NewMessageObservable(logger)
	go ChatroomPool.Borker()
	app := fiber.New(fiber.Config{
		Views: html.New("./views", ".html"),
	})

	app.Use("/relay", func(c *fiber.Ctx) error {

		if websocket.IsWebSocketUpgrade(c) {

			c.Locals("allowed", true)
			return c.Next()

		}

		return fiber.ErrUpgradeRequired

	})

	app.Use("/relay/chatroom", websocket.New(func(c *websocket.Conn) {
		timeOfConnection := time.Now().String()
		ChatroomPool.SubscriberChannel <- ObserverEntry{id: timeOfConnection, connection: c}
		defer ChatroomPool.Unsubscribe(ObserverEntry{id: timeOfConnection, connection: c})

		var (
			messageType int
			msg         []byte
			err         error
		)
		for {
			if messageType, msg, err = c.ReadMessage(); err != nil {
				log.Println("read:", err)
				break
			}
			log.Printf("recv: %s on mt: %d, observerCount %v", msg, messageType, len(ChatroomPool.observers))

			parsedJsonRequest, err := UnmarshaledJsonFromByteArray[struct {
				Sender  string
				Message string
			}](msg)

			if err != nil {
				log.Println("Failed to parse request data", err)
			}

			log.Printf("Parsed Request on chatroom2: %v %v \n", parsedJsonRequest.Sender, parsedJsonRequest.Message)

			ChatroomPool.MessagesChannel <- MessageData{
				sender:  parsedJsonRequest.Sender,
				message: parsedJsonRequest.Message,
			}

		}

	}))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("pages/index", fiber.Map{
			"Title":  "Home: Le Chat Room",
			"Sender": namegenerator.NewNameGenerator(time.Now().UTC().UnixNano()).Generate(),
		}, "layouts/main")
	})

	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", *host, *port)))
	ChatroomPool.Terminate()
}
