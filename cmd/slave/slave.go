package main

import (
	"ds-proj/handlers/SlaveOps"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/file", SlaveOps.DownloadFile)
	log.Fatal(http.ListenAndServe("127.0.0.1:8081", nil))
}
