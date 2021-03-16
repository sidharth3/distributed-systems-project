package senders

import (
	"ds-proj/master/structs"
	"fmt"
	"net/http"
	"time"
)

func HeartbeatSender(m *structs.Master) {
	for {
		if len(m.Slaves) != 0 {
			fmt.Println("Sending heartbeats ...")
			time.Sleep(time.Duration(5000) * time.Millisecond)

			for slave, status := range m.Slaves {
				if status != structs.DEAD {
					go func(slave *structs.Slave) {
						fmt.Println("Connecting to slave at ", slave.IP, "...")
						req, err := http.NewRequest("GET", "http://"+slave.IP+"/heartbeat", nil)
						client := &http.Client{
							Timeout: time.Second * 5,
						}
						resp, err := client.Do(req)
						if err != nil || resp.StatusCode != 200 {
							fmt.Println("Error", err)
							m.Lock.Lock()
							m.Slaves[slave] = structs.DEAD
							fmt.Println(slave.IP, " is DEAD.")
							m.Lock.Unlock()
						} else {
							fmt.Println("Success")
						}
						// else {
						// 	if resp.StatusCode == 404 {

						// 	}
						// 	fmt.Println("Response Status:", resp.Status)
						// 	fmt.Println("Response Headers:", resp.Header)
						// 	body, _ := ioutil.ReadAll(resp.Body)
						// 	fmt.Println("Response Body:", string(body))

						// 	if resp.Status == "404 Not Found" {
						// 		slaveLock.Lock()
						// 		delete(*slaveNodes, slave)
						// 		fmt.Println(slave, " is removed. New slaveNodes: ", *slaveNodes)
						// 		slaveLock.Unlock()
						// 	}
						// }
						// defer resp.Body.Close()

					}(slave)
				}
			}

		}
	}
}
