package protobuf

import (
	"digimon/pbprotocol"
	"digimon/svcregister"
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"reflect"
)

type Pack struct {
	Router string
	Msg    interface{}
}

type Protobuf struct{}

func (pbcdc *Protobuf) Marshal(msg []byte) error {
	return nil
}

func (pbcdc *Protobuf) UnMarshal(msg []byte) (*Pack, error) {
	mp := new(pbprotocol.MsgPack)
	err := proto.Unmarshal(msg, mp)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	router := mp.Router
	handler, err := svcregister.Get(router)
	if err != nil {
		log.Println("service request type is not registered!")
		//return nil, err
	}
	fmt.Println(handler)
	req := reflect.New(handler.Typ.Elem()).Interface()
	err = proto.Unmarshal(mp.Data, req.(proto.Message))
	pack := new(Pack)
	pack.Router = router
	pack.Msg = req
	return pack, nil
}
