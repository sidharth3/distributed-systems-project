package main

import (
	"log"
	"net/http"
)

func requestHandler(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		http.NotFound(w, req)
		return
	}
	switch req.Method {
	case "GET":
		w.Write([]byte("9000"))
	case "POST":
		w.Write([]byte("Received a POST request\n"))
	default:
		w.WriteHeader(http.StatusNotImplemented)
		w.Write([]byte(http.StatusText(http.StatusNotImplemented)))
	}
}

func main() {

	http.HandleFunc("/", requestHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// client will send the request for a file, master will tell client which slave
// to contact (for now assuming different slaves are different ports)
