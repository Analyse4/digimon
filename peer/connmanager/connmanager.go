package connmanager

import (
	"digimon/codec"
	"digimon/peer/connection"
	"fmt"
	"github.com/golang/glog"
	"log"
	"strconv"
	"sync"
)

const INVALIDID = 0

type ConnManager struct {
	connMap   map[int64]connection.Connection
	currentID int64
	mu        *sync.Mutex
	cleanConn chan int64
	codec     codec.Codec
}

func (cm *ConnManager) Add(c connection.Connection) {
	cm.mu.Lock()
	cm.currentID = cm.currentID + 1
	cm.connMap[cm.currentID] = c
	cm.mu.Unlock()
	c.SetID(cm.currentID)
	c.SetReqDeleteConn(cm.cleanConn)
	fmt.Println("Total connection num: " + strconv.Itoa(len(cm.connMap)))

	c.GetWaitGroup().Add(2)
	go c.ReadLoop(cm.codec)
	go c.WriteLoop()
	go func() {
		c.GetWaitGroup().Wait()
		if c.GetReqDeleteConn() != nil {
			c.GetReqDeleteConn() <- c.GetID()
		} else {
			glog.Error("connection haven't init successful!")
		}
		glog.Info("Connection %d\n is closed!", c.GetID())
	}()
}

func New(codecTyp string) *ConnManager {
	cd, _ := codec.Get(codecTyp)
	cm := &ConnManager{
		make(map[int64]connection.Connection),
		INVALIDID,
		new(sync.Mutex),
		make(chan int64),
		cd,
	}
	go func() {
		log.Println("conn cleaner start!")
		for {
			select {
			case connID := <-cm.cleanConn:
				if connID == INVALIDID {
					log.Println("Invalid connID")
				}
				delete(cm.connMap, connID)
				log.Printf("conn %d is deleted", connID)
			}
		}
	}()
	return cm
}

func (cm *ConnManager) GetCurrentConnID() int64 { return cm.currentID }
