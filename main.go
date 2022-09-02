package main

import(
	"fmt"
	// "io"
	// "log"
	"net/http"
	// "strings"
	"os"
	
	"github.com/gorilla/websocket"
	"github.com/xhfmvls/go-chat-app/pkg/controllers"
)

const (
	NewUserMessage = "new user"
	ChatMessage    = "chat"
	LeaveMessage   = "leave"
)

type WebSocketConnection struct {
	*websocket.Conn
	Username string
}

type SocketPayload struct {
	Message string
}

type ScoketResponse struct {
	User string
	MessageType string
	Message string
}

// var conn = make([]*WebSocketConnection, 0)

func main() {
	http.HandleFunc("/", controllers.Content)
	http.HandleFunc("/ws", controllers.WebSocket)
	
	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT not available")
	}
	port = fmt.Sprintf(":%s", port)

	fmt.Printf("[+] Server Listrning on Port %s", port)
	http.ListenAndServe(port, nil)
}