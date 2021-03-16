package structs

type Master struct {
	IP             string
	Slaves         map[*Slave]Status
	DirectoryTable map[string][]*Slave
}

type Slave struct {
	IP    string
	ID    int
	Files []string
}

// Status is an enumerated type
type Status string

const (
	OVERLOADED  Status = "Current load exceeds threshold"
	UNDERLOADED Status = "Current load does not exceed threshold"
	DEAD        Status = "Cannot be pinged"
)
