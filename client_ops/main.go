package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const MASTER_URL string = "http://localhost:8080/"

func main() {
	res, err := http.Get(MASTER_URL)

	if err != nil {
		log.Fatal(err)
	}

	// client := &http.Client{}

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	slave_port, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(slave_port))

}
