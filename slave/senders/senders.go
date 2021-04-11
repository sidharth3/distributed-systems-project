package senders

import (
	"bytes"
	"ds-proj/slave/config"
	"ds-proj/slave/helpers"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

func RegisterWithMaster() {
	dirList := helpers.ListDir()
	dirList["/"+helpers.IP()] = true // filnames would not have / character. Pass the return address.
	filesBytes, err := json.Marshal(dirList)
	if err != nil {
		log.Fatal()
	}
	numMasters := len(helpers.MasterIP())
	checkReg := make([]int, numMasters)
	for i :=0 ; i<numMasters;i++{
		checkReg = append(checkReg,0)
	}

	// keep looping until sum of checkReg == len(helpers.MasterIP())
	for sum(checkReg) != numMasters{
		for id,masterip := range helpers.MasterIP(){ // loop through and register to all masters
			if checkReg[id]!=1{
				req, err := http.NewRequest("POST", "http://"+masterip+"/register", bytes.NewBuffer(filesBytes))
				if err != nil {
					log.Fatal(err)
				}
			
				client := &http.Client{
					Timeout: time.Second * config.TIMEOUT,
				}
			
				resp, err := client.Do(req)
				if err != nil || resp.StatusCode != 200 {
					// log.Fatal("Failed to register with master.",masterip)
					log.Println("Failed to register with master.",masterip)
				} else {
					fmt.Println("Successfully registered with master.",masterip)
					checkReg[id] = 1
				}
				// defer resp.Body.Close()
			}		
		}
		time.Sleep(1 * time.Second)
	}
}

func ForceUpdateMaster(data url.Values) {
	//Questions - does forceUpdate need to send directory also or new upload information only?
	master_URL := "http://127.0.0.1:8080/update"
	res, err := http.PostForm(master_URL, data)
	fmt.Println(res.StatusCode)
	if err != nil || res.StatusCode != 200 {
		fmt.Println("File upload has failed.")
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
   
