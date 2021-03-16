package structs

import "sync"

// Map to bool represents a set. Easier to delete element.
type Master struct {
	IP             string
	Lock           *sync.Mutex
	Slaves         map[*Slave]bool
	DirectoryTable map[string](map[*Slave]bool)
}

type Slave struct {
	IP     string
	Files  map[string]bool
	Status Status
}

// Status is an enumerated type
type Status string

const (
	OVERLOADED  Status = "Current load exceeds threshold"
	UNDERLOADED Status = "Current load does not exceed threshold"
	DEAD        Status = "Cannot be pinged"
)
