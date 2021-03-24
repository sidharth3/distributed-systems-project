package periodic

import (
	"bytes"
	"ds-proj/master/config"
	"ds-proj/master/structs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func deadSlave(m *structs.Master, slave *structs.Slave) {
	m.SLock.Lock()
	delete(m.Slaves, slave)
	m.SLock.Unlock()
}

func aliveSlave(m *structs.Master, slave *structs.Slave, files map[string]bool) {
	m.SLock.Lock()
	slave.Files = files
	m.SLock.Unlock()
}

func HeartbeatSender(m *structs.Master) {
	for {
		time.Sleep(time.Duration(config.HBINTERVAL) * time.Second)
		fmt.Println("Sending heartbeats ...")
		m.SLock.Lock()
		for slave := range m.Slaves {
			go func(slave *structs.Slave) {
				// No need to lock because Slave IP won't change
				fmt.Println("Connecting to slave at ", slave.IP, "...")
				req, err := http.NewRequest("GET", "http://"+slave.IP+"/heartbeat", nil)
				client := &http.Client{
					Timeout: time.Second * config.TIMEOUT,
				}
				resp, err := client.Do(req)
				if err != nil || resp.StatusCode != 200 {
					fmt.Println(slave.IP, " is DEAD. Editing metadata.")
					deadSlave(m, slave)
					fmt.Println("Metadata edited.")
				} else {
					fmt.Println(slave.IP + " is alive. Updating metadata.")
					filesBytes, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Fatal(err)
					}

					files := make(map[string]bool)
					err = json.Unmarshal(filesBytes, &files)
					if err != nil {
						log.Fatal(err)
					}

					aliveSlave(m, slave, files)
					fmt.Println("Metadata updated.")
				}
			}(slave)
		}
		m.SLock.Unlock()
	}
}

func updateFileLocations(m *structs.Master) {
	m.SLock.Lock()
	updatedFileLocations := make(map[string]map[string]bool)
	for slave := range m.Slaves {
		for hash := range slave.Files {
			if updatedFileLocations[hash] == nil {
				updatedFileLocations[hash] = make(map[string]bool)
			}
			updatedFileLocations[hash][slave.IP] = true
		}
	}
	m.FLock.Lock()
	m.FileLocations = updatedFileLocations
	m.FLock.Unlock()
	m.SLock.Unlock()
}

func FileLocationsUpdater(m *structs.Master) {
	for {
		time.Sleep(time.Second * config.FLINTERVAL)
		fmt.Println("Updating file locations")
		updateFileLocations(m)
		fmt.Println("File locations updated")
	}
}

// for garbage collector: 
// periodically sends over the values of the namespaces in the Master struct
func SlaveGarbageCollector(m *structs.Master){
	for{
		time.Sleep(time.Duration(config.GCINTERVAL) * time.Second)
		fmt.Println("Sending garbage collection message ...")
		// prepare hashedContent
		m.NLock.Lock()
		hashedContent := make(map[string]bool)
		for _,v := range m.Namespace{
			hashedContent[v] = true
		}
		m.NLock.Unlock()

		filesBytes, err := json.Marshal(hashedContent)
		if err != nil {
			log.Fatal()
		}

		//send over hashedContent to each slave
		m.SLock.Lock()
		for slave := range m.Slaves {
			go func(slave *structs.Slave) {
				fmt.Println("Sending garbage collector msg to slave at", slave.IP, "...")
				req, err := http.NewRequest("POST", "http://"+slave.IP+"/garbagecollector", bytes.NewBuffer(filesBytes))
				if err != nil {
					log.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				client := &http.Client{
					Timeout: time.Second * config.TIMEOUT,
				}
				client.Do(req)
			}(slave)
		}
		m.SLock.Unlock()
	}
}