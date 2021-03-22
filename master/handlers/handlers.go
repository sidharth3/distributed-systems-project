package handlers

import (
	"ds-proj/master/structs"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Sends an array of strings over to the client. [ip1, ip2, ip3]
func HandleSlaveIPs(m *structs.Master) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
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

		slave := structs.Slave{slaveIP, structs.UNDERLOADED, files}
		newSlave(m, &slave)
	}
}
