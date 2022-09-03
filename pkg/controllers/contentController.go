package controllers

import(
	"fmt"
	"net/http"
	"os"
)

func Content(w http.ResponseWriter, r *http.Request) {
	f, err := os.ReadFile("index.html")

	if err != nil {
		fmt.Println("File not Found")
		return
	}

	fmt.Fprintf(w, "%s", f)
}