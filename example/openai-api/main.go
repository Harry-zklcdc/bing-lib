package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/v1/chat/completions", chatHandler)
	http.HandleFunc("/v1/images/generations", imageHandler)
	log.Println("Server Listening on :8080")
	http.ListenAndServe(":8080", nil)
}
