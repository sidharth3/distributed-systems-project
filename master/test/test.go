package test

import (
	"ds-proj/master/structs"
	"sync"
)

// Pure empty init
func EmptyCase() *structs.Master {
	var qItem []structs.QueueItem

	queue := structs.OperationQueue{
		Queue: qItem,
		QLock: &sync.RWMutex{},
	}

	master := structs.Master{IP: "127.0.0.1:8080",
		SLock:         &sync.Mutex{},
		FLock:         &sync.Mutex{},
		NLock:         &sync.Mutex{},
		Slaves:        make(map[*structs.Slave]bool),
		FileLocations: make(map[string]map[string]bool),
		Namespace:     make(map[string]string),
		Queue:         &queue,
	}
	return &master
}

// Initialize with single file in namespace
func SimpleCase() *structs.Master {
	namespace := make(map[string]string)
	namespace["test_file.txt"] = "d383caabf6289b8ad52e401dafb20fb301ec3b760d1708e2501e5a39f130a1fc"
	var qItem []structs.QueueItem

	queue := structs.OperationQueue{
		Queue: qItem,
		QLock: &sync.RWMutex{},
	}
	// uid := uuid.NewString()
	// queue.Enqueue(uid)
	master := structs.Master{IP: "127.0.0.1:8080",
		SLock:         &sync.Mutex{},
		FLock:         &sync.Mutex{},
		NLock:         &sync.Mutex{},
		Slaves:        make(map[*structs.Slave]bool),
		FileLocations: make(map[string]map[string]bool),
		Namespace:     make(map[string]string),
		Queue:         &queue,
	}
	return &master
}
