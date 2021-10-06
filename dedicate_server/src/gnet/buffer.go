package gnet

type MsgBuff struct {
	Data []byte
	Size int
}

func NewMsgBuff(data []byte, data_size int) *MsgBuff {
	var buff MsgBuff
	buff.Data = make([]byte, data_size)
	copy(buff.Data, data)
	buff.Size = data_size
	return &buff
}

func (self *MsgBuff) Write(data []byte, data_size int) {
	copy(self.Data[self.Size:], data)
	self.Size += data_size
}

// func Serialize(obj *gcore.GObject) []byte {
// 	builder := flatbuffers.NewBuilder(1024)

// 	name_offset := builder.CreateString(obj.Name)
// 	serialization.SzGObjectStart(builder)
// 	serialization.SzGObjectAddId(builder, int32(obj.Id))
// 	serialization.SzGObjectAddName(builder, name_offset)
// 	pos_offset := serialization.CreateSzVector2(builder, int32(obj.Pos.X), int32(obj.Pos.Y))
// 	serialization.SzGObjectAddPos(builder, pos_offset)
// 	serialization.SzGObjectAddRadius(builder, int32(obj.Radius))
// 	colision_offset := serialization.CreateSzRectangle(builder,
// 		int32(obj.CollisionArea.TopLeft.X), int32(obj.CollisionArea.TopLeft.Y),
// 		int32(obj.CollisionArea.BotRight.X), int32(obj.CollisionArea.BotRight.Y),
// 	)
// 	serialization.SzGObjectAddCollisionArea(builder, colision_offset)

// 	endpos := serialization.SzGObjectEnd(builder)
// 	builder.Finish(endpos)
// 	return builder.FinishedBytes()
// }
