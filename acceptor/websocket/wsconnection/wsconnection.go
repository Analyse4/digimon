package wsconnection

import (
	"digimon/pbprotocol"
	"fmt"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
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
	SendBuffer    chan []byte
}

func NewConnection(c *websocket.Conn) *Connection {
	nc := &Connection{Conn: c, wg: new(sync.WaitGroup), SendBuffer: make(chan []byte, SENDBUFFERSIZE)}
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
			close(c.SendBuffer)
			wg.Done()
			return
		} else {
			fmt.Println(data)
			err := c.processMsg(data)
			if err != nil {
				log.Println("server internal error!")
				log.Println(err)
			}
			//c.SenndBuffer <- data
		}
	}
}

func (c *Connection) writeLoop(wg *sync.WaitGroup) {
	for {
		select {
		case data, ok := <-c.SendBuffer:
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

// TODO: shoud digimon obj service function
func (c *Connection) processMsg(data []byte) error {
	req := new(pbprotocol.LoginReq)
	err := proto.Unmarshal(data, req)
	if err != nil {
		return err
	}
	log.Printf("login request-----username:%s, password:%s\n", req.Username, req.Password)
	return nil
	//TODO: codec
	//TODO: processFunc
	//TODO: send
}
