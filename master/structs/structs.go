package structs

import "sync"

type Master struct {
	IP             string
	Lock           *sync.Mutex
	Slaves         map[*Slave]Status
	DirectoryTable map[string][]*Slave
}

type Slave struct {
	IP    string
	Files []string
}

// Status is an enumerated type
type Status string

const (
	OVERLOADED  Status = "Current load exceeds threshold"
	UNDERLOADED Status = "Current load does not exceed threshold"
	DEAD        Status = "Cannot be pinged"
)
