package gcore

import (
	"fmt"
	"net"

	serialization "dedicate_server/serialization"

	flatbuffers "github.com/google/flatbuffers/go"
)

type GObject struct {
	Id            int
	Name          string
	Pos           Vector2
	Radius        int
	CollisionArea Rectangle
}

func NewGObject(id int, name string, pos Vector2, radius int) *GObject {
	return &GObject{
		Id:   id,
		Name: name,
		Pos:  Vector2{pos.X, pos.Y},
		CollisionArea: Rectangle{
			TopLeft:  Vector2{X: pos.X - radius, Y: pos.Y - radius},
			BotRight: Vector2{X: pos.X + radius, Y: pos.Y + radius},
		},
	}
}

func (obj *GObject) Serialize() ([]byte, int) {
	builder := flatbuffers.NewBuilder(1024)
	name_offset := builder.CreateString(obj.Name)
	serialization.SzGObjectStart(builder)
	serialization.SzGObjectAddId(builder, int32(obj.Id))
	serialization.SzGObjectAddName(builder, name_offset)
	pos_offset := serialization.CreateSzVector2(builder, int32(obj.Pos.X), int32(obj.Pos.Y))
	serialization.SzGObjectAddPos(builder, pos_offset)
	serialization.SzGObjectAddRadius(builder, int32(obj.Radius))
	colision_offset := serialization.CreateSzRectangle(builder,
		int32(obj.CollisionArea.TopLeft.X), int32(obj.CollisionArea.TopLeft.Y),
		int32(obj.CollisionArea.BotRight.X), int32(obj.CollisionArea.BotRight.Y),
	)
	serialization.SzGObjectAddCollisionArea(builder, colision_offset)
	endpos := serialization.SzGObjectEnd(builder)
	builder.Finish(endpos)
	bytes := builder.FinishedBytes()
	return bytes, len(bytes)
}

type Player struct {
	Obj      *GObject
	NickName string
	Addr     net.Addr
}

func (p *Player) ColisionRadius() int {
	return p.Obj.Radius
}
func NewPlayer(usrid int, nick_name string, address net.Addr, position Vector2, colision_radius int) *Player {
	return &Player{
		NewGObject(usrid, nick_name, position, colision_radius),
		nick_name,
		address,
	}
}
func (p *Player) Position() Vector2 {
	return p.Obj.Pos
}
func (w *World) UpdatePlayer(uid string, position Vector2) {
}

func (p Player) GetPositionStr() string {
	return "(" + fmt.Sprintf("%f", p.Obj.Pos.X) + ", " + fmt.Sprintf("%f", p.Obj.Pos.Y) + ")"
}
func (p *Player) UpdatePos(new_pos Vector2) { p.Obj.Pos = new_pos }
