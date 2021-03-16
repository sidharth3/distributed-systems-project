package main

import (
	"ds-proj/master/handlers"
	"ds-proj/master/senders"
	"ds-proj/master/test"
	"net/http"
)

func main() {
	master := test.EmptyCase()
	// master := test.TestCase1()

	go senders.HeartbeatSender(master)

	http.HandleFunc("/file", handlers.HandleFile(master))
	http.HandleFunc("/register", handlers.HandleNewSlave(master))
	http.ListenAndServe("127.0.0.1:8080", nil)
}
