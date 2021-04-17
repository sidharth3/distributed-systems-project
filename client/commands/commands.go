package commands

import (
	"bytes"
	"ds-proj/client/config"
	"ds-proj/client/helpers"
	"ds-proj/client/structs"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"
)

func DownloadDirectory(master_ip string, remote_directory string, local_directory string) {

	extension := path.Ext(remote_directory)

	if extension != "" {
		DownloadFile(master_ip, remote_directory, local_directory)
	} else {

		client := &http.Client{
			Timeout: time.Second * config.TIMEOUT,
		}

		res, err := client.Get("http://" + master_ip + "/ls?ls=" + remote_directory)
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

		err = os.MkdirAll(filepath.Dir(local_directory), os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		for _, temp_directory := range dir {
			DownloadFile(master_ip, remote_directory+temp_directory, local_directory+temp_directory)

		}
	}
}

func DeleteDirectory(master_ip string, remote_directory string) {
	extension := path.Ext(remote_directory)

	if extension != "" {
		DeleteFile(master_ip, remote_directory)
	} else {

		client := &http.Client{
			Timeout: time.Second * config.TIMEOUT,
		}

		res, err := client.Get("http://" + master_ip + "/ls?ls=" + remote_directory)
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

		for _, temp_directory := range dir {
			DeleteFile(master_ip, remote_directory+temp_directory)
		}
	}

}

func DownloadFile(master_ip string, remote_filename string, local_filename string) {

	ipArr := getFileMaster(master_ip, remote_filename)

	client := &http.Client{
		Timeout: time.Second * config.TIMEOUT,
	}

	var res *http.Response
	var err error
	for i := 1; i < len(ipArr); i++ { // Try all slave ips returned
		res, err = client.Get("http://" + ipArr[1] + "/file?file=" + ipArr[0])
		if err != nil {
			fmt.Println(err)
		} else if res.StatusCode != http.StatusOK {
			fmt.Println(res)
		} else {
			break
		}
	}
	if err != nil || res.StatusCode != http.StatusOK {
		log.Fatal("File download failed")
	}

	err = os.MkdirAll(filepath.Dir(local_filename), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create(local_filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Downloaded file")
}

func GetFile(master_ip string, remote_filename string) {
	ipArr := getFileMaster(master_ip, remote_filename)

	client := &http.Client{
		Timeout: time.Second * config.TIMEOUT,
	}

	var res *http.Response
	var err error
	for i := 1; i < len(ipArr); i++ { // Try all slave ips returned
		res, err = client.Get("http://" + ipArr[1] + "/file?file=" + ipArr[0])
		if err != nil {
			fmt.Println(err)
		} else if res.StatusCode != http.StatusOK {
			fmt.Println(res)
		} else {
			break
		}
	}
	if err != nil || res.StatusCode != http.StatusOK {
		log.Fatal("File cat failed")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	outputString := string(body)
	fmt.Println(outputString)
}

func PostFile(master_ip string, filename string, remote_filename string) {
	remote_filename = helpers.SanitizeInput(remote_filename)
	if remote_filename == "" {
		log.Fatal("File path invalid")
	}

	// Sends a GET request to the master for a list of available slave ips
	f := helpers.OpenFile(filename)
	hashValue := helpers.HashFileContent(f)
	f.Close()

	ipArr := getSlaveIPsMaster(master_ip, url.QueryEscape(remote_filename), hashValue) // pass remote filename to master
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
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(http.MethodDelete, masterURL, bytes.NewBuffer(jsonReq))
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Timeout: time.Second * config.TIMEOUT,
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatal(res)
	}

	fmt.Println("Sucessfully deleted file")
}

func ListDir(master_ip string, path string) string {
	client := &http.Client{
		Timeout: time.Second * config.TIMEOUT,
	}
	res, err := client.Get("http://" + master_ip + "/ls?ls=" + path)
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

	fileDir := structs.InitDir("/")
	for _, filename := range dir {
		fileDir.Insert(filename)
	}
	dirStr := fileDir.FormatString()
	// fmt.Println(dirStr)

	return dirStr
}

func getFileMaster(master_ip string, filename string) []string {
	client := &http.Client{
		Timeout: time.Second * config.TIMEOUT,
	}
	res, err := client.Get("http://" + master_ip + "/file?file=" + filename)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatal(res)
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

func getSlaveIPsMaster(master_ip string, remote_filename string, hash string) []string {
	// Sends a GET request to master for available slave ips
	client := &http.Client{
		Timeout: time.Second * config.TIMEOUT,
	}
	res, err := client.Get("http://" + master_ip + "/slaveips?file=" + remote_filename + "&hash=" + hash)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatal(res)
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
	if len(ipArr) == 0 {
		log.Fatal("Upload file failed")
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
	f.Close()
	w.Close()

	// Post request
	client := &http.Client{
		Timeout: time.Second * config.TIMEOUT,
	}
	res, err := client.Post(slaveURL, w.FormDataContentType(), &b)
	if err != nil {
		log.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		log.Fatal(res)
	}

	fmt.Println("Succeeded sending to slave")
}
