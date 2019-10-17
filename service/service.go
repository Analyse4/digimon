package service

import (
	"github.com/Analyse4/digimon/peer/sessionmanager"
)

type Service interface {
	Init(name, codecTyp, acceptorTyp, addr string)
	Start()
	GetAddr() string
	GetSessionManager() (*sessionmanager.SessionManager, error)
	GetName() string
	CleanerListen()
}
