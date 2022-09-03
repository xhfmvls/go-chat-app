package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

type SocketPayload struct {
	Message string
}

type SocketResponse struct {
	User        string
	MessageType string
	Message     string
}

type WebSocketConnection struct {
	*websocket.Conn
	Username string
}

var Conns = make([]*WebSocketConnection, 0)

const MESSAGE_NEW_USER = "New User"
const MESSAGE_CHAT = "Chat"
const MESSAGE_LEAVE = "Leave"

func main() {
	http.HandleFunc("/", Content)
	http.HandleFunc("/ws", WebSocket)

	godotenv.Load(".env")
	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT not available")
	}
	port = fmt.Sprintf(":%s", port)

	fmt.Printf("[+] Server Listrning on Port %s\n", port)
	http.ListenAndServe(port, nil)
}

func Content(w http.ResponseWriter, r *http.Request) {
	f, err := os.ReadFile("index.html")

	if err != nil {
		fmt.Println("File not Found")
		return
	}

	fmt.Fprintf(w, "%s", f)
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	gorillaConn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)

	if err != nil {
		fmt.Println("Connection Error")
		return
	}

	username := r.URL.Query().Get("username")
	NewConn := WebSocketConnection{Conn: gorillaConn, Username: username}
	Conns = append(Conns, &NewConn)

	go IoHandle(&NewConn, Conns)
}

func IoHandle(currConn *WebSocketConnection, conns []*WebSocketConnection) {
	log.Println("Input Output Handler")
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recover Error")
		}
	}()

	BroadcastMessage(currConn, MESSAGE_NEW_USER, "")

	for {
		payload := SocketPayload{}
		err := currConn.ReadJSON(&payload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				log.Println("Leave")
				BroadcastMessage(currConn, MESSAGE_LEAVE, "")
				// EjectConnection(currConn)
				return
			}

			log.Println("Internal Error")
			continue
		}

		BroadcastMessage(currConn, MESSAGE_CHAT, payload.Message)
	}
}

func BroadcastMessage(currConn *WebSocketConnection, msgType string, message string) {
	log.Println("Broadcast", message)
	for _, conn := range Conns {
		if conn == currConn {
			log.Println("Connection Valid")
			continue
		}

		conn.WriteJSON(SocketResponse{
			User:        currConn.Username,
			MessageType: msgType,
			Message:     message,
		})
	}
}

func EjectConnection(currConn *WebSocketConnection) {
	log.Println("Eject")
	for i, conn := range Conns {
		if conn == currConn {
			Conns[i] = Conns[len(Conns)-1]
			return
		}
	}
}
