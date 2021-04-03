package structs

import "sync"

// foo/bar.txt -> asdf332789asfj, purely controlled by client
type Namespace struct {
	rwLock    *sync.RWMutex
	namespace map[string]string
}

func (n *Namespace) SetHash(filename string, hash string) {
	n.rwLock.Lock()
	n.namespace[filename] = hash
	n.rwLock.Unlock()
}

func (n *Namespace) GetHash(filename string) string {
	n.rwLock.RLock()
	defer n.rwLock.RUnlock()
	return n.namespace[filename]
}

func (n *Namespace) LinkedHashes() map[string]bool {
	linkedHashes := make(map[string]bool)
	n.rwLock.RLock()
	for _, v := range n.namespace {
		linkedHashes[v] = true
	}
	defer n.rwLock.RUnlock()
	return linkedHashes
}
