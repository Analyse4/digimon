package room

import (
	"digimon/pbprotocol"
	"digimon/player"
	"sync"
)

const (
	TWO = 0
)

//TODO: rewrite room id generate
var currentRoomID uint64

type Room struct {
	Mu          *sync.Mutex
	Id          uint64
	IsStart     bool
	Type        pbprotocol.RoomInfo_RoomType
	CurrentNum  uint32
	PlayerInfos []*player.Player
}

func New() *Room {
	return &Room{
		Mu:          new(sync.Mutex),
		Id:          currentRoomID + 1,
		IsStart:     false,
		Type:        pbprotocol.RoomInfo_TWO,
		CurrentNum:  0,
		PlayerInfos: nil,
	}
}

func (r *Room) AddPlayer(p *player.Player) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	r.PlayerInfos = append(r.PlayerInfos, p)
	r.CurrentNum++
	if r.CurrentNum == 2 {
		r.IsStart = true
	}
}
