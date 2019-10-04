package dao

import (
	"digimon/player"
)

func InsertPlayerInfo(p *player.Player) error {
	_, err := db.Exec("insert visitor_player_info (user_id, nickname) values (?, ?)", p.Id, p.NickName)
	return err
}
