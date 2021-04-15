package handlers

import (
	"bytes"
	"ds-proj/master/config"
	"ds-proj/master/periodic"
	"ds-proj/master/structs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

// Sends an array of strings over to the client. [ip1, ip2, ip3]
func HandleSlaveIPs(m *structs.Master, masterList []string) http.HandlerFunc {

	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("handleslaveip")
		firstBecomeMaster(m, masterList)
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
func HandleFile(m *structs.Master, masterList []string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("HandleFile")
		firstBecomeMaster(m, masterList)
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
		fmt.Println("HandleDeleteFile")
		firstBecomeMaster(m, masterList)
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

		filenameBytes2, err := json.Marshal(filename)
		if err != nil {
			log.Fatal()
		}

		var wg sync.WaitGroup
		numofreplies := 0
		var numofreplieslock sync.Mutex
		for _, masterip := range masterList {
			wg.Add(1) //(len(masterList) -1)/ 2
			go masterSendForReply(masterip, filenameBytes2, "delfile", &wg, &numofreplies, &numofreplieslock)
		}
		wg.Wait()
		// check for majority
		if numofreplies >= (len(masterList) -1)/ 2{
			status := "DONE"
			fmt.Println("Reply from majority received")
		}else{
			status := "NOTDONE"
			fmt.Println("NOT enough reply from majority received")
		}

		// send status to client
		msgtoclient, err := json.Marshal(status)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(msgtoclient)
	}
}

func HandleListDir(m *structs.Master, masterList []string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("HandleListDir")

		firstBecomeMaster(m, masterList)
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
		fmt.Println(slaveIP, "slave is registered.")
		delete(files, "/"+slaveIP)
		m.Slaves.NewSlave(slaveIP, 0, files)
	}
}

func collateNS(masterip string, count map[[2]string]int, countLock *sync.Mutex, wg *sync.WaitGroup) {
	client := &http.Client{
		Timeout: time.Second * config.TIMEOUT,
	}

	resp, err := client.Post("http://"+masterip+"/master/namespace", "application/json", nil)
	if err != nil || resp.StatusCode != 200 {
		log.Println("Failed to reach master for namespace request at /master/namespace.", masterip)
		wg.Done()
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var incomingNS map[string]string
	err = json.Unmarshal(body, &incomingNS)
	if err != nil {
		log.Fatal(err)
	}

	for fn, hash := range incomingNS {
		incomingkey := [2]string{fn, hash}
		c, ok := count[incomingkey]
		if ok {
			newcount := c + 1
			count[incomingkey] = newcount
			// check if newcount reaches majority

		} else {
			count[incomingkey] = 1
		}
	}
	wg.Done()
}

// for master replica -----------------------------------

func firstBecomeMaster(m *structs.Master, masterList []string) {
	m.IsPrimaryLock.Lock()
	fmt.Println(m.IsPrimary)
	if !m.IsPrimary {
		fmt.Println("This is the new primary master.")
		// [filename, hash]: count
		count := make(map[[2]string]int)
		var countLock sync.Mutex
		//first populate this with your current namespace
		currentNS := m.Namespace.ReturnNamespace()
		for fn, hash := range currentNS {
			key := [2]string{fn, hash}
			count[key] = 1
		}
		// do collation to the other namespaces
		var wg sync.WaitGroup
		for _, masterip := range masterList {
			wg.Add(1)
			go collateNS(masterip, count, &countLock, &wg)
		}
		wg.Wait()
		// collating majority entries
		updatedNS := make(map[string]string)
		for key, c := range count {
			if c > len(masterList)/2+1 {
				updatedNS[key[0]] = key[1]
			}
		}
		m.Namespace.SetNamespace(updatedNS)
		// spawn the periodic gorountines
		go periodic.HeartbeatSender(m)
		go periodic.LoadChecker(m)
		go periodic.FileLocationsUpdater(m)
		go periodic.SlaveGarbageCollector(m)
		go periodic.CheckReplica(m)
		go periodic.MasterGarbageCollector(m)
		// change the IsPrimary to true
		m.IsPrimary = true
	}
	m.IsPrimaryLock.Unlock()
}

func masterSendForReply(masterip string, filenameBytes []byte, endpoint string, wg *sync.WaitGroup, numofreplies *int, &numofreplieslock *sync.Mutex) {
	fmt.Println("send reply to endpoint", endpoint, "to", masterip)
	req, err := http.NewRequest("POST", "http://"+masterip+"/master/"+endpoint, bytes.NewBuffer(filenameBytes)) //send over the string filename
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Timeout: time.Second * config.TIMEOUT,
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Println("Failed to reach master for reply.", masterip)
		wg.Done()
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var reply string
	err = json.Unmarshal(body, &reply)
	if err != nil {
		log.Fatal(err)
	}
	if reply == "OKAY" { // check that the reply is OKAY
		fmt.Println("Successfully get reply from master.", masterip)
		numofreplies.Lock()
		*numofreplies++
		numofreplies.Unlock()
		wg.Done()
	}
}

func MasterHandleNamespace(m *structs.Master) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		newNS := m.Namespace.ReturnNamespace()
		fmt.Println("Sending Namespace over...", newNS)
		data, err := json.Marshal(newNS)
		if err != nil {
			log.Fatal(err)
		}
		w.Write(data)
	}
}

func MasterHandleDelFile(m *structs.Master) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println("Receive filename to delete...")
		filesBytes, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		var filename string
		err = json.Unmarshal(filesBytes, &filename)
		if err != nil {
			log.Fatal(err)
		}
		// del filename from namespace
		m.Namespace.DelFile(filename)

		// reply back
		reply, err := json.Marshal("OKAY")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(reply)
	}
}

// HELPER FUNCTIONS ---------------------------------
func sum(array []int) int {
	result := 0
	for _, v := range array {
		result += v
	}
	return result
}
