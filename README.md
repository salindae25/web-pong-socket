## Phase-1

- App to communicate between frontend and go backend
  1. send set of message from single client one after from frontend [Phase 1]
     - build go based backend to implement web socket
     - send api request/frontend events and log them in backend

---

## learn in this section | 2022-03-04 |

    Implemented golang based backend using gorilla/websocket package instead of the
    standard package ( since standard package missing features).

```go
/* define the webSocket Upgrader read write buffer size there is no absolute value this is based on the
maximum message size you expected to send via websocket.
Note: experiment by changing these value and select value suitable.
*/
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
// this code is for upgrade normal http connection to websocket connection
ws, err := upgrader.Upgrade(w, r, nil) //w: http.ResponseWriter, r: *http.Request
defer ws.Close() // this wil close connection in case of any issue  free memory

```

once we start listening on the websocket we will receive a byte stream we can use `websocket`
package to read this stream.

```go
type Message struct {
	Username string `json:"userName"`
	Type     string `json:"type"`
	Content  string `json:"content"`
}
var msg Message
err := client.ws.ReadJSON(&msg)
 // this wil only read single message. in order to listen to all the message we need to insert this
 inside a infinite for loop
 // here we pass pointer reference to custom Message struct to ReadJSON which will read stream and set value
 to the msg object. incase message failed to set value it will throw an error
```

    from here once we read the message from socket we can process message and take action
    this should cover basic implementation of websocket.

---

2020-03-06

Learn the concept of attaching the method or `func` to a `struct` in golang.it could be used similar to
`this` keyword in other language.

```go
func (m *Message) post() {
	m.broadCast()
	messages = append(messages, m)
}
// above message could be invoke as below

newMsg:= Message{}
newMsg.post() // this method could access all properties in the struct object

// we can also do following
func (m Message) post(){}
 // this would only be copying the values not a reference

```

above concept reorganized code.

---

## Phase-2

1. render set of elements in canvas based on the streamed from backend.
   - when user input shape in text in frontend send it to backend and draw corresponding image
     using backend generated values.
   - based on the key input calculate next state of the object and render it on front end.
