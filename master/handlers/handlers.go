package handlers

import (
	"ds-proj/master/config"
	"ds-proj/master/structs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

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
		m.Queue.Confirm(receivedUid, receivedHash)
		m.Commit(receivedUid)
	}
}

// Sends an array of strings over to the client. [ip1, ip2, ip3]
func HandleSlaveIPs(m *structs.Master) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		filename := req.Form["file"][0]
		uid := uuid.NewString()
		m.Queue.Enqueue(uid, filename)

		ipArr := m.Slaves.GetFree()
		ipArr = append(ipArr, uid)

		data, err := json.Marshal(ipArr)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(data)

		// Handle queue timeouts
		go func() {
			<-time.After(time.Second * config.DQTIMEOUT)
			m.Queue.Timeout(uid)
			m.Commit(uid)
		}()
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

		uid := uuid.NewString()
		m.Queue.Enqueue(uid, filename)
		m.Queue.Confirm(uid, "delete")
		m.Commit(uid)
		w.WriteHeader(http.StatusOK)
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
