package handler

import (
	"digimon/acceptor"
)

type Digimon struct {
	Name     string
	Acceptor acceptor.Acceptor
}

func (dgm *Digimon) Start() {
	dgm.Acceptor.Accept()
}
