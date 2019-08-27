package wsconnection

import (
	"digimon/codec"
	"digimon/pbprotocol"
	"digimon/peer/session"
	"digimon/svcregister"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"log"
	"reflect"
	"sync"
)

const SENDBUFFERSIZE = 100

type WSConnection struct {
	ID            int64
	Conn          *websocket.Conn
	wg            *sync.WaitGroup
	ReqDeleteConn chan<- int64
	SendBuffer    chan []byte
}

func NewConnection(c *websocket.Conn) *WSConnection {
	nc := &WSConnection{Conn: c, wg: new(sync.WaitGroup), SendBuffer: make(chan []byte, SENDBUFFERSIZE)}
	return nc
}

func (c *WSConnection) ReadLoop(cd codec.Codec, sess *session.Session) {
	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			close(c.SendBuffer)
			c.wg.Done()
			return
		} else {
			fmt.Println(data)
			err := c.ProcessMsg(data, cd, sess)
			if err != nil {
				log.Println("server internal error!")
				log.Println(err)
			}
			//c.SenndBuffer <- data
		}
	}
}

func (c *WSConnection) WriteLoop() {
	for {
		select {
		case data, ok := <-c.SendBuffer:
			if !ok {
				log.Printf("connection %d write buffer is closed!\n", c.ID)
				c.wg.Done()
				return
			} else {
				err := c.Conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					log.Println(err)
				}
				log.Printf("connection %d send: ", c.ID)
				log.Println(data)
			}
		}
	}
}

func (c *WSConnection) SetID(id int64) {
	c.ID = id
}

func (c *WSConnection) GetID() int64 {
	return c.ID
}

func (c *WSConnection) GetReqDeleteConn() chan<- int64 {
	return c.ReqDeleteConn
}

func (c *WSConnection) SetReqDeleteConn(srd chan<- int64) {
	c.ReqDeleteConn = srd
}

func (c *WSConnection) GetWaitGroup() *sync.WaitGroup {
	return c.wg
}

func (c *WSConnection) ProcessMsg(msg []byte, cd codec.Codec, sess *session.Session) error {
	pack, _ := cd.UnMarshal(msg)
	log.Printf("msg pack-----router:%s, username: %s, password:%s\n", pack.Router, pack.Msg.(*pbprotocol.LoginReq).Username, pack.Msg.(*pbprotocol.LoginReq).Password)
	sess.Set(pack.Msg.(*pbprotocol.LoginReq).Username, pack.Msg.(*pbprotocol.LoginReq).Password)
	fmt.Println(sess.Get(pack.Msg.(*pbprotocol.LoginReq).Username))
	h, err := svcregister.Get(pack.Router)
	if err != nil {
		return err
	}
	f := h.Func
	rv := f.Func.Call([]reflect.Value{h.Receiver, reflect.ValueOf(pack.Msg)})
	if rv[1].Interface() != nil {
		fmt.Println(rv[1].Interface().(error))
	}
	ack, err := proto.Marshal(rv[0].Interface().(proto.Message))
	if err != nil {
		log.Println(err)
	}
	c.SendBuffer <- ack
	return nil
}

func (c *WSConnection) Close() {
	c.Conn.Close()
}
