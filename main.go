package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Message struct {
	Sender string `json:"sender"`
	Text   string `json:"text"`
}

// This variable is in-memory storage of client connections.
// Each WS connection will be put in here on /ws/messages handler.
// And then incoming message will be broadcast to all the connections, including sender.
var clients = make(map[*websocket.Conn]struct{})

func main() {
	// test is our app running
	fmt.Println("Hello from Go Chat App Backend 123")

	// Handles static files from public directory
	// In example http://localhost:3000/testfile.txt
	http.Handle("/", http.FileServer(http.Dir("./public")))

	http.HandleFunc("/ws/messages", wsMessages)

	http.ListenAndServe(":3000", nil)
}

func wsMessages(w http.ResponseWriter, r *http.Request) {
	// Create special upgrader entity. It will help us modify HTTP to WS connection
	upgrader := websocket.Upgrader{}

	// wsConnection is entity of WS connection needed to communicate bidirectionally
	wsConnection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WS upgrade failed ", err.Error())
		return
	}

	// Put newly established connection to in-memory storage
	clients[wsConnection] = struct{}{}

	// The value of the message is JSON with two keys:
	// sender - the author of message. For greetingMessage we used "SERVER" constant.
	// text - a text message a sender sent to all chat participants (receivers).
	greetingMessage := `{"sender":"SERVER", "text":"Hello THERE. You are logged in"}`

	// wsConnection.WriteMessage can be used to send any kind of bytes sequence.
	// Instead of this method you can use wsConnection.WriteJSON to write json.
	err = wsConnection.WriteMessage(websocket.TextMessage, []byte(greetingMessage))
	if err != nil {
		log.Println("Sending message is failed ", err.Error())
		return
	}

	// This way we declare infinite loop to "listen" the incoming messages
	for {
		var msg Message

		// Here the exact place where we do "listen" for next incoming message.
		// wsConnection.ReadJSON pauses execution of loop and return next message when it comes to the connection.
		err = wsConnection.ReadJSON(&msg)
		if err != nil {
			log.Println("Failed to read from connection", err.Error())
			delete(clients, wsConnection)
			return
		}

		// Broadcast the message to all connected clients
		for client := range clients {
			err = client.WriteJSON(msg)
			if err != nil {
				log.Println(err.Error())
				client.Close()
				delete(clients, client)
				return
			}
		}
	}
}
