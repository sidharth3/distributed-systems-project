package senders

import (
	"bytes"
	"ds-proj/slave/config"
	"ds-proj/slave/helpers"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
			
				client := &http.Client{
					Timeout: time.Second * config.TIMEOUT,
				}
				res, err := client.Post("POST", "http://"+masterip+"/register", "application/json", bytes.NewBuffer(filesBytes))
				if err != nil {
					log.Fatal(err)
				} else if res.StatusCode != http.StatusOK {
					log.Fatal(res)
				} else {
					fmt.Println("Successfully registered with master")
					checkReg[id] = 1
				}
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
   
