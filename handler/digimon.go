package handler

import (
	"digimon/dao"
	"digimon/errorhandler"
	"digimon/logger"
	"digimon/pbprotocol"
	"digimon/peer/acceptor"
	"digimon/peer/session"
	"digimon/peer/sessionmanager"
	"digimon/player"
	"digimon/playermanager"
	"digimon/roommanager"
	"digimon/svcregister"
	"fmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"strings"
)

var (
	TYPEOFERROR = reflect.TypeOf((*error)(nil)).Kind()
	log         *logrus.Entry
)

func init() {
	log = logger.GetLogger().WithField("pkg", "handler")
}

type CleanerMeta struct {
	ConnID   int64
	PlayerID uint64
}

type Digimon struct {
	Name           string
	Addr           string
	Acceptor       acceptor.Acceptor
	SessionManager *sessionmanager.SessionManager
	PlayerManager  *playermanager.PlayerManager
	RoomManager    *roommanager.RoomManager
	Cleaner        chan *CleanerMeta
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
	dgm.Register()
	log.WithFields(logrus.Fields{
		"name":     "digimon",
		"addr":     addr,
		"acceptor": acceptorTyp,
	}).Debug("init svc successful")
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
	room := dgm.RoomManager.GetIdleRoom()
	player, err := dgm.PlayerManager.Get(playerID.(uint64))
	if err != nil {
		log.WithFields(logrus.Fields{
			"player_id": playerID,
		}).Debug(err)
	}
	room.AddPlayer(player)
	//TODO: dao insert room info
	dao.InsertRoomInfo(room)
	if room.IsStart {
		ack := new(pbprotocol.StartGameAck)
		roominfo := new(pbprotocol.RoomInfo)
		ack.RoomInfo = roominfo
		ack.Identity = pbprotocol.DigimonIdentity_PALMON
		ack.RoomInfo.RoomId = room.Id
		ack.RoomInfo.Type = room.Type
		ack.RoomInfo.CurrentPlayerNum = room.CurrentNum
		ack.RoomInfo.IsStart = room.IsStart
		for i, p := range room.PlayerInfos {
			ack.RoomInfo.PlayerInfos[i].Id = p.Id
			ack.RoomInfo.PlayerInfos[i].Nickname = p.NickName
		}

		go room.BroadCast("digimon.startgame", ack)
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
		ack.RoomInfo.PlayerInfos = append(ack.RoomInfo.PlayerInfos, tmpPlayerInfo)
	}
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
