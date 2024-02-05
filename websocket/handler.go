package websocket

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
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

		// Listen for message and respond accordingly (pong incoming ping req)
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logErr(c, err, "DocsUpdate WS read error")
				} else {
					logErr(c, err, "DocsUpdate WS error")
				}
				return // Call defer and close the connection
			}

			if mt == websocket.PingMessage {
				err := c.WriteMessage(websocket.PongMessage, []byte("Pong"))
				if err != nil {
					return
				}
			}

			// Check if received msg is a ping and respond with a text pong
			if (string(msg)) == "ping" {
				err := c.WriteMessage(websocket.TextMessage, []byte("pong"))
				if err != nil {
					return
				}
			}
		}

		// Listen to new JSON encoded message and operate accordingly
		//for {
		//	var message map[string]interface{}
		//	err := c.ReadJSON(&message)
		//	if err != nil {
		//		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		//			logErr(c, err, "DocsUpdate WS read error")
		//		} else {
		//			logErr(c, err, "DocsUpdate WS error")
		//		}
		//		return // Call defer and close the connection
		//	}
		//}
	})
}

func IssueUpdate() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		// When the function returns, unregister the client and close the connection
		defer func() {
			unregister <- c
			c.Close()
		}()

		// Register the client
		register <- c

		// Listen for message and respond accordingly (pong incoming ping req)
		for {
			mt, _, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logErr(c, err, "DocsUpdate WS read error")
				} else {
					logErr(c, err, "DocsUpdate WS error")
				}
				return // Call defer and close the connection
			}

			if mt == websocket.PingMessage {
				fmt.Println("Ping received")
				err := c.WriteMessage(websocket.PongMessage, nil)
				if err != nil {
					return
				}
			}
		}

		// Listen to new JSON encoded message and operate accordingly
		//for {
		//	var message map[string]interface{}
		//	err := c.ReadJSON(&message)
		//	if err != nil {
		//		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
		//			logErr(c, err, "IssueUpdate WS read error")
		//		} else {
		//			logErr(c, err, "IssueUpdate WS error")
		//		}
		//		return // Call defer and close the connection
		//	}
		//}

	})
}

func NewsUpdate() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		// When the function returns, unregister the client and close the connection
		defer func() {
			unregister <- c
			c.Close()
		}()

		// Register the client
		register <- c

		// Listen for message and respond accordingly (pong incoming ping req)
		for {
			mt, msg, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logErr(c, err, "NewsUpdate WS read error")
				} else {
					logErr(c, err, "NewsUpdate WS error")
				}
				return // Call defer and close the connection
			}

			if mt == websocket.PingMessage {
				fmt.Println("Ping received")
				err := c.WriteMessage(websocket.PongMessage, nil)
				if err != nil {
					return
				}
			}

			// Check if received msg is a ping and respond with a text pong
			if (string(msg)) == "ping" {
				err := c.WriteMessage(websocket.TextMessage, []byte("pong"))
				if err != nil {
					return
				}
			}
		}
	})
}
