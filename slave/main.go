package main

import (
	"ds-proj/slave/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/file", handlers.DownloadFile)
	log.Fatal(http.ListenAndServe("127.0.0.1:8081", nil))
}
