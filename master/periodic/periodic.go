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

func HeartbeatSender(m *structs.Master) {
	for {
		time.Sleep(time.Duration(config.HBINTERVAL) * time.Second)
		fmt.Println("Sending heartbeats ...")
		f := func(slave *structs.Slave) {
			ip := slave.GetIP()
			fmt.Println("Connecting to slave at ", ip, "...")
			req, err := http.NewRequest("GET", "http://"+ip+"/heartbeat", nil)
			client := &http.Client{
				Timeout: time.Second * config.TIMEOUT,
			}
			resp, err := client.Do(req)
			if err != nil || resp.StatusCode != 200 {
				fmt.Println(ip, " is DEAD. Updating metadata.")
				m.Slaves.DelSlave(slave)
				fmt.Println("Metadata edited.")
			} else {
				fmt.Println(ip + " is alive. Updating metadata.")
				filesBytes, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}

				files := make(map[string]bool)
				err = json.Unmarshal(filesBytes, &files)
				if err != nil {
					log.Fatal(err)
				}
				slave.SetHashes(files)
				fmt.Println("Metadata updated.")
			}
		}

		m.Slaves.ForEvery(f)
	}
}

func FileLocationsUpdater(m *structs.Master) {
	for {
		time.Sleep(time.Second * config.FLINTERVAL)
		fmt.Println("Updating file locations")
		newFileLocations := m.Slaves.GenFileLocations()
		m.FileLocations.Remake(newFileLocations)
		fmt.Println("File locations updated")
	}
}

// for garbage collector:
// periodically sends over the values of the namespaces in the Master struct
func SlaveGarbageCollector(m *structs.Master) {
	for {
		time.Sleep(time.Duration(config.GCINTERVAL) * time.Second)
		fmt.Println("Sending garbage collection message ...")
		// prepare hashedContent
		hashedContent := m.Namespace.LinkedHashes()

		filesBytes, err := json.Marshal(hashedContent)
		if err != nil {
			log.Fatal()
		}

		//send over hashedContent to each slave
		f := func(slave *structs.Slave) {
			ip := slave.GetIP()
			fmt.Println("Sending garbage collector msg to slave at", ip, "...")
			req, err := http.NewRequest("POST", "http://"+ip+"/garbagecollector", bytes.NewBuffer(filesBytes))
			if err != nil {
				log.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			client := &http.Client{
				Timeout: time.Second * config.TIMEOUT,
			}
			client.Do(req)
		}
		m.Slaves.ForEvery(f)
	}
}

func CheckReplica(m *structs.Master) {
	for {
		time.Sleep(time.Duration(config.REPINTERVAL) * time.Second)
		fmt.Println("Replication cycle starting")
		toReplicate := make(map[string]map[string]string) // {slaveip: {fileHash: ip1, fileHash2: ip2}}

		for f, slaveips := range m.FileLocations.NeedReplication() {
			replicasLeft := config.REPLICAS - len(slaveips)
			for slaveip := range slaveips {
				for _, ip := range m.Slaves.FreeForReplication(f, replicasLeft) {
					if toReplicate[ip] == nil {
						toReplicate[ip] = make(map[string]string)
					}
					toReplicate[ip][f] = slaveip
				}
				break // Just getting the first slave from the map
			}
		}

		// Slave get replicas => {fileHash: ip1, fileHash, ip2}
		for slaveip, toGet := range toReplicate {
			slaveURL := "http://" + slaveip + "/replica"
			jsonReq, err := json.Marshal(toGet)
			if err != nil {
				log.Fatal(err)
			}
			req, err := http.NewRequest("POST", slaveURL, bytes.NewBuffer(jsonReq))
			client := &http.Client{
				Timeout: time.Second * config.TIMEOUT,
			}
			client.Do(req)
		}
	}
}
