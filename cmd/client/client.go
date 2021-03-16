package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const MASTER_URL string = "http://localhost:8080/"

func getfile_slave(slave_ip string, filename string) {
	res, err := http.Get("http://" + slave_ip + "/file?file=" + filename)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}

func getfile_master(master_ip string, filename string) {
	res, err := http.Get("http://" + master_ip + "/file?file=" + filename)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	ipArr := make([]string, 0)
	err = json.Unmarshal(body, &ipArr)
	if err != nil {
		log.Fatal(err)
	}

	//Add some way to determine best slave
	getfile_slave(ipArr[0], filename)
}

func main() {
	getfile_master("127.0.0.1:8080", "test_file.txt")
}
