package structs

import "sync"

type Master struct {
	Slaves        *Slaves
	FileLocations *FileLocations
	Namespace     *Namespace
	GCCount       *GCCount
}

func InitMaster() *Master {
	slaves := &Slaves{&sync.RWMutex{}, make(map[*Slave]bool), make([]*Slave, 0)}
	fileLocations := &FileLocations{&sync.RWMutex{}, make(map[string]map[string]bool)}
	namespace := &Namespace{&sync.RWMutex{}, make(map[string]string)}
	gccount := &GCCount{&sync.RWMutex{}, make(map[string]int)}
	return &Master{slaves, fileLocations, namespace, gccount}
}

func (m *Master) UnlinkedHashes() map[string]bool {
	unlinked := make(map[string]bool)
	linked := m.Namespace.LinkedHashes()
	m.FileLocations.rwLock.RLock()
	for hash := range m.FileLocations.fileLocations {
		if !linked[hash] {
			unlinked[hash] = true
		}
	}
	defer m.FileLocations.rwLock.RUnlock()
	return unlinked
}

func (m *Master) UnlinkedNamespace() map[string]bool {
	unlinked := make(map[string]bool)
	m.Namespace.rwLock.RLock()
	m.FileLocations.rwLock.RLock()
	for filename, hash := range m.Namespace.namespace {
		if m.FileLocations.fileLocations[hash] == nil {
			unlinked[filename] = true
		}
	}
	defer m.FileLocations.rwLock.RUnlock()
	defer m.Namespace.rwLock.RUnlock()
	return unlinked
}
