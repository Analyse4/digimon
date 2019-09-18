package player

import (
	"digimon/codec"
	"digimon/peer/session"
	"digimon/utils/randomid"
	"github.com/Pallinder/go-randomdata"
)

type Player struct {
	Id       uint64
	NickName string
	Sess     *session.Session
}

func New(sess *session.Session) (*Player, error) {
	p := new(Player)
	p.NickName = randomdata.FullName(randomdata.Male)
	p.Id = randomid.GetUniqueId()
	p.Sess = sess
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
