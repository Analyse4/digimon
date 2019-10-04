package dao

import (
	"digimon/pbprotocol"
	"digimon/room"
	"fmt"
	"time"
)

//TODO: create time need modify
func InsertRoomInfo(ri *room.Room) error {
	switch ri.Type {
	case pbprotocol.RoomInfo_TWO:
		if ri.CurrentNum == 1 {
			_, err := db.Exec("insert into room_info (id, status, type, current_num, create_time, player1_id) values ($1, $2, $3, $4, $5, $6)", ri.Id, ri.IsStart, ri.Type, ri.CurrentNum, time.Now(), ri.PlayerInfos[0].Id)
			if err != nil {
				return err
			}
		} else {
			_, err := db.Exec("insert into room_info (id, status, type, current_num, create_time, player2_id) values ($1, $2, $3, $4, $5, $6)", ri.Id, ri.IsStart, ri.Type, ri.CurrentNum, time.Now(), ri.PlayerInfos[1].Id)
			if err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("invailed room type")

	}
}
