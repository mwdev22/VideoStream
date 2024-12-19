package server

type Encoding int

const (
	Json Encoding = iota
	XML
	Protobuf
)

type Packet struct {
	IP   string
	Data []byte
	Enc  Encoding
}
