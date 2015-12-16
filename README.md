# WSChannel
Channel for websocket,  broadcast message to others in the same channel. 

This server is just used to redirect and send WebSocket packages without any modification.

## Channel
* URI: `ws://address:port/ws/<channel>`

Clients in the same channel will receive messages of other clients in the same channel.

## Dependencies
* [gorilla websocket](https://github.com/gorilla/websocket)

## Run
* Clone this repo and set up GOPATH
```
go get github.com/gorilla/websocket
go build -o wschannel main.go hub.go conn.go
./wschannel <address> <port>

//wschannel will listen on the port and waiting to connections.
```

  
