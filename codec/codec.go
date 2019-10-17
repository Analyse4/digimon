package codec

import (
	"github.com/Analyse4/digimon/codec/protobuf"
)

type Codec interface {
	Marshal(string, interface{}) ([]byte, error)
	UnMarshal([]byte) (*protobuf.Pack, error)
}

//TODO: Should perfect for general purpose
func Get(typ string) (Codec, error) {
	if typ == "protobuf" {
		return new(protobuf.Protobuf), nil
	}
	return nil, nil
}
