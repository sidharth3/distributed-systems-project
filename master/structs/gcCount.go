package structs

import "sync"

type GCCount struct {
	rwLock  *sync.RWMutex
	gcCount map[string]int
}

func (g *GCCount) Cycle(unlinked map[string]bool) map[string]bool {
	toDelete := make(map[string]bool)
	newGCCount := make(map[string]int)
	g.rwLock.Lock()
	for filename := range unlinked {
		if g.gcCount[filename] == 0 {
			newGCCount[filename] = 1
		} else {
			g.gcCount[filename] -= 1
			if g.gcCount[filename] == 0 {
				toDelete[filename] = true
			} else {
				newGCCount[filename] = g.gcCount[filename]
			}
		}
	}
	g.gcCount = newGCCount
	defer g.rwLock.Unlock()
	return toDelete
}

func (g *GCCount) NewFile(filename string) {
	g.rwLock.Lock()
	g.gcCount[filename] = 2
	g.rwLock.Unlock()
}
