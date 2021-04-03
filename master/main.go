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
<<<<<<< HEAD
	go periodic.SlaveGarbageCollector(master)
	go periodic.CheckReplica(master)
	go periodic.DeleteUidFromQueue(master)
=======
	go periodic.DeleteUidFromQueue(master)
	go periodic.SlaveGarbageCollector(master)
	go periodic.CheckReplica(master)
>>>>>>> 23718734330cfc879c3e2cce57dc9c6dc6f21b81

	http.HandleFunc("/file", handlers.HandleFile(master))
	http.HandleFunc("/delete", handlers.HandleDeleteFile(master))
	http.HandleFunc("/slaveips", handlers.HandleSlaveIPs(master))
	http.HandleFunc("/update", handlers.HandleUpdate(master))
	http.HandleFunc("/register", handlers.HandleNewSlave(master))
	http.ListenAndServe("127.0.0.1:8080", nil)
}
