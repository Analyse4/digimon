package wsconnection

import (
	"fmt"
	"github.com/gorilla/websocket"
	"sync"
	"traefik/log"
)

type Connection struct {
	Id            int64
	Conn          *websocket.Conn
	wg            *sync.WaitGroup
	ReqDeleteConn chan<- int64
}

func NewConnection(c *websocket.Conn) *Connection {
	nc := &Connection{Conn: c, wg: new(sync.WaitGroup)}
	nc.wg.Add(1)
	go nc.readLoop(nc.wg)
	go func() {
		nc.wg.Wait()
		if nc.ReqDeleteConn != nil {
			nc.ReqDeleteConn <- nc.Id
		} else {
			log.Error("connection haven't init successful!")
		}
		fmt.Printf("Connection %d\n is closed!", nc.Id)
	}()
	return nc

}

func (c *Connection) readLoop(wg *sync.WaitGroup) {
	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			wg.Done()
			return
		} else {
			fmt.Println(data)
		}
	}
}
