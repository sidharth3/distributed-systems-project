package senders

import (
	"bytes"
	"ds-proj/slave/config"
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

	client := &http.Client{
		Timeout: time.Second * config.TIMEOUT,
	}

	res, err := client.Post("http://"+helpers.MasterIP()+"/register", "application/json", bytes.NewBuffer(filesBytes))
	if err != nil {
		log.Fatal(err)
	} else if res.StatusCode != http.StatusOK {
		log.Fatal(res)
	} else {
		fmt.Println("Successfully registered with master")
	}
}
