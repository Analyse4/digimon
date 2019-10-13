package panel

import "sync"

type rpcSelectionCounterPanel struct {
	mu      *sync.Mutex
	counter map[uint64]bool
}

func NewRpcSelectionCounterPanel() *rpcSelectionCounterPanel {
	return &rpcSelectionCounterPanel{
		mu:      new(sync.Mutex),
		counter: make(map[uint64]bool),
	}
}

func (rpcSCP *rpcSelectionCounterPanel) Update(id uint64) {
	rpcSCP.mu.Lock()
	defer rpcSCP.mu.Unlock()
	rpcSCP.counter[id] = true
}

func (rpcSCP *rpcSelectionCounterPanel) IsEnd() bool {
	rpcSCP.mu.Lock()
	defer rpcSCP.mu.Unlock()
	if len(rpcSCP.counter) <= 2 {
		return false
	}
	for _, r := range rpcSCP.counter {
		if r == false {
			return false
		}
	}
	return true
}

func (rpcSCP *rpcSelectionCounterPanel) Refresh() {
	rpcSCP.mu.Lock()
	defer rpcSCP.mu.Unlock()
	for i := range rpcSCP.counter {
		rpcSCP.counter[i] = false
	}
}
