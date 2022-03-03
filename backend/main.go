package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func root(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "root route")
}

func reader(conn *websocket.Conn) {
	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			break
		}
		fmt.Printf("client says: %s\n", msg.Content)
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

	log.Println("Client Suceesfully connected...")

	reader(ws)
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
