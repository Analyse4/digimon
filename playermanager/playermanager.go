package playermanager

import (
	"digimon/player"
	"fmt"
	"sync"
)

type PlayerManager struct {
	Mu        *sync.Mutex
	PlayerMap map[uint64]*player.Player
}

func New() *PlayerManager {
	return &PlayerManager{
		Mu:        new(sync.Mutex),
		PlayerMap: make(map[uint64]*player.Player),
	}
}

func (pm *PlayerManager) Add(p *player.Player) error {
	pm.Mu.Lock()
	defer pm.Mu.Unlock()
	if pm.PlayerMap[p.Id] != nil {
		return fmt.Errorf("duplicate player")
	}
	pm.PlayerMap[p.Id] = p
	return nil
}

func (pm *PlayerManager) Get(id uint64) (*player.Player, error) {
	pm.Mu.Lock()
	defer pm.Mu.Unlock()
	p := pm.PlayerMap[id]
	if p == nil {
		return nil, fmt.Errorf("player info not found")
	}
	return p, nil
}
