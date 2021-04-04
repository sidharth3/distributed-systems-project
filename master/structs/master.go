package structs

import "sync"

type Master struct {
	Slaves        *Slaves
	FileLocations *FileLocations
	Namespace     *Namespace
	Queue         *OperationQueue
}

func InitMaster() *Master {
	slaves := &Slaves{&sync.RWMutex{}, make(map[*Slave]bool)}
	fileLocations := &FileLocations{&sync.RWMutex{}, make(map[string]map[string]bool)}
	namespace := &Namespace{&sync.RWMutex{}, make(map[string]string)}
	queue := &OperationQueue{&sync.RWMutex{}, make([]string, 0), make(map[string]*QueueItem)}
	return &Master{slaves, fileLocations, namespace, queue}
}

func (m *Master) Commit(uid string) {
	if m.Queue.FirstUID() == uid {
		commits := m.Queue.Dequeue()
		// Apply commits in queue order
		for _, item := range commits {
			if item.Hash == "delete" {
				m.Namespace.DelFile(item.Filename)
			} else {
				m.Namespace.SetHash(item.Filename, item.Hash)
			}
		}
	}
}
