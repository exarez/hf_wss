package main

import (
	"log"
	"net/http"
)

func InitServer() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":3333", nil))
}

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// 	CheckOrigin: func(r *http.Request) bool {
// 		return true // accept all origins for simplicity, adjust as needed
// 	},
// }

// func handler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// 	defer conn.Close()

// 	// go handleInput(conn)

// 	for {
// 		var message Message
// 		messageType, p, err := conn.ReadMessage()
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		msg := p
// 		err = json.Unmarshal(msg, &message)
// 		if err != nil {
// 			log.Println("Error in unmarshalling:", err)
// 			return
// 		}

// 		fmt.Printf("\n %s |\t%s: %s\n-----------------------", fmt.Sprint(messageType), message.Username, message.Content)
// 		if err := conn.WriteMessage(messageType, p); err != nil {
// 			log.Println(err)
// 			return
// 		}
// 	}
// }
