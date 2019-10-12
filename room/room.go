package room

import (
	"digimon/errorhandler"
	"digimon/pbprotocol"
	"digimon/player"
	"digimon/playermanager"
	"fmt"
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

type judgeFinalJudgeCondition struct {
	AIdentityLevel int32
	ASkillLevel    int32
	TIdentityLevel int32
	TSkillType     int32
	TSkillLevel    int32
	TEscape        bool
}

type JudgeInfo struct {
	PlayerID uint64
	Number   int32
}

type RoundResult struct {
	FinalJudgeID map[uint64][]*JudgeInfo
	DeadID       []uint64
}

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

func (r *Room) DeadAnalyse() (*RoundResult, error) {
	tmpATSet := make(map[uint64][]uint64)
	for _, v := range r.PlayerInfos {
		if v.DigiMonstor.SkillType == player.ATTACK {
			tmpATSet[v.Id] = make([]uint64, 0)
			tmpATSet[v.Id] = append(tmpATSet[v.Id], v.DigiMonstor.SkillTargets...)
		}
	}
	if len(tmpATSet) == 0 {
		return nil, nil
	}
	tmpRoundResult := new(RoundResult)
	tmpRoundResult.DeadID = make([]uint64, 0)
	tmpRoundResult.FinalJudgeID = make(map[uint64][]*JudgeInfo, 0)
	for attackerID, targets := range tmpATSet {
		apl, _ := r.GetPlayer(attackerID)
		for _, t := range targets {
			tpl, _ := r.GetPlayer(t)
			cond := new(judgeFinalJudgeCondition)
			cond.AIdentityLevel = apl.DigiMonstor.IdentityLevel
			cond.ASkillLevel = apl.DigiMonstor.SkillLevel
			cond.TIdentityLevel = tpl.DigiMonstor.SkillLevel
			cond.TSkillType = tpl.DigiMonstor.SkillType
			cond.TSkillLevel = tpl.DigiMonstor.SkillLevel
			if jnum, err := getFinalJudgeNum(cond); err != nil {
				return nil, err
			} else if jnum == 0 {
				tpl.DigiMonstor.IsDead = true
				tmpRoundResult.DeadID = append(tmpRoundResult.DeadID, t)
			} else {
				tmpJI := new(JudgeInfo)
				tmpJI.PlayerID = t
				tmpJI.Number = jnum
				tmpRoundResult.FinalJudgeID[attackerID] = append(tmpRoundResult.FinalJudgeID[attackerID], tmpJI)
			}
		}
	}
	return tmpRoundResult, nil
}

func getFinalJudgeNum(cond *judgeFinalJudgeCondition) (int32, error) {
	switch cond.AIdentityLevel {
	case player.ROOKIE:
		switch cond.TIdentityLevel {
		case player.ROOKIE:
			if (cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1)) && cond.TEscape {
				return 1, nil
			} else {
				return 0, nil
			}
		case player.CHAMPION:
			if (cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1)) && cond.TEscape {
				return 3, nil
			} else if cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1) {
				return 2, nil
			} else {
				return 0, nil
			}
		case player.ULTIMATE:
			if (cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1)) && cond.TEscape {
				return 4, nil
			} else if cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1) {
				return 3, nil
			} else {
				return 0, nil
			}
		case player.MEGA:
			if (cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1)) && cond.TEscape {
				return 6, nil
			} else if cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1) {
				return 5, nil
			} else {
				return 0, nil
			}
		default:
			return -1, errorhandler.ERR_PARAMETERINVALID_MSG
		}
	case player.CHAMPION:
		switch cond.TIdentityLevel {
		case player.ROOKIE:
			if (cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1)) && cond.TEscape {
				return 1, nil
			} else {
				return 0, nil
			}
		case player.CHAMPION:
			if (cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1)) && cond.TEscape {
				return 1, nil
			} else {
				return 0, nil
			}
		case player.ULTIMATE:
			if (cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1)) && cond.TEscape {
				return 3, nil
			} else if cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1) {
				return 2, nil
			} else {
				return 0, nil
			}
		case player.MEGA:
			if (cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1)) && cond.TEscape {
				return 4, nil
			} else if cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1) {
				return 3, nil
			} else {
				return 0, nil
			}
		default:
			return -1, errorhandler.ERR_PARAMETERINVALID_MSG
		}
	case player.ULTIMATE:
		switch cond.TIdentityLevel {
		case player.ROOKIE:
			if (cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1)) && cond.TEscape {
				return 1, nil
			} else {
				return 0, nil
			}
		case player.CHAMPION:
			if (cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1)) && cond.TEscape {
				return 1, nil
			} else {
				return 0, nil
			}
		case player.ULTIMATE:
			if (cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1)) && cond.TEscape {
				return 1, nil
			} else {
				return 0, nil
			}
		case player.MEGA:
			if (cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1)) && cond.TEscape {
				return 3, nil
			} else if cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1) {
				return 2, nil
			} else {
				return 0, nil
			}
		default:
			return -1, errorhandler.ERR_PARAMETERINVALID_MSG
		}
	case player.MEGA:
		switch cond.TIdentityLevel {
		case player.ROOKIE:
			if (cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1)) && cond.TEscape {
				return 1, nil
			} else {
				return 0, nil
			}
		case player.CHAMPION:
			if (cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1)) && cond.TEscape {
				return 1, nil
			} else {
				return 0, nil
			}
		case player.ULTIMATE:
			if (cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1)) && cond.TEscape {
				return 1, nil
			} else {
				return 0, nil
			}
		case player.MEGA:
			if (cond.TSkillType != player.DEFENCE || (cond.TSkillType == player.DEFENCE && cond.TSkillLevel != 1)) && cond.TEscape {
				return 1, nil
			} else {
				return 0, nil
			}
		default:
			return -1, errorhandler.ERR_PARAMETERINVALID_MSG
		}
	default:
		return -1, errorhandler.ERR_PARAMETERINVALID_MSG
	}
}

func (r *Room) GetPlayer(id uint64) (*player.Player, error) {
	for _, p := range r.PlayerInfos {
		if p.Id == id {
			return p, nil
		}
	}
	return nil, fmt.Errorf("player not found")
}

func (r *Room) IsAllDead() bool {
	for _, v := range r.PlayerInfos {
		if v.DigiMonstor.IsDead == false {
			return false
		}
	}
	return true
}

func (r *Room) RefreshAllHeroStatus() {
	for _, pl := range r.PlayerInfos {
		pl.DigiMonstor.RefreshHeroRoundStatus()
	}
}
