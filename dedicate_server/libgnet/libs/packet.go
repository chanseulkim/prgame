package gnet

// import "dedicate_server/gcore"

const (
	HEADERTYPE_MOVE  byte = 0
	HEADERTYPE_SYNC  byte = 1
	HEADERTYPE_ENTER byte = 2
	HEADERTYPE_LEAVE byte = 3
	HEADERTYPE_MSG   byte = 4
)

type GHeader struct {
	header_type byte
}

type GPacket struct {
	header   *GHeader
	buff     []byte
	buff_len int32
}

func NewGPacket(pack_type byte, data []byte, data_size int32) *GPacket {
	return &GPacket{
		header:   &GHeader{header_type: pack_type},
		buff:     data,
		buff_len: data_size,
	}
}

func (p *GPacket) Write(data []byte, data_size int32) {
	copy(p.buff[p.buff_len:], data)
	p.buff_len += data_size
}
