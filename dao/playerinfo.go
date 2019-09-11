package dao

import (
	"digimon/dao/entity"
)

func InsertPlayerInfo(pi *entity.PlayerInfo) error {
	_, err := db.Exec("insert player_info values (?, ?)", pi.Id, pi.NickName)
	return err
}
