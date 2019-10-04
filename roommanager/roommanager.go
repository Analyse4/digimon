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

func (rm *RoomManager) GetIdleRoom() (*room.Room, bool) {
	rm.Mu.Lock()
	defer rm.Mu.Unlock()
	for _, r := range rm.RoomMap {
		if !r.IsStart {
			return r, false
		}
	}
	r := room.New()
	rm.add(r)
	return r, true
}

func (rm *RoomManager) Get(id uint64) (*room.Room, error) {
	rm.Mu.Lock()
	defer rm.Mu.Unlock()
	room := rm.RoomMap[id]
	if room == nil {
		return nil, fmt.Errorf("room not found")
	}
	return room, nil
}

func (rm *RoomManager) Delete(id uint64) {
	rm.Mu.Lock()
	defer rm.Mu.Unlock()
	delete(rm.RoomMap, id)
}
