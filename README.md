- App to communicate between frontend and go backend
  1. send set of message from single client one after from frontend [Phase 1]
     - build go based backend to implement web socket
     - send api request/frontend events and log them in backend

---

## learn in this section | 2022-03-04 |

    Implemented golang based backend using gorilla/websocket package insted of the standard package ( since standard package missing features).

```go
/* define the webSocket Upgrader read write buffer size there is no absolute value this is based on the  maximum message size you expected to send via websocket.
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

s
once we start listening on the websocket we will recive a byte stream we can use `websocket` package to read this stream.

```go
type Message struct {
	Username string `json:"userName"`
	Type     string `json:"type"`
	Content  string `json:"content"`
}
var msg Message
err := client.ws.ReadJSON(&msg)
 // this wil only read single message. in order to listen to all the meesage we need to insert this inside a inifinite for loop
 // here we passa pointer reference to custom Message struct to ReadJSON which will read stream and set vlue to the msg object. incase message failed to set value it will throw an error
```

    from here once we read the message from socket we can process mesage and take action
    this should cover basic implementation of websocket.
