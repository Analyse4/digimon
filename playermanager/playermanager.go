package playermanager

import (
	"digimon/logger"
	"digimon/player"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

var (
	log *logrus.Entry
)

func init() {
	log = logger.GetLogger().WithField("pkg", "playermanager")
}

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

	log.WithFields(logrus.Fields{
		"player_id":        p.Id,
		"total_player_num": len(pm.PlayerMap),
	}).Debug("add player successful")
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

func (pm *PlayerManager) Delete(id uint64) {
	pm.Mu.Lock()
	defer pm.Mu.Unlock()
	delete(pm.PlayerMap, id)
}
