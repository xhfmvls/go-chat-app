package controllers

import (
	"log"
	"strings"

	"github.com/xhfmvls/go-chat-app/pkg/models"
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
	NewUserMessage = "New User"
	ChatMessage = "Chat"
	LeaveMessage = "Leave"
)

func IoHandle(currConn *WebSocketConnection, conns []*WebSocketConnection) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recover Error")
		}
	}()

	BroadcastMessage(currConn, NewUserMessage, "")

	for {
		payload := SocketPayload{}
		err := currConn.ReadJSON(&payload)
		if err != nil {
			if strings.Contains(err.Error(), "websocket: close") {
				BroadcastMessage(currConn, LeaveMessage, "")
				EjectConnection(currConn)
				return
			}

			continue
		}

		models.InsertNewChat(currConn.Username, payload.Message)
		BroadcastMessage(currConn, ChatMessage, payload.Message)
	}
}

func BroadcastMessage(currConn *WebSocketConnection, msgType string, message string) {
	log.Println("Broadcast", msgType, message)
	for _, conn := range Conns {
		if conn == currConn {
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
	for i, conn := range Conns {
		if conn == currConn {
			Conns[i] = Conns[len(Conns)-1]
			return
		}
	}
}
