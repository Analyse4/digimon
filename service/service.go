package service

import (
	"digimon/acceptor"
	"digimon/handler"
	"fmt"
)

type Service interface {
	Start()
}

//TODO
func New(name, codecTyp, acceptorTyp string, addr string) (Service, error) {
	fmt.Println("create service successful!")
	fmt.Println("name: " + name + "codec: " + codecTyp + "acceptor: " + acceptorTyp + "addr: " + addr)
	dgm := new(handler.Digimon)
	dgm.Name = name
	dgm.Acceptor, _ = acceptor.Get(acceptorTyp, addr, codecTyp)
	return dgm, nil
}
