package main

import (
	"fmt"
	"image/color"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type InputMessage struct {
	Username string `json:"userName"`
	Type     string `json:"type"`
	Content  string `json:"content"`
}
type OutMessage struct {
	Type  string    `json:"type"`
	State GameBoard `json:"state"`
}
type Position struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
}
type Ball struct {
	Position  `json:"position"`
	Radius    float32 `json:"radius"`
	XVelocity float32 `json:"xVelocity"`
	YVelocity float32 `json:"yVelocity"`
	Color     string  `json:"color"`
}
type Command struct {
	key string
}
type GameBoard struct {
	Element Ball `json:"element"`
}
type Client struct {
	username string
	id       string
	ws       *websocket.Conn
	command  chan Command
}

var game = GameBoard{}
var universalCommand = make(chan Command)
var clients []*Client
var messages []*InputMessage
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var frequency = 36

const (
	initBallVelocity = 5.0
	initPaddleSpeed  = 10.0
	speedUpdateCount = 6
	speedIncrement   = 0.5
)

var (
	BgColor  = color.Black
	ObjColor = color.RGBA{120, 226, 160, 255}
)

const (
	windowWidth  = 400
	windowHeight = 400
)
const (
	InitBallRadius = 10.0
)

func (b *Ball) Update() {
	h := windowHeight
	w := windowWidth
	b.X += b.XVelocity
	b.Y += b.YVelocity

	// bounce off edges when getting to top/bottom
	if b.Y-b.Radius > float32(h) {
		b.YVelocity = -b.YVelocity
		b.Y = float32(h) - b.Radius
	} else if b.Y+b.Radius < 0 {
		b.YVelocity = -b.YVelocity
		b.Y = b.Radius
	}
	if b.X-b.Radius > float32(w) {
		b.XVelocity = -b.XVelocity
		b.X = float32(w) - b.Radius
	} else if b.X+b.Radius < 0 {
		b.XVelocity = -b.XVelocity
		b.X = b.Radius
	}
	// // bounce off paddles
	// if b.X-b.Radius < leftPaddle.X+float32(leftPaddle.Width/2) &&
	// 	b.Y > leftPaddle.Y-float32(leftPaddle.Height/2) &&
	// 	b.Y < leftPaddle.Y+float32(leftPaddle.Height/2) {
	// 	b.XVelocity = -b.XVelocity
	// 	b.X = leftPaddle.X + float32(leftPaddle.Width/2) + b.Radius
	// } else if b.X+b.Radius > rightPaddle.X-float32(rightPaddle.Width/2) &&
	// 	b.Y > rightPaddle.Y-float32(rightPaddle.Height/2) &&
	// 	b.Y < rightPaddle.Y+float32(rightPaddle.Height/2) {
	// 	b.XVelocity = -b.XVelocity
	// 	b.X = rightPaddle.X - float32(rightPaddle.Width/2) - b.Radius
	// }
}

//--- Game --- //
func GameLoop(update func(), post func(), delay time.Duration) chan bool {
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
	return stop
}
func (g *GameBoard) Update(command chan Command) {

	press := <-command
	key := press.key
	switch key {
	case "ARROWUP":
		g.Element.Y -= 10
	case "ARROWDOWN":
		g.Element.Y += 10
	case "A":
		g.Element.X -= 10
	case "D":
		g.Element.X += 10
	}

}
func (g *GameBoard) Post() {
	for _, client := range clients {
		err := client.ws.WriteJSON(OutMessage{Type: "GAME_STATE", State: *g})
		if err != nil {
			log.Printf("ERROR:%+v", err)
		}
	}
}
func (g *GameBoard) RunGame(command chan Command) chan bool {
	update := func() {
		g.Update(command)

	}

	post := func() {
		g.Element.Update()
		g.Post()

	}
	return GameLoop(update, post, time.Duration(frequency)*time.Millisecond)
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

		if len(clients) == 1 {
			game.Element = Ball{
				Position: Position{
					X: float32(windowWidth / 2),
					Y: float32(windowHeight / 2)},
				Radius:    InitBallRadius,
				Color:     fmt.Sprintf("%+v", ObjColor),
				XVelocity: initBallVelocity,
				YVelocity: initBallVelocity,
			}
			game.RunGame(universalCommand)
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
