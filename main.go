package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/xhfmvls/go-chat-app/pkg/controllers"
)

func main() {
	http.HandleFunc("/", controllers.Content)
	http.HandleFunc("/ws", controllers.WebSocket)

	godotenv.Load(".env")
	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT not available")
	}
	port = fmt.Sprintf(":%s", port)

	fmt.Printf("[+] Server Listrning on Port %s\n", port)
	http.ListenAndServe(port, nil)
}