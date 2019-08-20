package protobuf

import (
	"digimon/pbprotocol"
	"digimon/svcregister"
	"github.com/golang/protobuf/proto"
	"reflect"
	"traefik/log"
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
	reqTyp, ok := svcregister.SVCR.Register["digimon.login"]
	if !ok {
		log.Println("service request type is not registered!")
	}
	req := reflect.New(reqTyp).Interface()
	err = proto.Unmarshal(mp.Data, req.(proto.Message))
	pack := new(Pack)
	pack.Router = router
	pack.Msg = req
	return pack, nil
}
