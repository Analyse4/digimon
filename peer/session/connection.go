package session

import (
	"github.com/Analyse4/digimon/codec"
	"github.com/Analyse4/digimon/peer/cleaner"
	"sync"
)

type Connection interface {
	ReadLoop(codec.Codec, *Session)
	WriteLoop()
	GetID() int64
	SetID(int64)
	GetReqDeleteConn() chan<- *cleaner.CleanerMeta
	SetReqDeleteConn(chan<- *cleaner.CleanerMeta)
	GetWaitGroup() *sync.WaitGroup
	Close()
	Send([]byte)
	CloseSendBuffer()
}
