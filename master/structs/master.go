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
	queue := &OperationQueue{make([]QueueItem, 0), &sync.RWMutex{}}
	return &Master{slaves, fileLocations, namespace, queue}
}
