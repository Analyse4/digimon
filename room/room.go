package room

import (
	"digimon/pbprotocol"
	"digimon/player"
	"digimon/playermanager"
	"github.com/sirupsen/logrus"
	"sync"
)

const (
	TWO = 0
)

//TODO: rewrite room id generate
var (
	mu            *sync.Mutex = new(sync.Mutex)
	currentRoomID uint64
)

type skillSet struct {
	mu     *sync.Mutex
	skills map[int32]bool
}

// IsStart: game is start or not
// NewSeated: new joined player's seat num
// IsOpen: room is open or not
type Room struct {
	Mu          *sync.Mutex
	Id          uint64
	IsStart     bool
	Type        pbprotocol.RoomInfo_RoomType
	CurrentNum  uint32
	PlayerInfos []*player.Player
	NewSeated   int8
	IsOpen      bool
	Skills      *skillSet
	round       int32
}

// temporary only have two-player room
func New() *Room {
	r := &Room{
		Mu:          new(sync.Mutex),
		Id:          currentRoomID + 1,
		IsStart:     false,
		Type:        pbprotocol.RoomInfo_TWO,
		CurrentNum:  0,
		PlayerInfos: make([]*player.Player, 1),
		NewSeated:   -1,
		IsOpen:      true,
		Skills:      &skillSet{mu: new(sync.Mutex), skills: make(map[int32]bool)},
		round:       1,
	}
	mu.Lock()
	currentRoomID = r.Id
	mu.Unlock()
	return r
}

func (r *Room) AddPlayer(p *player.Player) {
	p.RoomID = r.Id
	r.Mu.Lock()
	defer r.Mu.Unlock()
	if r.PlayerInfos[0] == nil {
		r.PlayerInfos[0] = new(player.Player)
		r.PlayerInfos[0].Id = p.Id
		r.PlayerInfos[0].NickName = p.NickName
		r.PlayerInfos[0].RoomID = p.RoomID
		r.PlayerInfos[0].DigiMonstor = p.DigiMonstor
		r.PlayerInfos[0].Sess = p.Sess
		r.NewSeated = 0
		p.Seat = 0
	} else {
		r.PlayerInfos = append(r.PlayerInfos, p)
		r.NewSeated = int8(len(r.PlayerInfos) - 1)
		p.Seat = int32(len(r.PlayerInfos) - 1)
	}
	r.CurrentNum++
	if r.CurrentNum == 2 {
		r.IsStart = true
	}
}

func (r *Room) DeletePlayer(id uint64) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	for i, v := range r.PlayerInfos {
		if v != nil && v.Id == id {
			r.PlayerInfos[i] = nil
		}
	}
	r.CurrentNum--
	r.IsStart = false
}

func (r *Room) BroadCast(router string, data interface{}) {
	for _, p := range r.PlayerInfos {
		p.Send(router, data)
	}
}

func (sks *skillSet) Update(seat int32) {
	sks.mu.Lock()
	defer sks.mu.Unlock()
	sks.skills[seat] = true
}

func (sks *skillSet) Refresh() {
	sks.mu.Lock()
	defer sks.mu.Unlock()
	for i := range sks.skills {
		sks.skills[i] = false
	}
}

func (sks *skillSet) IsSkillsReady(roomNum pbprotocol.RoomInfo_RoomType) bool {
	switch roomNum {
	case pbprotocol.RoomInfo_TWO:
		if len(sks.skills) < 2 {
			return false
		}
	default:
		logrus.Error("room type invalid")
		return false
	}
	for _, v := range sks.skills {
		if v == false {
			return false
		}
	}
	return true
}

func (r *Room) UpdateRound() { r.round++ }

func (r *Room) GetRound() int32 { return r.round }

func (r *Room) UpdatePlayerInfo(pm *playermanager.PlayerManager, ids ...uint64) {
	for id := range ids {
		for i, v := range r.PlayerInfos {
			if v.Id == uint64(id) {
				pl, _ := pm.Get(v.Id)
				r.PlayerInfos[i].DigiMonstor = pl.DigiMonstor
			}
		}
	}
}
