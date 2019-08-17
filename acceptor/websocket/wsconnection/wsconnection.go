package wsconnection

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

const SENDBUFFERSIZE = 100

type Connection struct {
	Id            int64
	Conn          *websocket.Conn
	wg            *sync.WaitGroup
	ReqDeleteConn chan<- int64
	SenndBuffer   chan []byte
}

func NewConnection(c *websocket.Conn) *Connection {
	nc := &Connection{Conn: c, wg: new(sync.WaitGroup), SenndBuffer: make(chan []byte, SENDBUFFERSIZE)}
	nc.wg.Add(2)
	go nc.readLoop(nc.wg)
	go nc.writeLoop(nc.wg)
	go func() {
		nc.wg.Wait()
		if nc.ReqDeleteConn != nil {
			nc.ReqDeleteConn <- nc.Id
		} else {
			glog.Error("connection haven't init successful!")
		}
		glog.Info("Connection %d\n is closed!", nc.Id)
	}()
	return nc

}

func (c *Connection) readLoop(wg *sync.WaitGroup) {
	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			close(c.SenndBuffer)
			wg.Done()
			return
		} else {
			fmt.Println(data)
			c.SenndBuffer <- data
		}
	}
}

func (c *Connection) writeLoop(wg *sync.WaitGroup) {
	for {
		select {
		case data, ok := <-c.SenndBuffer:
			if !ok {
				log.Printf("connection %d write buffer is closed!\n", c.Id)
				wg.Done()
				return
			} else {
				err := c.Conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					log.Println(err)
				}
				log.Printf("connection %d send: ", c.Id)
				log.Println(data)
			}
		}
	}
}
