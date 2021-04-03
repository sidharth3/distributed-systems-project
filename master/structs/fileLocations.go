package structs

import (
	"ds-proj/master/config"
	"sync"
)

// asdf332789asfj -> {ip1, ip2, ip3}, master periodically updates this based on Slaves
type FileLocations struct {
	rwLock        *sync.RWMutex
	fileLocations map[string]map[string]bool
}

func (f *FileLocations) GetIPs(hash string) []string {
	ips := make([]string, 0)
	f.rwLock.RLock()
	for ip := range f.fileLocations[hash] {
		ips = append(ips, ip)
	}
	defer f.rwLock.RUnlock()
	return ips
}

func (f *FileLocations) Remake(newFileLocations map[string]map[string]bool) {
	f.rwLock.Lock()
	f.fileLocations = newFileLocations
	f.rwLock.Unlock()
}

func (f *FileLocations) NeedReplication() map[string]map[string]bool {
	needReplication := make(map[string]map[string]bool)
	f.rwLock.RLock()
	for hash, ips := range f.fileLocations {
		if len(ips) < config.REPLICAS {
			needReplication[hash] = make(map[string]bool)
			for ip := range ips {
				needReplication[hash][ip] = true
			}
		}
	}
	defer f.rwLock.RUnlock()
	return needReplication
}
