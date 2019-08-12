package connmanager

import (
	"digimon/acceptor/websocket/wsconnection"
)

type ConnManager struct {
	connMap   map[int64]*wsconnection.Connection
	currentID int64
}

func (cm *ConnManager) Add(c *wsconnection.Connection) {
	cm.currentID = cm.currentID + 1
	c.Id = cm.currentID
	cm.connMap[cm.currentID] = c
}

func New() *ConnManager {
	return &ConnManager{
		make(map[int64]*wsconnection.Connection),
		0,
	}
}
