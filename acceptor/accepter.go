package acceptor

import (
	"digimon/acceptor/websocket"
	"digimon/codec"
)

type Acceptor interface {
	Accept()
}

//TODO: Should perfect for general purpose
func Get(typ string, addr string, codecTyp string) (Acceptor, error) {
	act := new(websocket.Websocket)
	if typ == "ws" && codecTyp == "protobufcdc" {
		act.Addr = addr
		act.Codec, _ = codec.Get(codecTyp)
		return act, nil
	}
	return nil, nil
}
