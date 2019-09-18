package randomid

import (
	"math/rand"
	"sync"
)

var mu *sync.Mutex
var generatedId map[uint64]bool

func init() {
	mu = new(sync.Mutex)
	generatedId = make(map[uint64]bool)
}

func GetUniqueId() uint64 {
	mu.Lock()
	defer mu.Unlock()
	for {
		id := uint64(rand.Intn(1000000))
		if !generatedId[(id)] {
			generatedId[id] = true
			return id
		}
	}
}
