package main

import (
	"ds-proj/master/handlers"
	"ds-proj/master/periodic"
	"ds-proj/master/test"
	"fmt"
	"net/http"
	"os"
)

//to run master - go run main.go 8080 8081 8082 -p
//to run master - go run main.go 8081 8080 8082
//to run master - go run main.go 8082 8080 8081

func main() {
	master := test.SimpleCase()

	// check if it is primary
	if len(os.Args) > 3 && os.Args[len(os.Args)-1] == "-p" {
		fmt.Println("This is the initial primary master.", "127.0.0.1:"+os.Args[1])
		go periodic.HeartbeatSender(master)
		go periodic.LoadChecker(master)
		go periodic.FileLocationsUpdater(master)
		go periodic.SlaveGarbageCollector(master)
		go periodic.CheckReplica(master)
		go periodic.MasterGarbageCollector(master)
		master.IsPrimary = true
	} else {
		fmt.Println("This is a master.", "127.0.0.1:"+os.Args[1])
	}

	// array of all other masters
	masterList := make([]string, 0)
	for i := 2; i < len(os.Args); i++ {
		if os.Args[i] != "-p" {
			masterList = append(masterList, "127.0.0.1:"+os.Args[i])
		}
	}

	http.HandleFunc("/file", handlers.HandleFile(master, masterList))
	http.HandleFunc("/delete", handlers.HandleDeleteFile(master, masterList))
	http.HandleFunc("/ls", handlers.HandleListDir(master, masterList))
	http.HandleFunc("/slaveips", handlers.HandleSlaveIPs(master, masterList))
	http.HandleFunc("/register", handlers.HandleNewSlave(master))

	// for other masters
	http.HandleFunc("/master/namespace", handlers.MasterHandleNamespace(master)) //tosend over namespaces
	http.HandleFunc("/master/delfile", handlers.MasterHandleDelFile(master))     //to del something in the namespace and sends reply back

	// http.ListenAndServe("127.0.0.1:8080", nil)
	http.ListenAndServe("127.0.0.1:"+os.Args[1], nil)
}
