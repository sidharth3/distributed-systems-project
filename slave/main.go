package main

import (
	"log"
	"net/http"
	"slave/handlers"
)

func main() {
	http.HandleFunc("/file", handlers.DownloadFileHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
