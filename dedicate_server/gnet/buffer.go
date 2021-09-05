package gnet

type MsgBuff struct {
	data     []byte
	size     uint32
	capacity uint32
}

func NewMsgBuff(data []byte, data_size uint32) *MsgBuff {
	var buff MsgBuff
	buff.data = make([]byte, data_size)
	copy(buff.data, data)
	buff.size = data_size
	return &buff
}

func (self *MsgBuff) Write(data []byte, data_size uint32) {
	copy(self.data[self.size:], data)
	self.size += data_size
}
