package svcregister

import (
	"fmt"
	"reflect"
	"sync"
)

type Handler struct {
	Typ      reflect.Type
	Func     reflect.Method
	Receiver reflect.Value
}

var SVCRegister sync.Map

func init() {
	SVCRegister = sync.Map{}
}

func Get(index string) (*Handler, error) {
	h, ok := SVCRegister.Load(index)
	if !ok {
		return nil, fmt.Errorf("get handler failed-----index: %s", index)
	}
	if h == nil {
		return nil, fmt.Errorf("handler is not registed----index: %s", index)
	}
	return h.(*Handler), nil
}

func Set(index string, handler *Handler) {
	SVCRegister.Store(index, handler)
}
