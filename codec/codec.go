package codec

import "digimon/codec/protobufcdc"

type Codec interface {
	Marshaler([]byte) error
	UnMarshaler([]byte) error
}

//TODO
func Get(typ string) (Codec, error) {
	if typ == "protobufcdc" {
		return new(protobufcdc.ProtobufCDC), nil
	}
	return nil, nil
}
