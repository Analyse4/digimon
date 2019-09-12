package roommanager

import (
	"digimon/room"
	"fmt"
	"sync"
)

type RoomManager struct {
	Mu      *sync.Mutex
	RoomMap map[uint64]*room.Room
}

func New() *RoomManager {
	return &RoomManager{
		Mu:      new(sync.Mutex),
		RoomMap: make(map[uint64]*room.Room),
	}
}

func (rm *RoomManager) add(r *room.Room) error {
	if rm.RoomMap[r.Id] != nil {
		return fmt.Errorf("duplicate player")
	}
	rm.RoomMap[r.Id] = r
	return nil
}

func (rm *RoomManager) GetIdleRoom() *room.Room {
	rm.Mu.Lock()
	defer rm.Mu.Unlock()
	for _, r := range rm.RoomMap {
		if !r.IsStart {
			return r
		}
	}
	r := room.New()
	rm.add(r)
	return r
}
