package svcregister

import (
	"digimon/pbprotocol"
	"reflect"
)

type SVCRegister struct {
	Register map[string]reflect.Type
}

var SVCR *SVCRegister

func init() {
	SVCR = new(SVCRegister)
	SVCR.Register = make(map[string]reflect.Type)
	SVCR.Register["digimon.login"] = reflect.TypeOf(pbprotocol.LoginReq{})
}
