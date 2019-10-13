package room

import (
	"digimon/player"
	"fmt"
	"strconv"
	"testing"
)

type condition struct {
	name      string
	condition *rulingCondition
}

func (con *condition) setCondition(AID uint64, AIdentityLevel int32, ASkillLevel int32, TID uint64, TIdentityLevel int32, TSkillType int32, TSkillLevel int32, TEscape bool) *condition {
	if TSkillType == player.DEFENCE {
		con.name = strconv.Itoa(int(AIdentityLevel)) + strconv.Itoa(int(ASkillLevel)) + "-" + strconv.Itoa(int(TIdentityLevel)) + "0" + strconv.Itoa(int(TSkillLevel))
	} else {
		con.name = strconv.Itoa(int(AIdentityLevel)) + strconv.Itoa(int(ASkillLevel)) + "-" + strconv.Itoa(int(TIdentityLevel)) + strconv.Itoa(int(TSkillLevel)) + "0"
	}
	if TEscape {
		con.name = con.name + "_e"
	}
	con.condition.AID = AID
	con.condition.AIdentityLevel = AIdentityLevel
	con.condition.ASkillLevel = ASkillLevel
	con.condition.TID = TID
	con.condition.TIdentityLevel = TIdentityLevel
	con.condition.TSkillType = TSkillType
	con.condition.TSkillLevel = TSkillLevel
	con.condition.TEscape = TEscape
	return con
}

func (con *condition) refresh() {
	con.name = ""
	con.condition.AID = 0
	con.condition.AIdentityLevel = 0
	con.condition.ASkillLevel = 0
	con.condition.TID = 0
	con.condition.TIdentityLevel = 0
	con.condition.TSkillType = 0
	con.condition.TSkillLevel = 0
	con.condition.TEscape = false
}

type Result struct {
	DeadID uint64
	RPCT   int32
	Err    error
}

func (r *Result) refresh() {
	r.DeadID = 0
	r.RPCT = -1
	r.Err = nil
}

type tAction struct {
	skillType  int32
	skillLevel int32
	escape     bool
}

//func TestSeparateRuling2(t *testing.T) {
//	tActionList := []tAction{
//		{
//			skillType:  player.ATTACK,
//			skillLevel: 1,
//			escape: false,
//		},
//		{
//			skillType: player.ATTACK,
//			skillLevel: 2,
//			escape: false,
//		},
//		{
//			skillType: player.DEFENCE,
//			skillLevel: 0,
//			escape: false,
//		},
//		{
//			skillType: player.DEFENCE,
//			skillLevel: 0,
//			escape: true,
//		},
//		{
//			skillType: player.DEFENCE,
//			skillLevel: 10,
//			escape: false,
//		},
//		{
//			skillType: player.DEFENCE,
//			skillLevel: 10,
//			escape: true,
//		},
//	}
//
//	ResultList := []Result{
//		{
//			DeadID: 0,
//			RPCT:   0,
//			Err:    nil,
//		},
//		{
//			DeadID: 1,
//			RPCT:   0,
//			Err:    nil,
//		},
//		{
//			DeadID: 0,
//			RPCT:   0,
//			Err:    nil,
//		},
//		{
//			DeadID: 0,
//			RPCT:   0,
//			Err:    nil,
//		},
//		{
//			DeadID: 2,
//			RPCT:   0,
//			Err:    nil,
//		},
//		{
//			DeadID: 0,
//			RPCT:   1,
//			Err:    nil,
//		},
//	}
//
//	cond := new(condition)
//	cond.condition= new(rulingCondition)
//	result := new(Result)
//
//	//cond.refresh()
//	//result.refresh()
//	//cond.setCondition(
//	//	1,
//	//	player.ROOKIE,
//	//	1,
//	//	2,
//	//	player.ULTIMATE,
//	//	player.DEFENCE,
//	//	3,
//	//	false,
//	//)
//	//result.DeadID, result.RPCT, result.Err = SeparateRuling(cond.condition)
//	//if result.DeadID != 0 || result.RPCT != 3 || result.Err != nil {
//	//	t.Errorf("name: %v, dead_id: %d, rpct: %d, err: %v", cond.name, result.DeadID, result.RPCT, result.Err)
//	//}
//	//fmt.Println(cond.name + " \u2713")
//}

