package commands

import (
	"bytes"
	"ds-proj/client/config"
	"ds-proj/client/helpers"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

func GetFile(master_ip string, filename string) {
	// Sends a GET request to the master for a list of slave ips with that filename
	ipArr := getFileMaster(master_ip, filename)

	// Sends a GET request to the slave for the file content
	// Currently just using first ip returned
	getFileSlave(ipArr[1], ipArr[0])
}

func PostFile(master_ip string, filename string) {
	// Sends a GET request to the master for a list of available slave ips
	f := helpers.OpenFile(filename)
	hashValue := helpers.HashFileContent(f)
	f.Close()

	ipArr := getSlaveIPsMaster(master_ip, filename, hashValue) // pass filename to master

	// Sends a POST request to the slave to upload the file content
	// sends file to all alive slaves
	for _, ip := range ipArr {
		postFileSlave(ip, filename, hashValue)
	}
}

func DeleteFile(master_ip string, filename string) {
	masterURL := "http://" + master_ip + "/delete"

	// Sends a DELETE request to master to delete the file
	jsonReq, err := json.Marshal(filename)
	req, err := http.NewRequest(http.MethodDelete, masterURL, bytes.NewBuffer(jsonReq))
	client := &http.Client{
		Timeout: time.Second * config.TIMEOUT,
	}
	res, err := client.Do(req)

	if err != nil || res.StatusCode != http.StatusOK {
		log.Fatal("File delete has failed.")
	} else {
		fmt.Println("Successfully deleted file.")
	}
}

func ListDir(master_ip string, path string) {
	res, err := http.Get("http://" + master_ip + "/ls?ls=" + path)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	dir := make([]string, 0)
	err = json.Unmarshal(body, &dir)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
}

func getFileMaster(master_ip string, filename string) []string {
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

	return ipArr
}

func getFileSlave(slave_ip string, hashValue string) {
	res, err := http.Get("http://" + slave_ip + "/file?file=" + hashValue)
	if err != nil {
		log.Fatal(err)
	}

	// only for verbose
	if res.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		outputString := string(body)
		fmt.Println(outputString)
	}

}

func getSlaveIPsMaster(master_ip string, filename string, hash string) []string {
	// Sends a GET request to master for available slave ips
	res, err := http.Get("http://" + master_ip + "/slaveips?file=" + filename + "&hash=" + hash)
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
	fmt.Println(ipArr)

	return ipArr
}

func postFileSlave(slave_ip string, filename string, hash string) {
	slaveURL := "http://" + slave_ip + "/upload"

	f := helpers.OpenFile(filename)

	// Prepare a form that you will submit to the URL
	var b bytes.Buffer
	w := multipart.NewWriter(&b)

	if fw, err := w.CreateFormFile("filename", hash); err != nil {
		log.Fatal(err)
	} else {
		if _, err := io.Copy(fw, f); err != nil {
			log.Fatal(err)
		}
	}
	w.Close()
	f.Close()

	// Post request
	res, err := http.Post(slaveURL, w.FormDataContentType(), &b)

	if err != nil || res.StatusCode != 200 {
		fmt.Println("File upload has failed.")
		log.Fatal(err)
	}

	// only for verbose
	if res.StatusCode == http.StatusOK {
		fmt.Println("Succeeded sending to slave")
	}
}
