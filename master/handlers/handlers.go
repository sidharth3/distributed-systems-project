package handlers

import (
	"ds-proj/master/structs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func 

// Sends an array of strings over to the client. [ip1, ip2, ip3]
func HandleSlaveIPs(m *structs.Master) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		filename := req.Form["file"][0]
		hash := req.Form["hash"][0]

		m.Namespace.SetHash(filename, hash)
		m.GCCount.NewFile(filename)

		ipArr := m.Slaves.GetFree()

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

func HandleDeleteFile(m *structs.Master, masterList []string) http.HandlerFunc {
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

		m.Namespace.DelFile(filename)

		if len(masterList)>0{
			// send this information to all the other masters
			for _,masterip := range masterList{
				// req, err := http.NewRequest("POST", "http://"+masterip+"/master/modifynamespace", bytes.NewBuffer(filesBytes))
				// if err != nil {
				// 	log.Fatal(err)
				// }
			
				// client := &http.Client{
				// 	Timeout: time.Second * config.TIMEOUT,
				// }
			
				// resp, err := client.Do(req)
				// if err != nil || resp.StatusCode != 200 {
				// 	// log.Fatal("Failed to register with master.",masterip)
				// 	log.Println("Failed to register with master.",masterip)
				// } else {
				// 	fmt.Println("Successfully registered with master.",masterip)
				// 	checkReg[id] = 1
				// }
				fmt.Println("send modified namespace to",masterip)
			}
		}
	}
}

func HandleListDir(m *structs.Master) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		req.ParseForm()
		path := req.Form["ls"][0]
		fmt.Println("Path", path)

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
		fmt.Println(slaveIP,"slave is registered.")
		delete(files, "/"+slaveIP)
		m.Slaves.NewSlave(slaveIP, 0, files)
	}
}

func MasterHandleNamespace(m *structs.Master) http.HandlerFunc{
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Sending Namespace over...", m.Namespace.namespace)
		data, err := json.Marshal(m.Namespace.namespace)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(data)
	}
}

func MasterModifyHandleNamespace(m *structs.Master) http.HandlerFunc{
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Receive new namespace...")
		filesBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		files := make(map[string]string)
		err = json.Unmarshal(filesBytes, &files)
		if err != nil {
			log.Fatal(err)
		}
		// change namespace
		m.Namespace.rwLock.Lock()
		m.Namespace.namespace = files
		m.Namespace.rwLock.unlock()
		// reply back
		reply, err := json.Marshal("OKAY")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(reply)
	}
}
