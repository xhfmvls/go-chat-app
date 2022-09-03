package controllers

import(
	"net/http"

	"github.com/gorilla/websocket"
)

type WebSocketConnection struct {
	*websocket.Conn
	Username string
}

var Conns = make([]*WebSocketConnection, 0)
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var NewConn WebSocketConnection

func WebSocket(w http.ResponseWriter, r *http.Request) {
	gorillaConn, err := Upgrader.Upgrade(w, r, w.Header())

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	username := r.URL.Query().Get("username")
	NewConn = WebSocketConnection{Conn: gorillaConn, Username: username}
	Conns = append(Conns, &NewConn)

	go IoHandle(&NewConn, Conns)
}