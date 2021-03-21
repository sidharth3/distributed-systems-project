package test

import (
	"ds-proj/master/structs"
	"sync"
)

// Pure empty init
func EmptyCase() *structs.Master {
	master := structs.Master{"127.0.0.1:8080",
		&sync.Mutex{},
		&sync.Mutex{},
		&sync.Mutex{},
		make(map[*structs.Slave]bool),
		make(map[string]map[string]bool),
		make(map[string]string),
	}
	return &master
}

// Initialize with single file in namespace
func SimpleCase() *structs.Master {
	namespace := make(map[string]string)
	namespace["test_file.txt"] = "d383caabf6289b8ad52e401dafb20fb301ec3b760d1708e2501e5a39f130a1fc"
	master := structs.Master{"127.0.0.1:8080",
		&sync.Mutex{},
		&sync.Mutex{},
		&sync.Mutex{},
		make(map[*structs.Slave]bool),
		make(map[string]map[string]bool),
		namespace,
	}
	return &master
}
