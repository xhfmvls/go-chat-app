package controllers

import (
	"net/http"
	"strings"
	"log"
)

type SocketPayload struct {
	Message string
}

type SocketResponse struct {
	User string
	MessageType string
	Message string
}


const (
	NewUserMessage = "new user"
	ChatMessage    = "chat"
	LeaveMessage   = "leave"
)

func BroadcastMessage(currConn *WebSocketConnection, msgType string, message string) {
	log.Println("Broadcast")
	for _, conn := range Conns {
		if conn == currConn {
			continue
		}
		
		conn.WriteJSON(SocketResponse{
			User: conn.Username,
			MessageType: msgType,
			Message: message,
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

func IoHandle(currConn *WebSocketConnection, conns []*WebSocketConnection) {
	log.Println("Input Output Handler")
	defer func() {
		if r := recover(); r != nil {
			log.Println(http.StatusInternalServerError)
			return
		}
	}()

	BroadcastMessage(currConn, NewUserMessage, "")

	for {
		payload := SocketPayload{}
		err := NewConn.ReadJSON(&payload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				BroadcastMessage(currConn, ChatMessage, payload.Message)
				EjectConnection(currConn)
				return
			}

			log.Println(http.StatusInternalServerError)
			continue
		}
		
		log.Println(payload.Message)
		BroadcastMessage(currConn, ChatMessage, payload.Message)
	}
}