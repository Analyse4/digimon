package handler

import (
	"digimon/dao"
	"digimon/errorhandler"
	"digimon/logger"
	"digimon/pbprotocol"
	"digimon/peer/acceptor"
	"digimon/peer/cleaner"
	"digimon/peer/session"
	"digimon/peer/sessionmanager"
	"digimon/player"
	"digimon/playermanager"
	"digimon/prometheus"
	"digimon/room"
	"digimon/roommanager"
	"digimon/svcregister"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"
)

const INVALIDID = 0

var (
	TYPEOFERROR = reflect.TypeOf((*error)(nil)).Kind()
	log         *logrus.Entry
)

func init() {
	log = logger.GetLogger().WithField("pkg", "handler")
}

type Digimon struct {
	Name           string
	Addr           string
	Acceptor       acceptor.Acceptor
	SessionManager *sessionmanager.SessionManager
	PlayerManager  *playermanager.PlayerManager
	RoomManager    *roommanager.RoomManager
	Cleaner        chan *cleaner.CleanerMeta
	SignalCapture  chan os.Signal
}

func (dgm *Digimon) Start() {
	dgm.Acceptor.Accept(dgm)
}

func (dgm *Digimon) GetAddr() string {
	return dgm.Addr
}

func (dgm *Digimon) Init(name, codecTyp, acceptorTyp, addr string) {
	dgm.Name = name
	dgm.Addr = addr
	acp, err := acceptor.Get(acceptorTyp)
	dgm.Acceptor = acp
	if err != nil {
		log.WithFields(logrus.Fields{
			"acceptor_type": acceptorTyp,
		}).Fatalln(err)
	}
	dgm.SessionManager = sessionmanager.New(codecTyp)
	dgm.PlayerManager = playermanager.New()
	dgm.RoomManager = roommanager.New()
	dgm.Cleaner = make(chan *cleaner.CleanerMeta, 100)
	dgm.SessionManager.SetCleaner(dgm.Cleaner)
	dgm.SignalCapture = make(chan os.Signal, 1)
	signal.Notify(dgm.SignalCapture, syscall.SIGINT, syscall.SIGTERM)
	dgm.Register()
	go dgm.CleanerListen()
	go dgm.signalHandler()
	log.WithFields(logrus.Fields{
		"name":     "digimon",
		"addr":     addr,
		"acceptor": acceptorTyp,
	}).Debug("init svc successful")
}

func (dgm *Digimon) CleanerListen() {
	log.Debug("connection cleaner start")
	for {
		select {
		case cmt := <-dgm.Cleaner:
			if cmt.ConnID == INVALIDID {
				log.WithFields(logrus.Fields{
					"request_clean_connection_id": cmt.ConnID,
				}).Warn("invalid connection id")
			}
			var rm *room.Room
			if cmt.PlayerID != 0 {
				player, err := dgm.PlayerManager.Get(cmt.PlayerID)
				if err == nil {
					roomID := player.RoomID
					rm, err = dgm.RoomManager.Get(roomID)
					if err == nil {
						rm.DeletePlayer(cmt.PlayerID)
						dao.UpdateRoomInfo(rm)
						if rm.CurrentNum == 0 {
							rm.IsOpen = false
							dao.UpdateRoomInfo(rm)
							dgm.RoomManager.Delete(roomID)
							prometheus.GetRoomGauge().Dec()
						} else {
							prometheus.GetInGameRoomGauge().Dec()
						}
					}
					dgm.PlayerManager.Delete(cmt.PlayerID)
					prometheus.GetPlayerGauge().Dec()
				}
			}
			sess := dgm.SessionManager.Get(cmt.ConnID)
			sess.Conn.Close()
			dgm.SessionManager.Delete(cmt.ConnID)

			if rm != nil {
				log.WithFields(logrus.Fields{
					"room_id":                 rm.Id,
					"player_id":               cmt.PlayerID,
					"session_id":              cmt.ConnID,
					"connection_id":           cmt.ConnID,
					"total_player":            len(dgm.PlayerManager.PlayerMap),
					"total_room":              len(dgm.RoomManager.RoomMap),
					"current_room_player_num": rm.CurrentNum,
				}).Debug("connection resource is released")
			} else {
				log.WithFields(logrus.Fields{
					"player_id":     cmt.PlayerID,
					"session_id":    cmt.ConnID,
					"connection_id": cmt.ConnID,
					"total_player":  len(dgm.PlayerManager.PlayerMap),
					"total_room":    len(dgm.RoomManager.RoomMap),
				}).Debug("connection resource is released")
			}
			log.Println(time.Now())
		}
	}
}

