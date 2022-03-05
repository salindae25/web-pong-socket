### Process

- client( browser) make a connection with server and establish websocket connection
- client was notified wesocket open event.
- web app: ask user to provide username to complete join chat
- once user enter the username web app send info with a commond "join".
- Server: look at the message
  - if "join": broadcast all client that the user new user join
  - if "msg": broadcast the message
  - if "quit": broadcast that user exiting chat and remove him from activ clients

```mermaid
sequenceDiagram
    participant c1 as Client1
    participant c2 as Client2
    participant c3 as Client3
    participant s as Server
    c1->>s: Initiate web Socket connection
    activate s
    s->>c1: Notify open event
    loop Username
        c1-->>c1: ask to enter username?
        c1-->>c1: user enter username
    end
        c1-)s: send user name with 'join' cmd
        s-->> c1: new user join
        s-->> c3: new user join
        s-->> c2: new user join

    loop communicate via chat
        c1-)s: send message with  "msg" cmd
        s-->> c3: new message
        s-->> c2: new message
    end


        c1-)s: send message with  "quit" cmd
        s-->> c3: notify client1 exiting chat
        s-->> c2: notify client1 exiting chat

    deactivate s
```
