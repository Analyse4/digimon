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
var (
	mu            *sync.Mutex = new(sync.Mutex)
	currentRoomID uint64
)

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
		r.PlayerInfos[0].Sess = p.Sess
		r.NewSeated = 0
	} else {
		r.PlayerInfos = append(r.PlayerInfos, p)
		r.NewSeated = int8(len(r.PlayerInfos) - 1)
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
