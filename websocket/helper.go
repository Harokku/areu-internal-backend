package websocket

import (
	"github.com/gofiber/websocket/v2"
	"log"
)

func logErr(c *websocket.Conn, err error, msg string) {
	log.Printf(" - %s |\t %s:\t%s", c.Locals("remoteIp").(string), msg, err)
}

func logEvent(c *websocket.Conn, msg string) {
	log.Printf(" - %s |\t %s", c.Locals("remoteIp").(string), msg)
}
