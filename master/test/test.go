package test

import (
	"ds-proj/master/structs"
	"sync"
)

func TestCase1() *structs.Master {
	var lock sync.Mutex
	defaultFileSet := make(map[string]bool)
	defaultFileSet["test_file.txt"] = true
	defaultSlave := &structs.Slave{"127.0.0.1:8081", defaultFileSet, structs.UNDERLOADED}
	defaultSlaveMap := make(map[*structs.Slave]bool)
	defaultSlaveMap[defaultSlave] = true
	defaultDirTable := make(map[string](map[*structs.Slave]bool))
	defaultDirTableEntry := make(map[*structs.Slave]bool)
	defaultDirTableEntry[defaultSlave] = true
	defaultDirTable["test_file.txt"] = defaultDirTableEntry
	master := structs.Master{"127.0.0.1:8080", &lock, defaultSlaveMap, defaultDirTable}
	return &master
}

func EmptyCase() *structs.Master {
	master := structs.Master{"127.0.0.1:8080", &sync.Mutex{}, make(map[*structs.Slave]bool), make(map[string]map[*structs.Slave]bool)}
	return &master
}
