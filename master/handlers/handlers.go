package handlers

import (
	"bytes"
	"ds-proj/master/config"
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
	firstBecomeMaster(m, masterList)

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
func HandleFile(m *structs.Master, masterList []string) http.HandlerFunc {
	firstBecomeMaster(m,, masterList)
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
	firstBecomeMaster(m, masterList)
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

		filenameBytes2, err := json.Marshal(filename)
		if err != nil {
			log.Fatal()
		}

		var wg sync.WaitGroup
		wg.Add(len(masterList)/2+1)
		for _,masterip := range masterList{
			go masterSendForReply(masterip,filenameBytes2,"delfile", &wg)
		}
		wg.Wait()

		// checkreply := make([]int, len(masterList))
		// for i :=0 ; i<len(masterList);i++{
		// 	checkreply = append(checkreply,0)
		// }

		// var checkreplyLock sync.RWMutex
		// for id,masterip := range masterList{
		// 	go masterSendForReply(id, masterip, &checkreply, &checkreplyLock,filenameBytes2,"delfile")
		// }

		// for sum(checkreply)<len(masterList)/2+1{ //blocking -- wait for majority of replies
		// 	fmt.Println("Waiting for Reply from majority")
		// } 
		fmt.Println("Reply from majority received")

		// send DONE to client
		msgtoclient, err := json.Marshal("DONE")
		if err != nil {
			log.Fatal(err)
		}
		w.Write(msgtoclient)
	}
}

func HandleListDir(m *structs.Master, masterList []string) http.HandlerFunc {
	firstBecomeMaster(m, masterList)
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

func collateNS(masterip string,  count *map[[]string]int, countLock *sync.Mutex, wg *sync.WaitGroup){
	req, err := http.NewRequest("POST", "http://"+masterip+"/master/namespace") 
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Timeout: time.Second * config.TIMEOUT,
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Println("Failed to reach master for namespace request at /master/namespace.",masterip)
		wg.Done()
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var incomingNS map[string] string
	err = json.Unmarshal(body, &incomingNS)
	if err != nil {
		log.Fatal(err)
	}

	for fn, hash:= range incomingNS{
		incomingkey := []string{fn,hash}
		count, ok := (*count)[incomingkey]
		if ok{
			newcount := count+1
			(*count)[incomingkey] = newcount 
			// check if newcount reaches majority

		}else{
			(*count)[incomingkey] = 1
		}
	}
	wg.Done()
}

// for master replica -----------------------------------

func firstBecomeMaster(m *structs.Master, masterList []string){
	m.isPrimaryLock.Lock()
	if !*m.isPrimary{
		fmt.Println("This is the new primary master.")
		// [filename, hash]: count
		count:= make(map[[]string] int)
		var countLock sync.Mutex
		//first populate this with your current namespace
		currentNS := m.Namespace.ReturnNamespace()
		for fn, hash:= range currentNS{
			key:=[]string{fn,hash}
			count[key] = 1
		}
		// do collation to the other namespaces
		var wg sync.WaitGroup
		for _,masterip :=range masterList{
			wg.Add(1)
			go collateNS(masterip, &count, &countLock,&wg)
		}
		wg.Wait()
		// collating majority entries
		updatedNS := make(map [string]string)
		for key, c:= range count{
			if c > len(masterList)/2+1{
				updatedNS[key[0]] = key[1]
			}
		}
		m.Namespace.setNamespace(updatedNS)
		// spawn the periodic gorountines
		go periodic.HeartbeatSender(m)
		go periodic.LoadChecker(m)
		go periodic.FileLocationsUpdater(m)
		go periodic.SlaveGarbageCollector(m)
		go periodic.CheckReplica(m)
		go periodic.MasterGarbageCollector(m)
		// change the isPrimary to true
		*m.isPrimary = true
	} 	
	m.isPrimaryLock.Unlock()
}

func masterSendForReply(masterip string, filenameBytes []byte, endpoint string, wg *sync.WaitGroup){
	fmt.Println("send reply to endpoint",endpoint,"to",masterip)
	req, err := http.NewRequest("POST", "http://"+masterip+"/master/"+endpoint, bytes.NewBuffer(filenameBytes)) //send over the string filename
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Timeout: time.Second * config.TIMEOUT,
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Println("Failed to reach master for reply.",masterip)
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
	if reply =="REPLY"{ // check that the reply is OKAY
		fmt.Println("Successfully get reply from master.",masterip)
		// checkreplyLock.Lock()
		// (*checkreply)[id] = 1
		// checkreplyLock.Unlock()
		wg.Done()
	}
}

func MasterHandleNamespace(m *structs.Master) http.HandlerFunc{
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

func MasterHandleDelFile(m *structs.Master) http.HandlerFunc{
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
