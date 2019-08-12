package codec

import "digimon/codec/protobufcdc"

type Codec interface {
	Marshaler([]byte) error
	UnMarshaler([]byte) error
}

//TODO: Should perfect for general purpose
func Get(typ string) (Codec, error) {
	if typ == "protobufcdc" {
		return new(protobufcdc.ProtobufCDC), nil
	}
	return nil, nil
}
