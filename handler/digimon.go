package handler

import (
	"digimon/pbprotocol"
	"digimon/peer/acceptor"
	"digimon/peer/connmanager"
	"digimon/svcregister"
	"fmt"
	"github.com/golang/glog"
	"reflect"
	"strings"
)

var (
	TYPEOFERROR = reflect.TypeOf((*error)(nil)).Kind()
)

type Digimon struct {
	Name        string
	Addr        string
	Acceptor    acceptor.Acceptor
	ConnManager *connmanager.ConnManager
}

func (dgm *Digimon) Start() {
	dgm.Acceptor.Accept(dgm)
}

func (dgm *Digimon) GetAddr() string {
	return dgm.Addr
}

func (dgm *Digimon) Init(name, codecTyp, acceptorTyp, addr string) error {
	glog.Info("create service successful!")
	glog.Info("name: " + name + "codec: " + codecTyp + "acceptor: " + acceptorTyp + "addr: " + addr)
	dgm.Name = name
	dgm.Addr = addr
	dgm.Acceptor, _ = acceptor.Get(acceptorTyp)
	dgm.ConnManager = connmanager.New(codecTyp)
	dgm.Register()
	fmt.Println(svcregister.SVCRegister)
	return nil
}

func (dgm *Digimon) GetConnManager() (*connmanager.ConnManager, error) {
	if dgm.ConnManager == nil {
		return nil, fmt.Errorf("connection manager is not allowcated!")
	}
	return dgm.ConnManager, nil
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
		handler.Typ = m.Type.In(1)
		svcregister.Set(index, handler)
	}
}

func (dgm *Digimon) Login(req *pbprotocol.LoginReq) (*pbprotocol.LoginAck, error) {
	ack := new(pbprotocol.LoginAck)
	ack.Result = 0
	ack.Message = "everything is ok!"
	return ack, nil
}

//TODO: verification is not accurate enough
func checkHandlerMethod(m reflect.Method) bool {
	if m.Type.NumIn() != 2 || m.Type.NumOut() != 2 {
		return false
	}
	if m.Type.In(1).Kind() != reflect.Ptr || m.Type.Out(0).Kind() != reflect.Ptr {
		return false
	}
	return true
}
