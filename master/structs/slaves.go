package structs

import (
	"ds-proj/master/config"
	"sync"
)

type Slaves struct {
	rwLock *sync.RWMutex
	slaves map[*Slave]bool
}

// Updated whenever master receives heartbeat from slave
type Slave struct {
	rwLock *sync.RWMutex
	ip     string
	status Status
	hashes map[string]bool // hashes that a slave has
}

// Status is an enumerated type
type Status string

const (
	OVERLOADED  Status = "Current load exceeds threshold"
	UNDERLOADED Status = "Current load does not exceed threshold"
)

// TODO: some way to select the most free slaves
func (s *Slaves) GetFree() []string {
	ips := make([]string, 0)
	s.rwLock.RLock()
	for slave := range s.slaves {
		ips = append(ips, slave.GetIP())
		if len(ips) >= config.REPLICAS {
			break
		}
	}
	defer s.rwLock.RUnlock()
	return ips
}

func (s *Slaves) NewSlave(ip string, status Status, hashes map[string]bool) {
	newSlave := &Slave{&sync.RWMutex{}, ip, status, hashes}
	s.rwLock.Lock()
	s.slaves[newSlave] = true
	s.rwLock.Unlock()
}

func (s *Slaves) DelSlave(slave *Slave) {
	s.rwLock.Lock()
	delete(s.slaves, slave)
	s.rwLock.Unlock()
}

func (s *Slaves) ForEvery(f func(*Slave)) {
	s.rwLock.RLock()
	for slave := range s.slaves {
		go f(slave)
	}
	s.rwLock.RUnlock()
}

func (s *Slaves) GenFileLocations() map[string]map[string]bool {
	updatedFileLocations := make(map[string]map[string]bool)
	s.rwLock.RLock()
	for slave := range s.slaves {
		slave.rwLock.RLock()
		for hash := range slave.hashes {
			if updatedFileLocations[hash] == nil {
				updatedFileLocations[hash] = make(map[string]bool)
			}
			updatedFileLocations[hash][slave.ip] = true
		}
		slave.rwLock.RUnlock()
	}
	defer s.rwLock.RUnlock()
	return updatedFileLocations
}

// TODO: select which slave to replicate to
func (s *Slaves) FreeForReplication(hash string, numNeeded int) []string {
	ips := make([]string, 0)
	s.rwLock.RLock()
	for slave := range s.slaves {
		slave.rwLock.RLock()
		if !slave.hashes[hash] {
			ips = append(ips, slave.ip)
			numNeeded -= 1
		}
		slave.rwLock.RUnlock()
		if numNeeded == 0 {
			break
		}
	}
	return ips
}

func (s *Slave) GetIP() string {
	// Don't need to lock because IP won't change. But leaving code commented out.
	// s.rwLock.RLock()
	// defer s.rwLock.RUnlock()
	return s.ip
}

func (s *Slave) SetHashes(hashes map[string]bool) {
	s.rwLock.Lock()
	s.hashes = hashes
	s.rwLock.Unlock()
}
