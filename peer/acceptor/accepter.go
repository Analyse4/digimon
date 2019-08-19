package acceptor

import (
	"digimon/peer/acceptor/websocket"
	"digimon/service"
)

type Acceptor interface {
	Accept(service.Service)
}

//TODO: Should perfect for general purpose
func Get(typ string) (Acceptor, error) {
	act := new(websocket.Websocket)
	if typ == "ws" {
		return act, nil
	}
	return nil, nil
}
