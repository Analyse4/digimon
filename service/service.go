package service

import (
	"digimon/acceptor"
	"digimon/handler"
	"github.com/golang/glog"
)

type Service interface {
	Start()
}

//TODO
func New(name, codecTyp, acceptorTyp, addr string) (Service, error) {
	glog.Info("create service successful!")
	glog.Info("name: " + name + "codec: " + codecTyp + "acceptor: " + acceptorTyp + "addr: " + addr)
	dgm := new(handler.Digimon)
	dgm.Name = name
	dgm.Acceptor, _ = acceptor.Get(acceptorTyp, addr, codecTyp)
	return dgm, nil
}
