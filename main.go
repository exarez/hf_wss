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
	port         = ":3333"
)

type Message struct {
	Username  string `json:"username"`
	Content   string `json:"content"`
	Usergroup int    `json:"usergroup"`
}

var shutDown = make(chan os.Signal, 1)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	title0 := "\033[38;5;40m   __ ______  _________  _  ___   ______  \033[0m"
	title1 := "\033[38;5;118m  / // / __/ / ___/ __ \\/ |/ / | / / __ \\ \033[0m"
	title2 := "\033[38;5;118m / _  / _/  / /__/ /_/ /    /| |/ / /_/ / \033[0m"
	title3 := "\033[38;5;40m/_//_/_/    \\___/\\____/_/|_/ |___/\\____/  \033[0m"
	title := fmt.Sprintf("%s\n%s\n%s\n%s\n", title0, title1, title2, title3)

	status := "\033[38;5;40m[+] \033[0mWaiting for incoming connection on 127.0.0.1" + port + "..."
	fmt.Println(title)
	fmt.Println(status)

	// interrupt := make(chan os.Signal, 1)
	signal.Notify(shutDown, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-shutDown
		fmt.Println("Received an interrupt signal, shutting down...")
		// perform clean up if necessary
		os.Exit(0)
	}()

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServeTLS(port, "cert.pem", "key.pem", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	select {
	case <-shutDown:
		// If we're shutting down, don't accept new requests.
		return
	default:

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

				// TODO: Get username from user input at launch
				message := Message{Username: "Ezocain", Content: msg}

				if testing {
					fmt.Print("Are you sure you want to send this message? (y/n) ")
					fmt.Scanln(&confirmation)
					// ----

					if strings.TrimSpace(strings.ToLower(confirmation)) != "y" {
						fmt.Println("Message not sent.")
						confirmation = ""
						continue
					}
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

			fmt.Printf("%v: %s\n%s-----------------------%s\n", printColoredUser(message.Username, message.Usergroup), message.Content, "\033[3;30m", "\033[0m")
		}
	}
}