func (dgm *Digimon) signalHandler() {
	_ = <-dgm.SignalCapture
	dgm.RoomManager.Mu.Lock()
	defer dgm.RoomManager.Mu.Unlock()
	for _, r := range dgm.RoomManager.RoomMap {
		r.IsOpen = false
		dao.UpdateRoomInfo(r)
	}
}

func (dgm *Digimon) GetSessionManager() (*sessionmanager.SessionManager, error) {
	if dgm.SessionManager == nil {
		return nil, fmt.Errorf("session manager haven't allocated")
	}
	return dgm.SessionManager, nil
}

func (dgm *Digimon) Register() {
	typ := reflect.TypeOf(dgm)
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		if ok := checkHandlerMethod(m); !ok {
			continue
		}
		index := strings.ToLower(typ.Elem().Name()) + "." + strings.ToLower(m.Name)
		handler := new(svcregister.Handler)
		handler.Receiver = reflect.ValueOf(dgm)
		handler.Func = m
		handler.Typ = m.Type.In(2)
		svcregister.Set(index, handler)

		log.WithFields(logrus.Fields{
			"service": dgm.Name,
			"router":  index,
			"func":    handler.Func.Name,
		}).Debug("service handler register successful")
	}
}

func (dgm *Digimon) Login(sess *session.Session, req *pbprotocol.LoginReq) (*pbprotocol.LoginAck, error) {
	baseack := new(pbprotocol.BaseAck)
	ack := new(pbprotocol.LoginAck)
	ack.Base = baseack
	ack.PlayerInfo = new(pbprotocol.PlayerInfo)

	if sess.Get("PLAYERID") == nil {
		if req.Type == pbprotocol.LoginReq_Visitor {
			log.WithFields(logrus.Fields{
				"is_new_player": "true",
				"login_type":    "visitor",
			}).Info("new player login")
			player, err := player.New(sess)
			dgm.PlayerManager.Add(player)
			prometheus.GetPlayerGauge().Inc()
			if err != nil {
				log.Println(err)
				ack.Base.Result = errorhandler.ERR_SERVICEBUSY
				ack.Base.Msg = errorhandler.GetErrMsg(errorhandler.ERR_SERVICEBUSY)
			}
			err = dao.InsertPlayerInfo(player)
			if err != nil {
				log.WithFields(logrus.Fields{
					"player_id": player.Id,
				}).Debug("insert player info failed")
			}
			ack.Base.Result = errorhandler.SUCESS
			ack.Base.Msg = errorhandler.GetErrMsg(errorhandler.SUCESS)
			ack.PlayerInfo.Nickname = player.NickName
			ack.PlayerInfo.Id = player.Id
			sess.Set("PLAYERID", player.Id)
			return ack, err
		}
	} else {
		log.Println("already login")
	}
	return ack, nil
}

