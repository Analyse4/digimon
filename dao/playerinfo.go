package dao

import (
	"digimon/player"
)

func InsertPlayerInfo(pi *player.Player) error {
	_, err := db.Exec("insert player_info values (?, ?)", pi.Id, pi.NickName)
	return err
}
