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
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "root route")
}
func handleMessage(msg Message, client Client) {
	// depend on th type of the message activate diffrent function
	// case type='join' set msg.content as username and broadcast that new user join
	// case type='msg" broad cast message from user
	// case type='quit' remove user from clients list and broadcast client leaving chat
	broadCastMessage(client, Message{Type: msg.Type, Username: client.username, Content: msg.Content})
	// challenges each client is handle by seperate go routine
}
func removeClientActiveList(client Client) {
	// remove current user from clients
	// and broadcast user quit to chat
	broadCastMessage(client, Message{Type: "quit", Username: "System", Content: fmt.Sprintf("%v has exited chat", client.username)})

}
func broadCastMessage(currentClient Client, msg Message) {
	// broadcast to all users except currentclient
	// fmt.Printf("broadcast this to every one except %s:\n%+v\n", currentClient.username, msg)
	for _, eachClient := range clients {
		eachClient.ws.WriteJSON(&msg)
	}
}
func reader(client Client) {
	for {
		var msg Message
		err := client.ws.ReadJSON(&msg)
		if err != nil {

			removeClientActiveList(client)
			break
		}
		handleMessage(msg, client)
	}
}

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
	clients = append(clients, &newUser)
	broadCastMessage(newUser, Message{Type: "newJoin", Username: newUser.username})
	log.Println("Client Suceesfully connected...")

	reader(newUser)
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
