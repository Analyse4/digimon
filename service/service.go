package service

import (
	"digimon/peer/sessionmanager"
)

//TODO: GetconnManger should return connection interface
type Service interface {
	Init(name, codecTyp, acceptorTyp, addr string) error
	Start()
	GetAddr() string
	GetSessionManager() (*sessionmanager.SessionManager, error)
}
