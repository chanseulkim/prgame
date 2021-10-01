package gnet

import (
	"fmt"

	flatbuffers "github.com/google/flatbuffers/go"
)

type MsgBuff struct {
	data     []byte
	size     int32
	capacity uint32
}

func NewMsgBuff(data []byte, data_size int32) *MsgBuff {
	var buff MsgBuff
	buff.data = make([]byte, data_size)
	copy(buff.data, data)
	buff.size = data_size
	return &buff
}

func (self *MsgBuff) Write(data []byte, data_size int32) {
	copy(self.data[self.size:], data)
	self.size += data_size
}
func Serialize(buff *MsgBuff) {
	builder := flatbuffers.NewBuilder(1024)
	weaponOne := builder.CreateString("Sword")
	fmt.Println(weaponOne)
	builder.Finish(weaponOne)
}
