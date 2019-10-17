package dao

import (
	"github.com/Analyse4/digimon/room"
	"strconv"
	"time"
)

//TODO: create time need modify
func InsertRoomInfo(ri *room.Room) error {
	_, err := db.Exec("insert into room_info (room_id, game_status, room_status, type, current_num, create_time, player1_id) values (?, ?, ?, ?, ?, ?, ?)", ri.Id, ri.IsStart, ri.IsOpen, ri.Type, ri.CurrentNum, time.Now(), ri.PlayerInfos[0].Id)
	if err != nil {
		return err
	}
	return nil

}

func UpdateRoomInfo(ri *room.Room) error {
	var err error
	if !ri.IsOpen {
		_, err = db.Exec("update room_info set room_status=? where room_id = ? and room_status = ?", ri.IsOpen, ri.Id, 1)
		if err != nil {
			return err
		}
		return nil
	}
	//playerItem := "player" + strconv.Itoa(int(ri.NewSeated + 1)) + "_id"
	//_, err := db.Exec("update room_info set game_status=?, current_num=?, " + playerItem + "=? where id = ?", ri.IsStart, ri.CurrentNum, ri.PlayerInfos[ri.NewSeated].Id, ri.Id)
	//if err != nil {
	//	return err
	//}
	//return nil

	for i, v := range ri.PlayerInfos {
		if v == nil {
			playerItem := "player" + strconv.Itoa(i+1) + "_id"
			_, err = db.Exec("update room_info set game_status=?, current_num=?, "+playerItem+"=? where room_id=? and room_status=?", ri.IsStart, ri.CurrentNum, nil, ri.Id, 1)
		} else {
			playerItem := "player" + strconv.Itoa(i+1) + "_id"
			_, err = db.Exec("update room_info set game_status=?, current_num=?, "+playerItem+"=? where room_id=? and room_status=?", ri.IsStart, ri.CurrentNum, ri.PlayerInfos[i].Id, ri.Id, 1)
		}
	}
	if err != nil {
		return err
	}
	return nil
}
