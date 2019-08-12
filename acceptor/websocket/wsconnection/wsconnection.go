package wsconnection

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Connection struct {
	Id   int64
	Conn *websocket.Conn
}

func NewConnection(c *websocket.Conn) *Connection {
	nc := &Connection{Conn: c}
	go nc.read()
	return nc

}

func (c *Connection) read() {
	for {
		_, data, err := c.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Println(data)
		}
	}
}
