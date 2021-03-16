package main

import (
	"ds-proj/master/handlers"
	"ds-proj/master/structs"
	"net/http"
)

func main() {
	defaultFileArr := []string{"test_file.txt"}
	defaultSlave := &structs.Slave{"127.0.0.1:8081", 0, defaultFileArr}
	defaultSlaveMap := make(map[*structs.Slave]structs.Status)
	defaultSlaveMap[defaultSlave] = structs.UNDERLOADED
	defaultDirTable := make(map[string][]*structs.Slave)
	defaultDirTable["test_file.txt"] = []*structs.Slave{defaultSlave}
	master := structs.Master{"127.0.0.1:8080", defaultSlaveMap, defaultDirTable}

	http.HandleFunc("/file", handlers.HandleFile(&master))
	http.ListenAndServe("127.0.0.1:8080", nil)
}
