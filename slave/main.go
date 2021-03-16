package main

import (
	"ds-proj/slave/handlers"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	IP = "127.0.0.1"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Please enter a port number.")
	} else {
		err := os.MkdirAll("files_"+os.Args[1], os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		http.HandleFunc("/file", handlers.DownloadFile)
		log.Fatal(http.ListenAndServe("127.0.0.1:"+os.Args[1], nil))
	}
}
