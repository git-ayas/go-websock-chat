package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/goombaio/namegenerator"
)

var host, port = flag.String("host", "192.168.1.14", "Name of host"), flag.Int("port", 3174, "Port to listen on")

func main() {
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
			log.Printf("recv: %s on mt: %d", msg, messageType)

			parsedJsonRequest, err := UnmarshaledJsonFromByteArray[struct {
				Sender  string
				Message string
			}](msg)

			if err != nil {
				log.Println("Failed to parse request data", err)
			}

			log.Printf("Parsed Request: %v %v \n", parsedJsonRequest.Sender, parsedJsonRequest.Message)
			parsedCompomnent, err := GetParsedMessageComponent(MesageDataType{
				Sender:  parsedJsonRequest.Sender,
				Content: parsedJsonRequest.Message,
			})
			if err != nil {
				log.Println("Failed to parse message template", err)
			}
			if err = c.WriteMessage(messageType, parsedCompomnent); err != nil {
				log.Println("write:", err)
				break
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
}
