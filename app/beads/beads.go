package beads

import (
	"sync"
)

type Beads struct {
	mu    sync.Mutex
	Beads map[string]map[string]int
}

func (b *Beads) Set(roomID string, direction string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.Beads[roomID][direction] = 10
}

func (b *Beads) Exist(roomID string, direction string) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.Beads[roomID][direction] == 10 {
		return true
	} else {
		return false
	}
}

func (b *Beads) DeleteWithKey(roomID string, key string) {
	delete(b.Beads[roomID], key)
}
