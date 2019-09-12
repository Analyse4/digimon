package player

import (
	"digimon/pbprotocol"
)

type Player struct {
	Id       uint64
	NickName string
}

func New() (*Player, error) {
	p := new(pbprotocol.PlayerInfo)
	p.Nickname = "Joker"
	p.Id = 4
	return p, nil
}
