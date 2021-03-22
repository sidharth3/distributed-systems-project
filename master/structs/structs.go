package structs

import "sync"

// Map to bool represents a set. Easier to delete element.
type Master struct {
	IP            string
	SLock         *sync.Mutex
	FLock         *sync.Mutex
	NLock         *sync.Mutex
	Slaves        map[*Slave]bool            // updated every heartbeat
	FileLocations map[string]map[string]bool // asdf332789asfj -> {ip1, ip2, ip3}, master periodically updates this based on Slaves
	Namespace     map[string]string          // foo/bar.txt -> asdf332789asfj, purely controlled by client
}

type Slave struct {
	IP     string
	Status Status
	Files  map[string]bool // hashes that a slave has
}

// Status is an enumerated type
type Status string

const (
	OVERLOADED  Status = "Current load exceeds threshold"
	UNDERLOADED Status = "Current load does not exceed threshold"
)
