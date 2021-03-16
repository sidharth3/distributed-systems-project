package senders

import (
	"ds-proj/master/structs"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func deadSlave(m *structs.Master, slave *structs.Slave) {
	m.Lock.Lock()
	delete(m.Slaves, slave)
	for file := range slave.Files {
		delete(m.DirectoryTable[file], slave)
	}
	m.Lock.Unlock()
}

func aliveSlave(m *structs.Master, slave *structs.Slave, files map[string]bool) {
	m.Lock.Lock()
	// Maybe add some checksum to check whether any changes, then can skip the update.
	for file := range slave.Files {
		delete(m.DirectoryTable[file], slave)
	}
	slave.Files = files
	for file := range slave.Files {
		if m.DirectoryTable[file] == nil {
			m.DirectoryTable[file] = make(map[*structs.Slave]bool)
		}
		m.DirectoryTable[file][slave] = true
	}
	m.Lock.Unlock()
}

func HeartbeatSender(m *structs.Master) {
	for {
		time.Sleep(time.Duration(5000) * time.Millisecond)
		if len(m.Slaves) != 0 {
			fmt.Println("Sending heartbeats ...")

			for slave := range m.Slaves {
				if slave.Status != structs.DEAD {
					go func(slave *structs.Slave) {
						fmt.Println("Connecting to slave at ", slave.IP, "...")
						req, err := http.NewRequest("GET", "http://"+slave.IP+"/heartbeat", nil)
						client := &http.Client{
							Timeout: time.Second * 5,
						}
						resp, err := client.Do(req)
						if err != nil || resp.StatusCode != 200 {
							fmt.Println("Error", err)
							fmt.Println(slave.IP, " is DEAD. Editing metadata.")
							deadSlave(m, slave)
							fmt.Println("Metadata edited.")
							fmt.Println(m)
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
			}
		}
	}
}
