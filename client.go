package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
)

/////////////////////////////////////////////////////
// THIS IS MOSTLY FOR DEALING WITH TUI STUFF LATER //
/////////////////////////////////////////////////////
type UserData struct {
	Uid         int    `json:"uid"`
	Username    string `json:"username"`
	Usergroup   int    `json:"usergroup"`
	Group       int    `json:"group"`
	Avatar      string `json:"avatar"`
	Regdate     int    `json:"regdate"`
	Postnum     int    `json:"postnum"`
	Level       int    `json:"level"`
	Comparename string `json:"comparename"`
	Rank        int    `json:"rank"`
}

type InputModel struct {
	Content string
}

type model struct {
	textInput textinput.Model
	err       error
	MessageModel
	InputModel
}

// MessageModel is a Bubbletea model for a list of messages.
type MessageModel struct {
	Messages []Message
	// add a field here for your websocket connection
	Conn *websocket.Conn
}

// Model is the main model for your application.
type Model struct {
	MessageModel
	InputModel
}

func parseUser(userJson string) User {
	var user User
	err := json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		fmt.Println("Error in unmarshalling:", err)
	}
	return user
}

type Rank struct {
	name  string // Ub3r = 28
	order int
}

type User struct {
	username  string
	avatar    string
	usergroup int
	uid       int
	rank      Rank
	postnum   int
}

func (m model) Init() tea.Cmd {
	// Establish a WebSocket connection
	u := url.URL{Scheme: "wss", Host: "localhost:3333", Path: "/"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatalf("Dial failed: %s", err)
	}
	m.MessageModel.Conn = conn

	// Create a command to continuously read incoming messages from the server
	cmd := (func() tea.Msg {
		_, message, err := m.MessageModel.Conn.ReadMessage()
		if err != nil {
			log.Printf("Read failed: %s", err)
			// return an error message if necessary
		}

		var msg Message
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Printf("Unmarshal failed: %s", err)
			// return an error message if necessary
		}

		return msg
	})

	return cmd
}

func InitialModel() model {
	ti := textinput.New()
	ti.Placeholder = "Pikachu"
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		textInput: ti,
		err:       nil,
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case Message:
		m.Messages = append(m.Messages, msg)
	}

	// update your InputModel here if necessary

	return m, nil
}

func (m model) View() string {
	// Render your view here using the lipgloss package
	// For now, let's just print all the messages and the current input
	var output string
	for _, msg := range m.Messages {
		output += fmt.Sprintf("%s: %s\n", msg.Username, msg.Content)
	}
	output += fmt.Sprintf("> %s", m.InputModel.Content)

	return output
}
