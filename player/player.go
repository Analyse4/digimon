package player

import (
	"digimon/codec"
	"digimon/peer/session"
)

type Player struct {
	Id       uint64
	NickName string
	Sess     *session.Session
}

func New(sess *session.Session) (*Player, error) {
	p := new(Player)
	p.NickName = "Joker"
	p.Id = 4
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
