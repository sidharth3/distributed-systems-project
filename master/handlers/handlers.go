package handlers

import (
	"ds-proj/master/structs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

//listens to post requests from slave for uid and file hash
func HandleUpdate(m *structs.Master) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if err := req.ParseForm(); err != nil {
			log.Fatal(err)
		}

		receivedHash := req.Form["filename"][0]
		receivedUid := strings.Trim(fmt.Sprint(req.Form["uid"][0]), "[]")
		q := m.Queue.ReturnObj()

		for _, qElement := range q {
			if qElement.Uid == receivedUid {
				m.Namespace.SetHash(qElement.Filename, receivedHash)
				break
			}

		}
	}
}

// Sends an array of strings over to the client. [ip1, ip2, ip3]
func HandleSlaveIPs(m *structs.Master) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		filename := req.Form["file"][0]
		ipArr := m.Slaves.GetFree()

		uid := uuid.NewString()
		ipArr = append(ipArr, uid)

		qItem := structs.QueueItem{Uid: uid, Filename: filename, Hash: ""}

		m.Queue.Enqueue(qItem)

		data, err := json.Marshal(ipArr)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(data)
	}
}

// Sends an array of strings over to the client. [hashValue, ip1, ip2, ip3]
func HandleFile(m *structs.Master) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		filename := req.Form["file"][0]
		ipArr := make([]string, 0)

		ipArr = append(ipArr, m.Namespace.GetHash(filename))

		ipArr = append(ipArr, m.FileLocations.GetIPs(ipArr[0])...)

		data, err := json.Marshal(ipArr)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(data)
	}
}

func HandleDeleteFile(m *structs.Master) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		filenameBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}

		var filename string
		err = json.Unmarshal(filenameBytes, &filename)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)

		// Check if queue is empty, delete file
		// Else, append file to queue
	}
}

func HandleNewSlave(m *structs.Master) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		filesBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		files := make(map[string]bool)
		err = json.Unmarshal(filesBytes, &files)
		if err != nil {
			log.Fatal(err)
		}
		var slaveIP string
		for file := range files {
			if file[:1] == "/" {
				slaveIP = file[1:]
				break
			}
		}
		delete(files, "/"+slaveIP)
		m.Slaves.NewSlave(slaveIP, structs.UNDERLOADED, files)
	}
}
