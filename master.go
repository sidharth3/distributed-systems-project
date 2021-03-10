package main
import (
	"fmt"
	"net/http"
	"bytes"
	"sync"
	"time"
	"io/ioutil"
	"encoding/json"
)

type Message struct {
	Type string
	Body string
}

var memory map[string] []string
// var slaveNodes []string
var slaveNodes map[string] int

func ping(slaveLock *sync.Mutex, slaveNodes *map[string] int){
	for {
		if len(*slaveNodes)!=0{
			fmt.Println("Start pinging ...")
			time.Sleep(time.Duration(5000) * time.Millisecond)

			for slave, _ := range *slaveNodes{
				fmt.Println(slave)
				go func(slave string) {
					msg := []byte("Alive")
				
					fmt.Println("Connecting to slave at ", slave, "...\n")
					req, err := http.NewRequest("POST", slave, bytes.NewBuffer(msg))
					client := &http.Client{
						Timeout: time.Second * 5,
					}

					resp, err := client.Do(req)
					if err != nil {
						// panic(err)
						fmt.Println("Error",err)
						
						slaveLock.Lock()
						delete(*slaveNodes, slave)
						fmt.Println(slave," is removed. New slaveNodes: ",*slaveNodes)
						slaveLock.Unlock()
					}else{
						fmt.Println("Response Status:", resp.Status)
						fmt.Println("Response Headers:", resp.Header)
						body, _ := ioutil.ReadAll(resp.Body)
						fmt.Println("Response Body:", string(body))
						
						if resp.Status == "404 Not Found" {
							slaveLock.Lock()
							delete(*slaveNodes, slave)
							fmt.Println(slave," is removed. New slaveNodes: ",*slaveNodes)
							slaveLock.Unlock()
						}
					}
					// defer resp.Body.Close()
					
				}(slave)
			}
		}
	}
}

func conn (memLock *sync.Mutex, slaveLock *sync.Mutex, slaveNodes *map[string] int) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request){
		fmt.Println("Handling a connection ...")
			
		buf := new(bytes.Buffer)
		buf.ReadFrom(req.Body)
		// buf2 := buf.Bytes()

		var msg Message
			
		if err := json.Unmarshal(buf.Bytes(), &msg); err != nil {
			panic(err)
		}
		// fmt.Println("Directory: ", location)
		fmt.Println("Message received",msg)
		fmt.Println(msg.Type)

		if msg.Type == "Register"{
			slaveLock.Lock()
			(*slaveNodes)[msg.Body] = 1
			fmt.Println("Slave node",msg.Body," is registered. New slaveNodes ", *slaveNodes)
			slaveLock.Unlock()
		}else if msg.Type == "Directory"{
			memLock.Lock()
			if slaveNodes, ok := memory[msg.Body]; ok {
				for _, slave := range slaveNodes {
					w.Write([]byte(slave+"\n"))
				}
			} else {
				fmt.Fprintf(w, "Slave node not found ...")
			}
			memLock.Unlock()
		}

	}
}

func main(){
	var memLock sync.Mutex
	var slaveLock sync.Mutex
	memory = make(map[string][]string)
	memory["/foo/bar.txt"] = []string{"http://127.0.0.1:8090/1","http://127.0.0.1:8090/2","http://127.0.0.1:8090/3"}
	
	slaveNodes = make(map[string] int)
	// slaveNodes["http://127.0.0.1:8090/1"] = 1
	// slaveNodes["http://127.0.0.1:8090/2"] = 1

	go Listen(&memLock, &slaveLock, &slaveNodes)

	go ping(&slaveLock, &slaveNodes)

	for{}
}

func Listen(memLock *sync.Mutex, slaveLock *sync.Mutex, slaveNodes *map[string] int){
	fmt.Println("Starting server ...")
	http.HandleFunc("/",conn(memLock, slaveLock, slaveNodes))
	http.ListenAndServe("127.0.0.1:8080", nil)
	fmt.Println("Listening on :8080 ...")
}