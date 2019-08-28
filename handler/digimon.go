package handler

import (
	"digimon/errorhandler"
	"digimon/logger"
	"digimon/pbprotocol"
	"digimon/peer/acceptor"
	"digimon/peer/session"
	"digimon/peer/sessionmanager"
	"digimon/player"
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

type Digimon struct {
	Name           string
	Addr           string
	Acceptor       acceptor.Acceptor
	SessionManager *sessionmanager.SessionManager
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
	if sess.Get("PLAYERID") == nil {
		log.Println("new player login")
		if req.Type == pbprotocol.LoginReq_Visitor {
			log.Println("login type: visitor")
			userinfo, err := player.New()
			if err != nil {
				log.Println(err)
				ack.Base.Result = errorhandler.ERR_SERVICEBUSY
				ack.Base.Msg = errorhandler.GetErrMsg(errorhandler.ERR_SERVICEBUSY)
				//TODO: close connection
			}
			ack.Base.Result = errorhandler.SUCESS
			ack.Base.Msg = errorhandler.GetErrMsg(errorhandler.SUCESS)
			ack.Nickname = userinfo.NickName
			sess.Set("PLAYERID", userinfo.PlayerId)
			return ack, err
		}
	} else {
		log.Println("already login")
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
