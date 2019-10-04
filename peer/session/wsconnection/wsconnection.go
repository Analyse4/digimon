package wsconnection

import (
	"digimon/codec"
	"digimon/logger"
	"digimon/peer/cleaner"
	"digimon/peer/session"
	"digimon/svcregister"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"reflect"
	"sync"
)

//TODO: buffer size should bigger
const SENDBUFFERSIZE = 100

var (
	log *logrus.Entry
)

func init() {
	log = logger.GetLogger().WithField("pkg", "wsconnection")
}

type WSConnection struct {
	ID            int64
	Conn          *websocket.Conn
	wg            *sync.WaitGroup
	ReqDeleteConn chan<- *cleaner.CleanerMeta
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
			log.WithFields(logrus.Fields{
				"connection_id": c.ID,
			}).Debug(err)

			close(c.SendBuffer)
			c.wg.Done()

			log.WithFields(logrus.Fields{
				"connection_id": c.ID,
			}).Debug("read loop finished")

			return
		} else {
			log.WithFields(logrus.Fields{
				"connection_id": c.ID,
				"data_len":      len(data),
			}).Debug("receive data")

			c.ProcessMsg(data, cd, sess)
		}
	}
}

func (c *WSConnection) WriteLoop() {
	for {
		select {
		case data, ok := <-c.SendBuffer:
			if !ok {
				c.wg.Done()

				log.WithFields(logrus.Fields{
					"connection_id": c.ID,
				}).Debug("write loop finished")

				return
			} else {
				err := c.Conn.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					log.WithFields(logrus.Fields{
						"connection_id": c.ID,
					}).Error(err)
				}

				log.WithFields(logrus.Fields{
					"connection_id": c.ID,
					"data_len":      len(data),
				}).Debug("send data")
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

func (c *WSConnection) GetReqDeleteConn() chan<- *cleaner.CleanerMeta {
	return c.ReqDeleteConn
}

func (c *WSConnection) SetReqDeleteConn(srd chan<- *cleaner.CleanerMeta) {
	c.ReqDeleteConn = srd
}

func (c *WSConnection) GetWaitGroup() *sync.WaitGroup {
	return c.wg
}

func (c *WSConnection) ProcessMsg(msg []byte, cd codec.Codec, sess *session.Session) {
	pack, err := cd.UnMarshal(msg)
	if err != nil {
		log.WithFields(logrus.Fields{
			"func": "unmarshal",
		}).Error(err)
		return
	}
	h, err := svcregister.Get(pack.Router)
	if err != nil {
		log.WithFields(logrus.Fields{
			"router": pack.Router,
		}).Error("router not found")
		return
	}
	f := h.Func
	rv := f.Func.Call([]reflect.Value{h.Receiver, reflect.ValueOf(sess), reflect.ValueOf(pack.Msg)})
	if rv[1].Interface() != nil {
		log.WithFields(logrus.Fields{
			"func": f.Name,
		}).Error(rv[1].Interface().(error))
		return
	}
	//TODO: use codec marshal
	//ack, err := proto.Marshal(rv[0].Interface().(proto.Message))
	ack, _ := cd.Marshal(pack.Router, rv[0].Interface())
	if err != nil {
		log.WithFields(logrus.Fields{
			"func": "marshal",
		}).Error(err)
	}
	c.SendBuffer <- ack
}

func (c *WSConnection) Close() {
	c.Conn.Close()
}

func (c *WSConnection) Send(data []byte) {
	c.SendBuffer <- data
}
