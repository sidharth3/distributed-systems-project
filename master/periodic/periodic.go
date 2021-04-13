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
	"strconv"
	"time"
)

func HeartbeatSender(m *structs.Master) {
	for {
		time.Sleep(time.Duration(config.HBINTERVAL) * time.Second)
		fmt.Println("Sending heartbeats ...")
		f := func(slave *structs.Slave) {
			ip := slave.GetIP()
			client := &http.Client{
				Timeout: time.Second * config.TIMEOUT,
			}
			res, err := client.Get("http://" + ip + "/heartbeat")
			if err != nil || res.StatusCode != http.StatusOK {
				fmt.Println(ip, " is DEAD")
				m.Slaves.DelSlave(slave)
			} else {
				filesBytes, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Fatal(err)
				}

				files := make(map[string]bool)
				err = json.Unmarshal(filesBytes, &files)
				if err != nil {
					log.Fatal(err)
				}
				slave.SetHashes(files)
			}
		}
		m.Slaves.ForEvery(f)
	}
}

func LoadChecker(m *structs.Master) {
	for {
		time.Sleep(time.Second * config.LDINTERVAL)
		fmt.Println("Checking loads...")
		f := func(slave *structs.Slave) {
			ip := slave.GetIP()
			client := &http.Client{
				Timeout: time.Second * config.TIMEOUT,
			}
			res, err := client.Get("http://" + ip + "/load")

			if err != nil || res.StatusCode != http.StatusOK {
				fmt.Println(ip, " is DEAD")
				m.Slaves.DelSlave(slave)
			} else {
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					log.Fatal(err)
				}
				load, err := strconv.Atoi(string(body))
				if err != nil {
					log.Fatal(err)
				}
				slave.SetLoad(load)
			}
		}
		m.Slaves.SortLoad()
		m.Slaves.ForEvery(f)
	}
}

func FileLocationsUpdater(m *structs.Master) {
	for {
		time.Sleep(time.Second * config.FLINTERVAL)
		fmt.Println("Updating file locations")
		newFileLocations := m.Slaves.GenFileLocations()
		m.FileLocations.Remake(newFileLocations)
	}
}

// for garbage collector:
// periodically sends over the values of the namespaces in the Master struct
func SlaveGarbageCollector(m *structs.Master) {
	for {
		time.Sleep(time.Duration(config.GCINTERVAL) * time.Second)
		fmt.Println("Sending garbage collection message ...")
		// prepare hashedContent
		hashedContent := m.UnlinkedHashes()

		filesBytes, err := json.Marshal(hashedContent)
		if err != nil {
			log.Fatal(err)
		}

		//send over hashedContent to each slave
		f := func(slave *structs.Slave) {
			ip := slave.GetIP()
			client := &http.Client{
				Timeout: time.Second * config.TIMEOUT,
			}
			client.Post("http://"+ip+"/garbagecollector", "application/json", bytes.NewBuffer(filesBytes))
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
			client := &http.Client{
				Timeout: time.Second * config.TIMEOUT,
			}
			client.Post(slaveURL, "application/json", bytes.NewBuffer(jsonReq))
		}
	}
}

func MasterGarbageCollector(m *structs.Master) {

	for {
		time.Sleep(time.Second * config.MGCINTERVAL)
		fmt.Println("Master Garbage Collection Cycle Starting")
		unlinked := m.UnlinkedNamespace()
		unlinked = m.GCCount.Cycle(unlinked)
		m.Namespace.CollectGarbage(unlinked)
	}
}
