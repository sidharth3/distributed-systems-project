package structs

import (
	"fmt"
	"sync"
)

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
