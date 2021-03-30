package main

import (
	"ds-proj/master/handlers"
	"ds-proj/master/periodic"
	"ds-proj/master/test"
	"net/http"
)

func main() {
	master := test.SimpleCase()

	go periodic.HeartbeatSender(master)
	go periodic.FileLocationsUpdater(master)
	go periodic.SlaveGarbageCollector(master)
	go periodic.CheckReplica(master)

	http.HandleFunc("/file", handlers.HandleFile(master))
	http.HandleFunc("/delete", handlers.HandleDeleteFile(master))
	http.HandleFunc("/slaveips", handlers.HandleSlaveIPs(master))
	http.HandleFunc("/register", handlers.HandleNewSlave(master))
	http.ListenAndServe("127.0.0.1:8080", nil)
}
