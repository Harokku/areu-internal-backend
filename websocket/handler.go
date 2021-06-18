package websocket

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"log"
	"time"
)

func DocsUpdate() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		// When the function returns, unregister the client and close the connection
		defer func() {
			unregister <- c
			c.Close()
		}()

		// Register the client
		register <- c

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("%s - DocsUpdate WS read error: %s", time.Now().Format(time.ANSIC), err)
				}
				return // Call defer and close the connection
			}

			//TODO: Remove in production DEBUG only
			if messageType == websocket.TextMessage {
				// Broadcast the received message
				Broadcast <- string(message)
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}
	})
}