func (dgm *Digimon) JoinRoom(sess *session.Session, req *pbprotocol.JoinRoomReq) (*pbprotocol.JoinRoomAck, error) {
	baseack := new(pbprotocol.BaseAck)
	roominfo := new(pbprotocol.RoomInfo)
	ack := new(pbprotocol.JoinRoomAck)
	ack.Base = baseack
	ack.RoomInfo = roominfo

	playerID := sess.Get("PLAYERID")
	if playerID == nil {
		logrus.Debug("user not login")
		ack.Base.Result = errorhandler.ERR_USERNOTLOGIN
		ack.Base.Msg = errorhandler.GetErrMsg(errorhandler.ERR_USERNOTLOGIN)
		return ack, nil
	}
	room, isNew := dgm.RoomManager.GetIdleRoom()
	pl, err := dgm.PlayerManager.Get(playerID.(uint64))
	if err != nil {
		log.WithFields(logrus.Fields{
			"player_id": playerID,
		}).Debug(err)
	}
	room.AddPlayer(pl)
	//TODO: dao update old room info
	if isNew {
		prometheus.GetRoomGauge().Inc()
		dao.InsertRoomInfo(room)
	} else {
		dao.UpdateRoomInfo(room)
	}
	if room.IsStart {
		ack := new(pbprotocol.StartGameAck)
		roominfo := new(pbprotocol.RoomInfo)
		ack.RoomInfo = roominfo
		ack.RoomInfo.RoomId = room.Id
		ack.RoomInfo.Type = room.Type
		ack.RoomInfo.CurrentPlayerNum = room.CurrentNum
		ack.RoomInfo.IsStart = room.IsStart
		for _, p := range room.PlayerInfos {
			p.DigiMonstor.Identity = pbprotocol.DigimonIdentity_PALMON
			p.DigiMonstor.IdentityLevel = player.ROOKIE
			p.DigiMonstor.SkillType = player.NULL
			p.DigiMonstor.SkillLevel = player.NULL
			p.DigiMonstor.SkillName = ""

			tmpPlayerInfo := new(pbprotocol.PlayerInfo)
			tmpPlayerInfo.Hero = new(pbprotocol.Hero)
			tmpPlayerInfo.Id = p.Id
			tmpPlayerInfo.Nickname = p.NickName
			tmpPlayerInfo.RoomId = p.RoomID
			tmpPlayerInfo.Seat = p.Seat
			tmpPlayerInfo.Hero.Identity = p.DigiMonstor.Identity
			tmpPlayerInfo.Hero.IdentityLevel = p.DigiMonstor.IdentityLevel
			tmpPlayerInfo.Hero.SkillType = p.DigiMonstor.SkillType
			tmpPlayerInfo.Hero.SkillLevel = p.DigiMonstor.SkillLevel
			tmpPlayerInfo.Hero.SkillName = p.DigiMonstor.SkillName

			ack.RoomInfo.PlayerInfos = append(ack.RoomInfo.PlayerInfos, tmpPlayerInfo)
		}

		go room.BroadCast("digimon.startgame", ack)
		prometheus.GetInGameRoomGauge().Inc()
	}

	ack.Base.Result = errorhandler.SUCESS
	ack.Base.Msg = errorhandler.GetErrMsg(errorhandler.SUCESS)
	ack.RoomInfo.RoomId = room.Id
	ack.RoomInfo.Type = room.Type
	ack.RoomInfo.CurrentPlayerNum = room.CurrentNum
	ack.RoomInfo.IsStart = room.IsStart
	for _, p := range room.PlayerInfos {
		tmpPlayerInfo := new(pbprotocol.PlayerInfo)
		tmpPlayerInfo.Id = p.Id
		tmpPlayerInfo.Nickname = p.NickName
		tmpPlayerInfo.RoomId = p.RoomID
		tmpPlayerInfo.Seat = p.Seat
		ack.RoomInfo.PlayerInfos = append(ack.RoomInfo.PlayerInfos, tmpPlayerInfo)
	}

	log.WithFields(logrus.Fields{
		"room_id":            room.Id,
		"current_player_num": room.CurrentNum,
	}).Debug("room info")

	return ack, nil
}

