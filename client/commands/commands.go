package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func getFileSlave(slave_ip string, filename string) {
	res, err := http.Get("http://" + slave_ip + "/file?file=" + filename)
	if err != nil {
		log.Fatal(err)
	}

	//only for verbose
	if res.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		outputString := string(body)
		fmt.Println(outputString)

	}

	fmt.Println(res)
}

func GetFileMaster(master_ip string, filename string) {
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
	getFileSlave(ipArr[0], filename)
}
