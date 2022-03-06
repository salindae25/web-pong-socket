package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Message struct {
	Username string `json:"userName"`
	Type     string `json:"type"`
	Content  string `json:"content"`
}
type Client struct {
	username string
	id       string
	ws       *websocket.Conn
}

var clients []*Client
var messages []*Message
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

//---- Message ----//
func (m *Message) post() {
	m.broadCast()
	messages = append(messages, m)
}
func (m *Message) broadCast() {
	for _, client := range clients {
		m.broadCastTo(client)
	}
}

func (m *Message) broadCastTo(client *Client) {
	err := client.ws.WriteJSON(m)
	if err != nil {
		log.Printf("ERROR:%+v", err)
	}
}
func sendPreviousMessageTo(client *Client) {
	for _, m := range messages {
		m.broadCastTo(client)
	}
}
func handleIncomeMessage(client *Client, m Message) {
	switch m.Type {
	case "INITIAL_CONN":
		newMsgToChat := Message{
			Type:     "NEW_USER",
			Username: "System",
			Content:  client.username,
		}
		newMsgToChat.post()
		sendPreviousMessageTo(client)
		newMsgToClient := Message{
			Type:     "USER_NAME",
			Username: "System",
			Content:  client.username,
		}
		newMsgToClient.broadCastTo(client)

	case "LEAVE":
		exitMsg := Message{
			Type:     "EXIT",
			Username: "System",
			Content:  fmt.Sprintf("%v has left the chat", client.username),
		}
		exitMsg.post()

	case "MESSAGE":
		msg := Message{
			Type:     "MESSAGE",
			Username: client.username,
			Content:  m.Content,
		}
		msg.post()

	}
}

//---- Client ----//
func (client *Client) startListening() {
	for {
		var msg Message
		err := client.ws.ReadJSON(&msg)
		if err != nil {
			releaseConnection(client)
			handleIncomeMessage(client, Message{Type: "LEAVE"})
			break
		}

		handleIncomeMessage(client, msg)
	}
}
func releaseConnection(client *Client) {
	index := -1
	for idx, val := range clients {
		if client == val {
			index = idx
			break
		}
	}
	if index >= 0 {
		clients = append(clients[:index], clients[index+1:]...)
	} else {
		log.Println("Try to RemoveClient not existing")
	}
	client.ws.Close()
}

//------ server -------//
func wsEndPoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	// add connection to list of connections
	newUser := Client{
		id:       uuid.New().String(),
		username: strings.Split(uuid.New().String(), "-")[0],
		ws:       ws,
	}
	handleIncomeMessage(&newUser, Message{Type: "INITIAL_CONN"})
	clients = append(clients, &newUser)

	newUser.startListening()

}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "root route")
}

func setupRoutes() {
	http.HandleFunc("/", root)
	http.HandleFunc("/ws", wsEndPoint)
}

func main() {
	setupRoutes()
	println("Starting Web Server ")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
