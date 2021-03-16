package handlers

import (
	"ds-proj/master/structs"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func HandleFile(m *structs.Master) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		filename := req.Form["file"][0]
		w.Header().Set("Content-Type", "application/json")
		ipArr := make([]string, 0)
		for slave := range m.DirectoryTable[filename] {
			ipArr = append(ipArr, slave.IP)
		}
		data, err := json.Marshal(ipArr)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(data)
	}
}

func newSlave(m *structs.Master, slave *structs.Slave) {
	m.Lock.Lock()
	m.Slaves[slave] = true
	for file := range slave.Files {
		if m.DirectoryTable[file] == nil {
			m.DirectoryTable[file] = make(map[*structs.Slave]bool)
		}
		m.DirectoryTable[file][slave] = true
	}
	m.Lock.Unlock()
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

		slave := structs.Slave{slaveIP, files, structs.UNDERLOADED}
		newSlave(m, &slave)
	}
}
