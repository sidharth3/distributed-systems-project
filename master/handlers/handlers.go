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
		err := req.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		filename := req.Form["file"][0]
		hash := req.Form["hash"][0]

		m.Namespace.SetHash(filename, hash)
		m.GCCount.NewFile(filename)

		ipArr := m.Slaves.GetFree()
		if len(ipArr) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

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
		err := req.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		filename := req.Form["file"][0]
		ipArr := make([]string, 0)
		ipArr = append(ipArr, m.Namespace.GetHash(filename))
		if ipArr[0] == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		ipArr = append(ipArr, m.FileLocations.GetIPs(ipArr[0])...)
		if len(ipArr) == 1 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

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
		if !m.Namespace.DelFile(filename) {
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

func HandleListDir(m *structs.Master) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		path := req.Form["ls"][0]

		file := m.Namespace.GetFile(path)

		data, err := json.Marshal(file)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(data)
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
		m.Slaves.NewSlave(slaveIP, 0, files)
	}
}
