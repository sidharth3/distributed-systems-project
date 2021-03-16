package test

import (
	"ds-proj/master/structs"
	"sync"
)

func Test_case1() *structs.Master {
	var lock sync.Mutex
	defaultFileArr := []string{"test_file.txt"}
	defaultSlave := &structs.Slave{"127.0.0.1:8081", defaultFileArr}
	defaultSlaveMap := make(map[*structs.Slave]structs.Status)
	defaultSlaveMap[defaultSlave] = structs.UNDERLOADED
	defaultDirTable := make(map[string][]*structs.Slave)
	defaultDirTable["test_file.txt"] = []*structs.Slave{defaultSlave}
	master := structs.Master{"127.0.0.1:8080", &lock, defaultSlaveMap, defaultDirTable}
	return &master
}
