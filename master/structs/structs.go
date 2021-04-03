package structs

import (
	"fmt"
	"sync"
)

// Map to bool represents a set. Easier to delete element.
type Master struct {
	IP            string
	SLock         *sync.Mutex
	FLock         *sync.Mutex
	NLock         *sync.Mutex
	Slaves        map[*Slave]bool            // updated every heartbeat
	FileLocations map[string]map[string]bool // asdf332789asfj -> {ip1, ip2, ip3}, master periodically updates this based on Slaves
	Namespace     map[string]string          // foo/bar.txt -> asdf332789asfj, purely controlled by client
	Queue         *OperationQueue
}

type Slave struct {
	IP     string
	Status Status
	Files  map[string]bool // hashes that a slave has
}

type OperationQueue struct {
	Queue []QueueItem
	QLock *sync.RWMutex
}

type QueueItem struct {
	Uid      string
	Filename string
	Hash     string
}

func (c *OperationQueue) Enqueue(qItem QueueItem) {
	c.QLock.Lock()
	defer c.QLock.Unlock()
	c.Queue = append(c.Queue, qItem)
}

func (c *OperationQueue) Dequeue() error {
	if len(c.Queue) > 0 {
		c.QLock.Lock()
		defer c.QLock.Unlock()
		c.Queue = c.Queue[1:]
		return nil
	}
	return fmt.Errorf("Queue is empty")
}

func (c *OperationQueue) Front() (QueueItem, error) {
	if len(c.Queue) > 0 {
		c.QLock.Lock()
		defer c.QLock.Unlock()
		return c.Queue[0], nil
	}
	var empty QueueItem
	return empty, fmt.Errorf("Queue is empty")
}

func (c *OperationQueue) Size() int {
	return len(c.Queue)
}

func (c *OperationQueue) Empty() bool {
	return len(c.Queue) == 0
}

func (c *OperationQueue) ReturnObj() []QueueItem {
	return c.Queue
}

// Status is an enumerated type
type Status string

const (
	OVERLOADED  Status = "Current load exceeds threshold"
	UNDERLOADED Status = "Current load does not exceed threshold"
)
