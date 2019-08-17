package connmanager

import (
	"digimon/acceptor/websocket/wsconnection"
	"fmt"
	"strconv"
	"sync"
	"traefik/log"
)

const INVALIDID = 0

type ConnManager struct {
	connMap   map[int64]*wsconnection.Connection
	currentID int64
	mu        *sync.Mutex
	cleanConn chan int64
}

func (cm *ConnManager) Add(c *wsconnection.Connection) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.currentID = cm.currentID + 1
	c.Id = cm.currentID
	c.ReqDeleteConn = cm.cleanConn
	cm.connMap[cm.currentID] = c
	fmt.Println("Total connection num: " + strconv.Itoa(len(cm.connMap)))
}

func New() *ConnManager {
	cm := &ConnManager{
		make(map[int64]*wsconnection.Connection),
		INVALIDID,
		new(sync.Mutex),
		make(chan int64),
	}
	go func() {
		log.Println("conn cleaner start!")
		for {
			select {
			case connID := <-cm.cleanConn:
				if connID == INVALIDID {
					log.Error("Invalid connID")
				}
				delete(cm.connMap, connID)
				log.Printf("conn %d is deleted", connID)
			}
		}
	}()
	return cm
}

func (cm *ConnManager) GetCurrentConnID() int64 { return cm.currentID }
