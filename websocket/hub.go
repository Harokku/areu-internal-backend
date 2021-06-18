package websocket

import (
	"github.com/gofiber/websocket/v2"
	"log"
	"time"
)

type client struct{}

var clients = make(map[*websocket.Conn]client) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
var register = make(chan *websocket.Conn)
var Broadcast = make(chan string)
var unregister = make(chan *websocket.Conn)

//RunHub actually start the connection hub to manage reg/unreg and broadcast of message. Use channels to control connected clients and messages broadcast
func RunHub() {
	for {
		select {
		case connection := <-register:
			clients[connection] = client{}
			log.Printf("%s - New client connected", time.Now().Format(time.ANSIC))

		case message := <-Broadcast:
			log.Printf("%s - Received new update, broadcasting...", time.Now().Format(time.ANSIC))

			// send message to all registered clients
			for connection := range clients {
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Printf("%s - Error sending udpate wia ws: %s", time.Now().Format(time.ANSIC), err)

					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
					delete(clients, connection)
				}
			}

		case connection := <-unregister:
			delete(clients, connection)
			log.Printf("%s - Client disconnected", time.Now().Format(time.ANSIC))
		}
	}
}
