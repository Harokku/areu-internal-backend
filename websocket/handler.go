package websocket

import (
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

		// Listen to new JSON encoded message and operate accordingly
		for {
			var message map[string]interface{}
			err := c.ReadJSON(&message)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logErr(c, err, "DocsUpdate WS read error")
				} else {
					logErr(c, err, "DocsUpdate WS error")
				}
				return // Call defer and close the connection
			}
		}
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

		// Listen to new JSON encoded message and operate accordingly
		for {
			var message map[string]interface{}
			err := c.ReadJSON(&message)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					logErr(c, err, "IssueUpdate WS read error")
				} else {
					logErr(c, err, "IssueUpdate WS error")
				}
				return // Call defer and close the connection
			}
		}

	})
}