func (dgm *Digimon) ReleaseSkill(sess *session.Session, req *pbprotocol.ReleaseSkillReq) (*pbprotocol.ReleaseSkillAck, error) {
	ack := new(pbprotocol.ReleaseSkillAck)
	ack.Base = new(pbprotocol.BaseAck)
	ack.Hero = new(pbprotocol.Hero)
	playerID := sess.Get("PLAYERID")
	if playerID == nil {
		logrus.Debug("user not login")
		ack.Base.Result = errorhandler.ERR_USERNOTLOGIN
		ack.Base.Msg = errorhandler.GetErrMsg(errorhandler.ERR_USERNOTLOGIN)
		return ack, nil
	}
	pl, err := dgm.PlayerManager.Get(playerID.(uint64))
	if err != nil {
		log.WithFields(logrus.Fields{
			"player_id": playerID,
		}).Debug(err)
	}
	if req.SkillLevel < 0 || req.SkillType < 0 {
		logrus.Debug("parameter invalid")
		err = errorhandler.ERR_PARAMETERINVALID_MSG
	}
	switch req.SkillType {
	case player.POWERUP:
		pl.PowerUp()
	case player.DEFENCE:
		pl.DigiMonstor.SkillLevel = req.SkillLevel
	case player.ESCAPE:
		pl.DigiMonstor.IsEscape = true
	case player.ATTACK:
		err := pl.PowerDown(req.SkillLevel)
		if err != nil {
			logrus.Println(err)
			break
		}
		pl.DigiMonstor.SkillLevel = req.SkillLevel
		pl.DigiMonstor.SkillName = pl.GetAttackName(req.SkillLevel)
		pl.DigiMonstor.SetSkillTargets(req.SkillTargets)
	case player.EVOLVE:
		//mega-evolve
		//super-evolve
		//evolve
		if req.SkillLevel == 3 {
			err = pl.PowerDown(5)
		} else if req.SkillLevel == 2 {
			err = pl.PowerDown(3)
		} else if req.SkillLevel == 1 {
			err = pl.PowerDown(2)
		}
		if err != nil {
			logrus.Println(err)
			break
		}
		pl.DigiMonstor.SkillLevel = req.SkillLevel
		pl.Evolve(req.SkillLevel)
	default:
		logrus.Debug("parameter invalid")
		err = errorhandler.ERR_SKILLPOINTNOTENOUGH_MSG
	}
	if err == errorhandler.ERR_PARAMETERINVALID_MSG {
		ack.Base.Result = errorhandler.ERR_PARAMETERINVALID
		ack.Base.Msg = errorhandler.GetErrMsg(errorhandler.ERR_PARAMETERINVALID)
		return ack, nil
	} else if err == errorhandler.ERR_SKILLPOINTNOTENOUGH_MSG {
		ack.Base.Result = errorhandler.ERR_SKILLPOINTNOTENOUGH
		ack.Base.Msg = errorhandler.GetErrMsg(errorhandler.ERR_SKILLPOINTNOTENOUGH)
		return ack, nil
	}
	pl.DigiMonstor.SkillType = req.SkillType
	rm, err := dgm.RoomManager.Get(pl.RoomID)
	rm.UpdatePlayerInfo(dgm.PlayerManager, pl.Id)
	if err != nil {
		logrus.Debug("room not find")
		ack.Base.Result = errorhandler.ERR_SERVICEBUSY
		ack.Base.Msg = errorhandler.GetErrMsg(errorhandler.ERR_SERVICEBUSY)
		return ack, nil
	}
	rm.Skills.Update(pl.Seat)
	if rm.Skills.IsSkillsReady(rm.Type) {
		bcAck := new(pbprotocol.RoundResultAck)
		bcAck.RoundResult.DeadId = make([]uint64, 0)
		bcAck.RoundResult.Judges = make(map[uint64]*pbprotocol.RoundResultAck_RoundResult_ListJudge)
		if roundResult, err := rm.DeadAnalyse(); err != nil {
			logrus.Debug(err)
			ack.Base.Result = errorhandler.ERR_PARAMETERINVALID
			ack.Base.Msg = errorhandler.GetErrMsg(errorhandler.ERR_PARAMETERINVALID)
			return ack, nil
		} else if rm.IsAllDead() {
			bcAck.IsEnd = true
			bcAck.Round = rm.GetRound()
		} else {
			rm.UpdateRound()
			rm.Skills.Refresh()
			rm.RefreshAllHeroStatus()
			bcAck.IsEnd = false
			bcAck.Round = rm.GetRound()
			bcAck.RoundResult.DeadId = append(bcAck.RoundResult.DeadId, roundResult.DeadID...)
			for attackerID, targets := range roundResult.FinalJudgeID {
				tmpTargets := new(pbprotocol.RoundResultAck_RoundResult_ListJudge)
				tmpTargets.Judge = make([]*pbprotocol.RoundResultAck_RoundResult_Judge, 0)
				for _, target := range targets {
					tmpTarget := new(pbprotocol.RoundResultAck_RoundResult_Judge)
					tmpTarget.PlayerId = target.PlayerID
					tmpTarget.Number = target.Number
					tmpTargets.Judge = append(tmpTargets.Judge, tmpTarget)
				}
				bcAck.RoundResult.Judges[attackerID] = tmpTargets
			}
		}
		go rm.BroadCast("digimon.roundresult", bcAck)
	}
	ack.Base.Result = errorhandler.SUCESS
	ack.Base.Msg = errorhandler.GetErrMsg(errorhandler.SUCESS)
	ack.Hero.Identity = pl.DigiMonstor.Identity
	ack.Hero.IdentityLevel = pl.DigiMonstor.IdentityLevel
	ack.Hero.SkillPoint = pl.DigiMonstor.SkillPoint
	ack.Hero.SkillType = pl.DigiMonstor.SkillType
	ack.Hero.SkillLevel = pl.DigiMonstor.SkillLevel
	ack.Hero.SkillName = pl.DigiMonstor.SkillName
	ack.Hero.IsEscape = pl.DigiMonstor.IsEscape
	ack.Hero.IsDead = pl.DigiMonstor.IsDead
	return ack, nil
}

//TODO: verification is not accurate enough
func checkHandlerMethod(m reflect.Method) bool {
	if m.Type.NumIn() != 3 || m.Type.NumOut() != 2 {
		return false
	}
	if m.Type.In(1).Kind() != reflect.Ptr || m.Type.Out(0).Kind() != reflect.Ptr {
		return false
	}
	return true
}

func (dgm *Digimon) GetName() string {
	return dgm.Name
}
