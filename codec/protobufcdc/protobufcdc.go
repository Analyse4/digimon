package protobufcdc

type ProtobufCDC struct{}

func (pbcdc *ProtobufCDC) Marshaler([]byte) error {
	return nil
}

func (pbcdc *ProtobufCDC) UnMarshaler([]byte) error {
	return nil
}
