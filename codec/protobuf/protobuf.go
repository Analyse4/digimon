package protobuf

import (
	"fmt"
	"github.com/Analyse4/digimon/pbprotocol"
	"github.com/Analyse4/digimon/svcregister"
	"github.com/golang/protobuf/proto"
	"reflect"
)

type Pack struct {
	Router string
	Msg    interface{}
}

type Protobuf struct{}

func (pbcdc *Protobuf) Marshal(router string, msg interface{}) ([]byte, error) {
	m, ok := msg.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("invaild data")
	}
	data, err := proto.Marshal(m)
	if err != nil {
		return nil, err
	}
	bm := new(pbprotocol.MsgPack)
	bm.Router = router
	bm.Data = data
	ack, err := proto.Marshal(bm)
	return ack, nil
}

func (pbcdc *Protobuf) UnMarshal(msg []byte) (*Pack, error) {
	mp := new(pbprotocol.MsgPack)
	err := proto.Unmarshal(msg, mp)
	if err != nil {
		return nil, err
	}
	router := mp.Router
	handler, err := svcregister.Get(router)
	if err != nil {
		return nil, err
	}
	req := reflect.New(handler.Typ.Elem()).Interface()
	err = proto.Unmarshal(mp.Data, req.(proto.Message))
	if err != nil {
		return nil, err
	}
	pack := new(Pack)
	pack.Router = router
	pack.Msg = req
	return pack, nil
}
