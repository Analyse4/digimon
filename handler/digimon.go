package handler

import (
	"digimon/peer/acceptor"
	"digimon/peer/connmanager"
	"fmt"
	"github.com/golang/glog"
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
	return nil
}

func (dgm *Digimon) GetConnManager() (*connmanager.ConnManager, error) {
	if dgm.ConnManager == nil {
		return nil, fmt.Errorf("connection manager is not allowcated!")
	}
	return dgm.ConnManager, nil
}
