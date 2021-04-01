package commands

import (
	"bytes"
	"ds-proj/client/helpers"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
)

var responses = make(chan int, 3)

func GetFile(master_ip string, filename string) {
	// Sends a GET request to the master for a list of slave ips with that filename
	ipArr := getFileMaster(master_ip, filename)

	// Sends a GET request to the slave for the file content
	// Currently just using first ip returned
	// fmt.Println(ipArr)
	getFileSlave(ipArr[1], ipArr[0])
}

func PostFile(master_ip string, filename string) {
	// Sends a GET request to the master for a list of available slave ips
	ipArr := getSlaveIPsMaster(master_ip)

	// Sends a POST request to the slave to upload the file content
	// sends file to all alive slaves
	// fmt.Println(ipArr)
	fmt.Println(ipArr)
	//Now, ipArr's last element is the uid of the operation so must remove the last slice.

	uid := ipArr[len(ipArr)-1]
	ipArr = ipArr[:len(ipArr)-1]
	// fmt.Println(uid)
	for _, ip := range ipArr {
		postFileSlave(ip, filename, uid)
	}

	// fmt.Println(len(responses))

	if len(responses) < (len(ipArr) - 1) {
		fmt.Println("Upload operation Failed.")
	} else {
		fmt.Println("Upload operation Success.")
	}
	for len(responses) > 0 {
		<-responses
	}
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

	// fmt.Println(res)
}

func getSlaveIPsMaster(master_ip string) []string {
	// Sends a GET request to master for available slave ips
	res, err := http.Get("http://" + master_ip + "/slaveips")
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

func postFileSlave(slave_ip string, filename string, uid string) (err error) {
	slaveURL := "http://" + slave_ip + "/upload"

	f := helpers.OpenFile(path.Join(helpers.StorageDir(), filename))
	hashValue := helpers.HashFileContent(f)

	// prepare the reader instances to encode
	values := map[string]io.Reader{
		"filename":  helpers.OpenFile(path.Join(helpers.StorageDir(), filename)),
		"uid":       strings.NewReader(uid),
	}

	// Prepare a form that you will submit to the URL
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	for key, r := range values {
		var fw io.Writer
		if x, ok := r.(io.Closer); ok {
			defer x.Close()
		}
		// Add file
		if _, ok := r.(*os.File); ok {
			if fw, err = w.CreateFormFile(key, hashValue); err != nil {
				return err
			}
		} else { // Add other fields
			if fw, err = w.CreateFormField(key); err != nil {
				return err
			}
		}
		if _, err = io.Copy(fw, r); err != nil {
			return err
		}
	}
	w.Close()

	// Post request
	res, err := http.Post(slaveURL, w.FormDataContentType(), &b)

	if err != nil || res.StatusCode != 200 {
		fmt.Println("File upload has failed.")
		return err
	}

	// only for verbose
	if res.StatusCode == http.StatusOK {
		responses <- 1
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		outputString := string(body)
		fmt.Println(outputString)
	}
	return
}
