package panel

import (
	"digimon/errorhandler"
	"digimon/player"
	"sync"
)

type rpcSelectionCounterPanel struct {
	mu      *sync.Mutex
	counter map[uint64]int32
}

func NewRpcSelectionCounterPanel() *rpcSelectionCounterPanel {
	return &rpcSelectionCounterPanel{
		mu:      new(sync.Mutex),
		counter: make(map[uint64]int32),
	}
}

func (rpcSCP *rpcSelectionCounterPanel) Update(id uint64, rpc int32) {
	rpcSCP.mu.Lock()
	defer rpcSCP.mu.Unlock()
	rpcSCP.counter[id] = rpc
}

func (rpcSCP *rpcSelectionCounterPanel) IsEnd() bool {
	rpcSCP.mu.Lock()
	defer rpcSCP.mu.Unlock()
	if len(rpcSCP.counter) <= 2 {
		return false
	}
	for _, r := range rpcSCP.counter {
		if r == 0 {
			return false
		}
	}
	return true
}

func (rpcSCP *rpcSelectionCounterPanel) Refresh() {
	rpcSCP.mu.Lock()
	defer rpcSCP.mu.Unlock()
	for i := range rpcSCP.counter {
		rpcSCP.counter[i] = 0
	}
}

func (rpcSCP *rpcSelectionCounterPanel) GetRPC(id uint64) int32 {
	rpcSCP.mu.Lock()
	defer rpcSCP.mu.Unlock()
	return rpcSCP.counter[id]
}

func (rpcSCP *rpcSelectionCounterPanel) RPCCompute() (uint64, error) {
	rpcIDInfos := rpcSetToSlice(rpcSCP.counter)
	if rpcIDInfos[0].rpc == rpcIDInfos[1].rpc {
		return 0, nil
	}
	switch rpcIDInfos[0].rpc {
	case player.ROCK:
		switch rpcIDInfos[1].rpc {
		case player.PAPER:
			return rpcIDInfos[1].id, nil
		case player.SCISSORS:
			return rpcIDInfos[0].id, nil
		default:
			return 0, errorhandler.ERR_PARAMETERINVALID_MSG
		}
	case player.PAPER:
		switch rpcIDInfos[1].rpc {
		case player.ROCK:
			return rpcIDInfos[0].id, nil
		case player.SCISSORS:
			return rpcIDInfos[1].id, nil
		default:
			return 0, errorhandler.ERR_PARAMETERINVALID_MSG
		}
	case player.SCISSORS:
		switch rpcIDInfos[1].rpc {
		case player.ROCK:
			return rpcIDInfos[1].id, nil
		case player.PAPER:
			return rpcIDInfos[0].id, nil
		default:
			return 0, errorhandler.ERR_PARAMETERINVALID_MSG
		}
	default:
		return 0, errorhandler.ERR_PARAMETERINVALID_MSG
	}
}

type rpcIDInfo struct {
	id  uint64
	rpc int32
}

func rpcSetToSlice(set map[uint64]int32) []*rpcIDInfo {
	rpcIDInfos := make([]*rpcIDInfo, 0)
	for id, rpc := range set {
		tmpRPCIDInfo := new(rpcIDInfo)
		tmpRPCIDInfo.id = id
		tmpRPCIDInfo.rpc = rpc
		rpcIDInfos = append(rpcIDInfos, tmpRPCIDInfo)
	}
	return rpcIDInfos
}
