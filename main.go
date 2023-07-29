package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // accept all origins for simplicity, adjust as needed
	},
}

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

func main() {
	fmt.Println("Starting server on port 3333")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServeTLS(":3333", "cert.pem", "key.pem", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	for {
		var message Message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		msg := p
		err = json.Unmarshal(msg, &message)
		if err != nil {
			log.Println("Error in unmarshalling:", err)
			return
		}

		fmt.Printf("\n %s |\t%s: %s\n-----------------------", fmt.Sprint(messageType), message.Username, message.Content)
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}
