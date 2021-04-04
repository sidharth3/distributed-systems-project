package structs

import (
	"sync"
)

type OperationQueue struct {
	rwLock     *sync.RWMutex
	uids       []string              // UIDs in order of request [uid1, uid2, uid3]
	queueItems map[string]*QueueItem // Map from UID to QueueItem {uid1 : QueueItem1, uid2 : QueueItem2}
}

type QueueItem struct {
	Filename string
	Hash     string
	timedOut bool
}

func (c *OperationQueue) Enqueue(uid string, filename string) {
	c.rwLock.Lock()
	defer c.rwLock.Unlock()
	c.uids = append(c.uids, uid)
	c.queueItems[uid] = &QueueItem{filename, "", false}
}

func (c *OperationQueue) Dequeue() []*QueueItem {
	toCommit := make([]*QueueItem, 0)
	c.rwLock.Lock()
	newFront := len(c.uids)
	for i, uid := range c.uids {
		if c.queueItems[uid].Hash != "" { // If hash exists means its confirmed by slave
			toCommit = append(toCommit, c.queueItems[uid])
			delete(c.queueItems, uid)
		} else if c.queueItems[uid].timedOut {
			delete(c.queueItems, uid)
		} else {
			newFront = i
			break
		}
	}
	c.uids = c.uids[newFront:]
	defer c.rwLock.Unlock()
	return toCommit
}

func (c *OperationQueue) Timeout(uid string) {
	c.rwLock.Lock()
	if c.queueItems[uid] != nil {
		c.queueItems[uid].timedOut = true
	}
	c.rwLock.Unlock()
}

func (c *OperationQueue) Confirm(uid string, hash string) {
	c.rwLock.Lock()
	c.queueItems[uid].Hash = hash
	c.rwLock.Unlock()
}

func (c *OperationQueue) FirstUID() string {
	c.rwLock.RLock()
	defer c.rwLock.RUnlock()
	if len(c.uids) == 0 {
		return ""
	}
	return c.uids[0]
}