func TestSeparateRuling(t *testing.T) {
	cond := new(condition)
	cond.condition = new(rulingCondition)
	cond.setCondition(
		1,
		player.ROOKIE,
		1,
		2,
		player.ROOKIE,
		player.ATTACK,
		1,
		false,
	)
	result := new(Result)
	result.DeadID, result.RPCT, result.Err = SeparateRuling(cond.condition)
	if result.DeadID != 0 || result.RPCT != 0 || result.Err != nil {
		t.Errorf("name: %v, dead_id: %d, rpct: %d, err: %v\n", cond.name, result.DeadID, result.RPCT, result.Err)
	}

	cond.refresh()
	result.refresh()
	cond.setCondition(
		1,
		player.ROOKIE,
		1,
		2,
		player.ROOKIE,
		player.DEFENCE,
		1,
		false,
	)
	result.DeadID, result.RPCT, result.Err = SeparateRuling(cond.condition)
	if result.DeadID != 0 || result.RPCT != 0 || result.Err != nil {
		t.Errorf("name: %v, dead_id: %d, rpct: %d, err: %v", cond.name, result.DeadID, result.RPCT, result.Err)
	}

	cond.refresh()
	result.refresh()
	cond.setCondition(
		1,
		player.ROOKIE,
		1,
		2,
		player.ROOKIE,
		player.DEFENCE,
		2,
		false,
	)
	result.DeadID, result.RPCT, result.Err = SeparateRuling(cond.condition)
	if result.DeadID != 2 || result.RPCT != 0 || result.Err != nil {
		t.Errorf("name: %v, dead_id: %d, rpct: %d, err: %v", cond.name, result.DeadID, result.RPCT, result.Err)
	}

	cond.refresh()
	result.refresh()
	cond.setCondition(
		1,
		player.ROOKIE,
		1,
		2,
		player.ROOKIE,
		player.DEFENCE,
		2,
		true,
	)
	result.DeadID, result.RPCT, result.Err = SeparateRuling(cond.condition)
	if result.DeadID != 0 || result.RPCT != 1 || result.Err != nil {
		t.Errorf("name: %v, dead_id: %d, rpct: %d, err: %v", cond.name, result.DeadID, result.RPCT, result.Err)
	}

	cond.refresh()
	result.refresh()
	cond.setCondition(
		1,
		player.ROOKIE,
		1,
		2,
		player.ROOKIE,
		player.DEFENCE,
		1,
		true,
	)
	result.DeadID, result.RPCT, result.Err = SeparateRuling(cond.condition)
	if result.DeadID != 0 || result.RPCT != 0 || result.Err != nil {
		t.Errorf("name: %v, dead_id: %d, rpct: %d, err: %v", cond.name, result.DeadID, result.RPCT, result.Err)
	}

	cond.refresh()
	result.refresh()
	cond.setCondition(
		1,
		player.ROOKIE,
		1,
		2,
		player.CHAMPION,
		player.ATTACK,
		1,
		false,
	)
	result.DeadID, result.RPCT, result.Err = SeparateRuling(cond.condition)
	if result.DeadID != 1 || result.RPCT != 0 || result.Err != nil {
		t.Errorf("name: %v, dead_id: %d, rpct: %d, err: %v", cond.name, result.DeadID, result.RPCT, result.Err)
	}

	cond.refresh()
	result.refresh()
	cond.setCondition(
		1,
		player.ROOKIE,
		1,
		2,
		player.CHAMPION,
		player.DEFENCE,
		1,
		false,
	)
	result.DeadID, result.RPCT, result.Err = SeparateRuling(cond.condition)
	if result.DeadID != 0 || result.RPCT != 0 || result.Err != nil {
		t.Errorf("name: %v, dead_id: %d, rpct: %d, err: %v", cond.name, result.DeadID, result.RPCT, result.Err)
	}

	cond.refresh()
	result.refresh()
	cond.setCondition(
		1,
		player.ROOKIE,
		1,
		2,
		player.CHAMPION,
		player.DEFENCE,
		2,
		false,
	)
	result.DeadID, result.RPCT, result.Err = SeparateRuling(cond.condition)
	if result.DeadID != 0 || result.RPCT != 2 || result.Err != nil {
		t.Errorf("name: %v, dead_id: %d, rpct: %d, err: %v", cond.name, result.DeadID, result.RPCT, result.Err)
	}

	cond.refresh()
	result.refresh()
	cond.setCondition(
		1,
		player.ROOKIE,
		1,
		2,
		player.CHAMPION,
		player.DEFENCE,
		2,
		true,
	)
	result.DeadID, result.RPCT, result.Err = SeparateRuling(cond.condition)
	if result.DeadID != 0 || result.RPCT != 3 || result.Err != nil {
		t.Errorf("name: %v, dead_id: %d, rpct: %d, err: %v", cond.name, result.DeadID, result.RPCT, result.Err)
	}

	cond.refresh()
	result.refresh()
	cond.setCondition(
		1,
		player.ROOKIE,
		1,
		2,
		player.CHAMPION,
		player.EVOLVE,
		2,
		true,
	)
	result.DeadID, result.RPCT, result.Err = SeparateRuling(cond.condition)
	if result.DeadID != 0 || result.RPCT != 3 || result.Err != nil {
		t.Errorf("name: %v, dead_id: %d, rpct: %d, err: %v", cond.name, result.DeadID, result.RPCT, result.Err)
	}

	cond.refresh()
	result.refresh()
	cond.setCondition(
		1,
		player.ROOKIE,
		1,
		2,
		player.ULTIMATE,
		player.ATTACK,
		1,
		false,
	)
	result.DeadID, result.RPCT, result.Err = SeparateRuling(cond.condition)
	if result.DeadID != 1 || result.RPCT != 0 || result.Err != nil {
		t.Errorf("name: %v, dead_id: %d, rpct: %d, err: %v", cond.name, result.DeadID, result.RPCT, result.Err)
	}
	fmt.Println(cond.name + " \u2713")

	cond.refresh()
	result.refresh()
	cond.setCondition(
		1,
		player.ROOKIE,
		1,
		2,
		player.ULTIMATE,
		player.DEFENCE,
		1,
		false,
	)
	result.DeadID, result.RPCT, result.Err = SeparateRuling(cond.condition)
	if result.DeadID != 0 || result.RPCT != 0 || result.Err != nil {
		t.Errorf("name: %v, dead_id: %d, rpct: %d, err: %v", cond.name, result.DeadID, result.RPCT, result.Err)
	}
	fmt.Println(cond.name + " \u2713")

	cond.refresh()
	result.refresh()
	cond.setCondition(
		1,
		player.ROOKIE,
		1,
		2,
		player.ULTIMATE,
		player.DEFENCE,
		3,
		false,
	)
	result.DeadID, result.RPCT, result.Err = SeparateRuling(cond.condition)
	if result.DeadID != 0 || result.RPCT != 3 || result.Err != nil {
		t.Errorf("name: %v, dead_id: %d, rpct: %d, err: %v", cond.name, result.DeadID, result.RPCT, result.Err)
	}
	fmt.Println(cond.name + " \u2713")

	cond.refresh()
	result.refresh()
	cond.setCondition(
		1,
		player.ROOKIE,
		1,
		2,
		player.ULTIMATE,
		player.DEFENCE,
		3,
		false,
	)
	result.DeadID, result.RPCT, result.Err = SeparateRuling(cond.condition)
	if result.DeadID != 0 || result.RPCT != 3 || result.Err != nil {
		t.Errorf("name: %v, dead_id: %d, rpct: %d, err: %v", cond.name, result.DeadID, result.RPCT, result.Err)
	}
	fmt.Println(cond.name + " \u2713")

	cond.refresh()
	result.refresh()
	cond.setCondition(
		1,
		player.ROOKIE,
		1,
		2,
		player.MEGA,
		player.DEFENCE,
		2,
		true,
	)
	result.DeadID, result.RPCT, result.Err = SeparateRuling(cond.condition)
	if result.DeadID != 0 || result.RPCT != 6 || result.Err != nil {
		t.Errorf("name: %v, dead_id: %d, rpct: %d, err: %v", cond.name, result.DeadID, result.RPCT, result.Err)
	}
	fmt.Println(cond.name + " \u2713")

	cond.refresh()
	result.refresh()
	cond.setCondition(
		1,
		player.ROOKIE,
		1,
		2,
		player.CHAMPION,
		player.ATTACK,
		1,
		false,
	)
	result.DeadID, result.RPCT, result.Err = SeparateRuling(cond.condition)
	if result.DeadID != 1 || result.RPCT != 0 || result.Err != nil {
		t.Errorf("name: %v, dead_id: %d, rpct: %d, err: %v", cond.name, result.DeadID, result.RPCT, result.Err)
	}
	fmt.Println(cond.name + " \u2713")
}
