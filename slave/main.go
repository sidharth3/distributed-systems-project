package main

import (
	"ds-proj/slave/handlers"
	"ds-proj/slave/helpers"
	"ds-proj/slave/senders"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Please enter a port number and a master IP address.")
	} else {
		err := os.MkdirAll(helpers.StorageDir(), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		go senders.RegisterWithMaster()

		http.HandleFunc("/file", handlers.DownloadFile)
		http.HandleFunc("/heartbeat", handlers.HeartbeatHandler)
		log.Fatal(http.ListenAndServe(helpers.IP(), nil))
	}
}
