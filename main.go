package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/gorilla/websocket"
)

var (
	testing bool = true
)

type Message struct {
	Username string `json:"username"`
	Content  string `json:"content"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // accept all origins for simplicity, adjust as needed
	},
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interrupt
		fmt.Println("Received an interrupt signal, shutting down...")
		// perform clean up if necessary
		os.Exit(0)
	}()

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServeTLS(":3333", "cert.pem", "key.pem", nil))

	// u := url.URL{Scheme: "wss", Host: "127.0.0.1:3333", Path: "/"}

	// c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	// if err != nil {
	// 	log.Fatal("dial:", err)
	// }
	// go func() {
	// 	defer c.Close()
	// 	for {
	// 		_, message, err := c.ReadMessage()
	// 		if err != nil {
	// 			log.Println("read:", err)
	// 			return
	// 		}
	// 		var msg Message
	// 		err = json.Unmarshal(message, &msg)
	// 		if err != nil {
	// 			log.Println("Unmarshal failed: ", err)
	// 			return
	// 		}
	// 		fmt.Printf("%s: %s\n", msg.Username, msg.Content)
	// 	}
	// }()

	// take input in a loop

	// for {
	// 	fmt.Print("> ")
	// 	fmt.Scanln(&input)

	// 	if strings.TrimSpace(input) == "" {
	// 		continue
	// 	}

	// 	// create a Message object and send it to the server
	// 	message := Message{Username: "username", Content: input}

	// 	// Here I want to prompt the user to verify the message before sending it
	// 	// to the server. If the user enters "y" then send the message, otherwise
	// 	// don't send it.

	// 	// ...
	// 	fmt.Print("Are you sure you want to send this message? (y/n) ")
	// 	fmt.Scanln(&confirmation)

	// 	if strings.TrimSpace(strings.ToLower(confirmation)) != "y" {
	// 		fmt.Println("Message not sent.")
	// 		input = "" // reset input
	// 		continue
	// 	}

	// 	err := c.WriteJSON(message)
	// 	if err != nil {
	// 		log.Println("write:", err)
	// 		return
	// 	}

	// 	input = "" // reset input
	// 	confirmation = ""
	// }
}

// 	fmt.Println("Starting server on port 3333")
// 	http.HandleFunc("/", handler)
// 	log.Fatal(http.ListenAndServeTLS(":3333", "cert.pem", "key.pem", nil))

// 	if testing {
// 		// Testing bubbletea
// 		program := tea.NewProgram(InitialModel())
// 		if _, err := program.Run(); err != nil {
// 			// handle error
// 			log.Fatalf("Error starting program: %s\n", err)
// 		}
// 	}
// }

func handleInput(conn *websocket.Conn) {
	// Start a goroutine to accept user input
	go func() {
		var input string
		for {
			fmt.Scanln(">", &input)

			// create a Message object and send it to the server
			message := Message{Username: "username", Content: input}
			err := conn.WriteJSON(message)
			if err != nil {
				log.Printf("Error sending message: %v", err)
				break
			}
		}
	}()

	// wait for the user to quit the application
	select {}
}

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

func handler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	confirmation := ""

	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := scanner.Text()
			message := Message{Username: "Ezocain", Content: msg}

			fmt.Print("Are you sure you want to send this message? (y/n) ")
			fmt.Scanln(&confirmation)

			if strings.TrimSpace(strings.ToLower(confirmation)) != "y" {
				fmt.Println("Message not sent.")
				confirmation = ""
				continue
			}

			if err := conn.WriteJSON(message); err != nil {
				log.Println("write:", err)
				return
			}
		}
	}()

	for {
		var message Message
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		err = json.Unmarshal(p, &message)
		if err != nil {
			log.Println("Error in unmarshalling:", err)
			return
		}

		fmt.Printf("\n%s: %s\n-----------------------\n", message.Username, message.Content)
	}
}
