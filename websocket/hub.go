package websocket

import (
	"github.com/gofiber/websocket/v2"
	"log"
)

type client struct {
	remoteIp string
}

// Constant definition for WS message types
const (
	Filewatcher = "Filewatcher event"
	Issue       = "Issue Event"
)

var clients = make(map[*websocket.Conn]client) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
var register = make(chan *websocket.Conn)
var Broadcast = make(chan map[string]interface{})
var unregister = make(chan *websocket.Conn)

// RunHub actually start the connection hub to manage reg/unreg and broadcast of message. Use channels to control connected clients and messages broadcast
func RunHub() {
	for {
		select {
		case connection := <-register:
			clients[connection] = client{remoteIp: connection.Locals("remoteIp").(string)}
			logEvent(connection, "New client connected")

		case message := <-Broadcast:
			log.Printf(" - Received new update, broadcasting...")

			// send message to all registered clients
			for connection := range clients {
				if err := connection.WriteJSON(message); err != nil {
					logErr(connection, err, "Error sending update wia ws")

					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
					delete(clients, connection)
				}
			}

		case connection := <-unregister:
			delete(clients, connection)
			logEvent(connection, "Client disconnected")
		}
	}
}
