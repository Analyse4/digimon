package service

import (
	"digimon/peer/connmanager"
)

//TODO: GetconnManger should return connection interface
type Service interface {
	Init(name, codecTyp, acceptorTyp, addr string) error
	Start()
	GetAddr() string
	GetConnManager() (*connmanager.ConnManager, error)
}
