package service

import (
	"digimon/peer/sessionmanager"
)

type Service interface {
	Init(name, codecTyp, acceptorTyp, addr string)
	Start()
	GetAddr() string
	GetSessionManager() (*sessionmanager.SessionManager, error)
	GetName() string
	CleanerListen()
}
