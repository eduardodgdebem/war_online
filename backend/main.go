package main

import (
	"log"
	"net/http"

	ws "war-backend/websocket"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ws", ws.WsHandler)

	log.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "home.html")
}
