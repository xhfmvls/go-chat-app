package models

import (
	"log"

	"github.com/xhfmvls/go-chat-app/pkg/config"
)

type Chat struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func InsertNewChat(username, msg string) {
	db, err := config.Connect()
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	_, err = db.Collection("chat").InsertOne(config.Ctx, Chat{
		Name:    username,
		Message: msg,
	})

	if err != nil {
		log.Fatal(err.Error())
		return
	}
	log.Println("Data save success")
}
