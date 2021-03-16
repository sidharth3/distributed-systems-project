package senders

import (
	"bytes"
	"ds-proj/slave/helpers"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func RegisterWithMaster() {
	dirList := helpers.ListDir()
	dirList["/"+helpers.IP()] = true // filnames would not have / character. Pass the return address.
	filesBytes, err := json.Marshal(dirList)
	if err != nil {
		log.Fatal()
	}

	req, err := http.NewRequest("POST", "http://"+helpers.MasterIP()+"/register", bytes.NewBuffer(filesBytes))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout: time.Second * helpers.TIMEOUT,
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Fatal("Failed to register with master.")
	} else {
		fmt.Println("Successfully registered with master.")
	}
}
