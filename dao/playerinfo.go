package dao

import (
	"digimon/player"
)

func InsertPlayerInfo(p *player.Player) error {
	_, err := db.Exec("insert player_info values (?, ?)", p.Id, p.NickName)
	return err
}
