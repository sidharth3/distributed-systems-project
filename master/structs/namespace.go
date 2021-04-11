package structs

import (
	"fmt"
	"strings"
	"sync"
)

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

func (n *Namespace) DelFile(filename string) {
	n.rwLock.Lock()
	delete(n.namespace, filename)
	n.rwLock.Unlock()
}

func (n *Namespace) GetFile(path string) []string {
	files := make([]string, 0)
	n.rwLock.Lock()
	if path == "" {
		for filename := range n.namespace {
			files = append(files, filename)
		}
	} else {
		for filename := range n.namespace {
			fmt.Println(strings.SplitAfter(filename, path))
			if strings.SplitAfter(filename, path)[0] == path && strings.SplitAfter(filename, path)[1][0] == '/' {
				files = append(files, strings.Split(filename, path)[1])
			}
		}
	}
	n.rwLock.Unlock()
	return files
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

func (n *Namespace) CollectGarbage(unlinked map[string]bool) {
	n.rwLock.Lock()
	for filename := range unlinked {
		delete(n.namespace, filename)
	}
	n.rwLock.Unlock()
}

// for master replica------------
func (n *Namespace) ReturnNamespace() map[string] string{
	return n.namespace
}