package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Master struct {
	IP             string
	Slaves         map[*Slave]Status   // Slave instance => Status
	DirectoryTable map[string][]*Slave // "/foo/bar.txt" => Slave instance
}

type Slave struct {
	IP    string
	ID    int
	Files []string // "/foo/bar.txt" => byte array
}

// Status is an enumerated type
type Status string

const (
	OVERLOADED  Status = "Current load exceeds threshold"
	UNDERLOADED Status = "Current load does not exceed threshold"
	DEAD        Status = "Cannot be pinged"
)

func (m *Master) handleFile(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	filename := req.Form["file"][0]
	w.Header().Set("Content-Type", "application/json")
	ipArr := make([]string, 0)
	for _, slave := range m.DirectoryTable[filename] {
		ipArr = append(ipArr, slave.IP)
	}
	data, err := json.Marshal(ipArr)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(data)
}

func main() {
	defaultFileArr := []string{"test_file.txt"}
	defaultSlave := &Slave{"127.0.0.1:8081", 0, defaultFileArr}
	defaultSlaveMap := make(map[*Slave]Status)
	defaultSlaveMap[defaultSlave] = UNDERLOADED
	defaultDirTable := make(map[string][]*Slave)
	defaultDirTable["test_file.txt"] = []*Slave{defaultSlave}
	master := Master{"127.0.0.1:8080", defaultSlaveMap, defaultDirTable}

	http.HandleFunc("/file", master.handleFile)
	http.ListenAndServe("127.0.0.1:8080", nil)
}
