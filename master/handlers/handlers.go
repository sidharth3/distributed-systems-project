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
				m.NLock.Lock()
				m.Namespace[qElement.Filename] = receivedHash
				m.NLock.Unlock()
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
		ipArr := make([]string, 0)
		m.SLock.Lock()
		for slave := range m.Slaves {
			// TODO: some way to select the 3 most free slaves
			ipArr = append(ipArr, slave.IP)
			if len(ipArr) == 3 {
				break
			}
		}
		m.SLock.Unlock()

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
		w.Header().Set("Content-Type", "application/json")
		ipArr := make([]string, 0)

		m.NLock.Lock()
		ipArr = append(ipArr, m.Namespace[filename])

		m.NLock.Unlock()

		m.FLock.Lock()
		for ip := range m.FileLocations[ipArr[0]] {
			ipArr = append(ipArr, ip)
		}
		m.FLock.Unlock()

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

func newSlave(m *structs.Master, slave *structs.Slave) {
	m.SLock.Lock()
	m.Slaves[slave] = true
	m.SLock.Unlock()
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

		slave := structs.Slave{IP: slaveIP, Status: structs.UNDERLOADED, Files: files}
		newSlave(m, &slave)
	}
}
