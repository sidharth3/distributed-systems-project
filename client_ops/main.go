package main

import (
	"fmt"
	// "io/ioutil"
	// "log"
	"net/http"
)

const MASTER_URL string = "http://localhost:8080/"

func getfile_slave( slave_ip string,  filename string) {
	res, err := http.Get(slave_ip + "/file?file=" + filename)
	fmt.Println(res)
	fmt.Println(err)
}

func getfile_master(master_ip string, filename string) { 
	res, err := http.Get(master_ip + "/file?file=" + filename)
	fmt.Println(res)
	fmt.Println(err)
}

func main() {
	getfile_slave("http://localhost:8081", "test_file.txt")
	// res, err := http.Get(MASTER_URL)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // client := &http.Client{}

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// defer res.Body.Close()

	// slave_port, err := ioutil.ReadAll(res.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println(string(slave_port))

}
