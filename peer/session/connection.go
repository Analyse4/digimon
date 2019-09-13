package session

import (
	"digimon/codec"
	"sync"
)

type Connection interface {
	ReadLoop(codec.Codec, *Session)
	WriteLoop()
	GetID() int64
	SetID(int64)
	GetReqDeleteConn() chan<- int64
	SetReqDeleteConn(chan<- int64)
	GetWaitGroup() *sync.WaitGroup
	Close()
	Send([]byte)
}
