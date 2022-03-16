package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/salindae25/web-pong/backend/pkg/pong"
)

type InputMessage struct {
	Username string `json:"userName"`
	Type     string `json:"type"`
	Content  string `json:"content"`
}
type OutMessage struct {
	Type  string        `json:"type"`
	State pong.GameSend `json:"state"`
}
type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}

type Command struct {
	key string
}

type Client struct {
	username string
	id       string
	ws       *websocket.Conn
	command  chan Command
}
type Room struct {
	GameBoard *pong.Game
}

var postchannel = make(chan bool)
var game = pong.Game{
	Post: postchannel,
}
var room = Room{
	GameBoard: &game,
}
var universalCommand = make(chan Command)
var clients []*Client
var messages []*InputMessage
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var frequency = 20

//--- Game --- //
func GameLoop(update func(), post func(), broad func(), delay time.Duration) chan bool {
	stop := make(chan bool)
	go func() {
		for {
			post()
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}

		}
	}()
	go func() {
		for {
			update()
		}
	}()
	go func() {
		for {
			broad()
		}
	}()
	return stop
}

func (r *Room) RunGame(command chan Command) chan bool {
	update := func() {
		fmt.Println("update inside RunGame")
		press := <-command
		key := press.key
		r.GameBoard.Update(key)

	}

	post := func() {
		if r.GameBoard.State == pong.PlayState {
			r.GameBoard.Draw()
		}

	}
	broad := func() {
		if <-postchannel {
			r.Post()
		}

	}
	return GameLoop(update, post, broad, time.Duration(frequency)*time.Millisecond)
}
func (r *Room) Post() {
	gameVar := pong.GameSend{
		State:    r.GameBoard.State,
		Ball:     r.GameBoard.Ball,
		Player1:  r.GameBoard.Player1,
		Player2:  r.GameBoard.Player2,
		Rally:    r.GameBoard.Rally,
		Level:    r.GameBoard.Level,
		MaxScore: r.GameBoard.MaxScore,
	}
	for _, client := range clients {
		err := client.ws.WriteJSON(OutMessage{Type: "GAME_STATE", State: gameVar})
		if err != nil {
			log.Printf("ERROR:%+v", err)
		}
	}
}

//---- InputMessage ----//
func (m *InputMessage) post() {
	m.broadCast()
	messages = append(messages, m)
}
func (m *InputMessage) broadCast() {

	for _, client := range clients {
		m.broadCastTo(client)
	}
}
func (m *InputMessage) broadCastTo(client *Client) {
	err := client.ws.WriteJSON(m)
	if err != nil {
		log.Printf("ERROR:%+v", err)
	}
}
func sendPreviousInputMessageTo(client *Client) {
	for _, m := range messages {
		m.broadCastTo(client)
	}
}

func handleIncomeInputMessage(client *Client, m InputMessage) {
	stop := make(chan bool)
	switch m.Type {
	case "INITIAL_CONN":
		newMsgToChat := InputMessage{
			Type:     "NEW_USER",
			Username: "System",
			Content:  client.username,
		}
		newMsgToChat.post()
		sendPreviousInputMessageTo(client)
		newMsgToClient := InputMessage{
			Type:     "USER_NAME",
			Username: "System",
			Content:  client.username,
		}
		newMsgToClient.broadCastTo(client)

		if len(clients) == 0 {
			room.GameBoard.Init()
			room.RunGame(universalCommand)
		}

	case "LEAVE":
		exitMsg := InputMessage{
			Type:     "EXIT",
			Username: "System",
			Content:  fmt.Sprintf("%v", client.username),
		}
		exitMsg.post()
		stop <- true

	case "MESSAGE":
		msg := InputMessage{
			Type:     "MESSAGE",
			Username: client.username,
			Content:  m.Content,
		}
		msg.post()
	case "GAME_CMD":
		universalCommand <- Command{key: m.Content}
	}
}

//---- Client ----//
func (client *Client) startListening() {
	for {
		var msg InputMessage
		err := client.ws.ReadJSON(&msg)
		if err != nil {
			releaseConnection(client)
			handleIncomeInputMessage(client, InputMessage{Type: "LEAVE"})
			break
		}

		handleIncomeInputMessage(client, msg)
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
	fmt.Println("Remove User")
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
		command:  make(chan Command),
	}
	handleIncomeInputMessage(&newUser, InputMessage{Type: "INITIAL_CONN"})
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
