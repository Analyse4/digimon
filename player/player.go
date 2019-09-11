package player

import "digimon/dao/entity"

type Player struct {
	PlayerId int64
	NickName string
}

func New() (*entity.PlayerInfo, error) {
	p := new(entity.PlayerInfo)
	p.NickName = "Joker"
	p.Id = 4
	//TODO: add player manager
	return p, nil
}
