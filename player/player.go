package player

import (
	"digimon/codec"
	"digimon/pbprotocol"
	"digimon/peer/session"
	"digimon/utils/randomid"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"sync"
)

const (
	NULL    = -1
	POWERUP = iota + 1
	DEFENCE
	ESCAPE
	_
	ATTACK
	EVOLVE
	ROOKIE
	CHAMPION
	ULTIMATE
	MEGA
)

type Hero struct {
	mu            *sync.Mutex
	Identity      pbprotocol.DigimonIdentity
	IdentityLevel int32
	SkillPoint    int32
	SkillType     int32
	SkillLevel    int32
	SkillName     string
}

type Player struct {
	Id          uint64
	NickName    string
	RoomID      uint64
	DigiMonstor *Hero
	Sess        *session.Session
}

func New(sess *session.Session) (*Player, error) {
	p := new(Player)
	p.NickName = randomdata.FullName(randomdata.Male)
	p.Id = randomid.GetUniqueId()
	p.Sess = sess
	p.DigiMonstor = new(Hero)
	p.DigiMonstor.mu = new(sync.Mutex)
	return p, nil
}

func (p *Player) Send(router string, data interface{}) {
	c, _ := codec.Get("protobuf")
	msg, err := c.Marshal(router, data)
	if err != nil {
		//TODO: log
	}
	p.Sess.Conn.Send(msg)
}

func (p *Player) PowerUp() {
	p.DigiMonstor.mu.Lock()
	defer p.DigiMonstor.mu.Unlock()
	p.DigiMonstor.SkillPoint++
}

func (p *Player) PowerDown(num int32) error {
	p.DigiMonstor.mu.Lock()
	defer p.DigiMonstor.mu.Unlock()
	if (p.DigiMonstor.SkillPoint - num) < 0 {
		return fmt.Errorf("skill point not enough")
	}
	p.DigiMonstor.SkillPoint = p.DigiMonstor.SkillPoint - 2
	return nil
}

func (p *Player) Evolve(typ int32) {
	p.DigiMonstor.IdentityLevel = typ
	switch typ {
	case 1:
		p.DigiMonstor.Identity++
	case 2:
		p.DigiMonstor.Identity++
	case 3:
		if p.DigiMonstor.IdentityLevel == 0 {
			p.DigiMonstor.Identity = p.DigiMonstor.Identity + 3
		} else if p.DigiMonstor.IdentityLevel == 1 {
			p.DigiMonstor.Identity = p.DigiMonstor.Identity + 2
		} else if p.DigiMonstor.IdentityLevel == 2 {
			p.DigiMonstor.Identity = p.DigiMonstor.Identity + 1
		}
	}
}

func (p *Player) GetAttackName(typ int32) string {
	// TODO:
	return ""
}
