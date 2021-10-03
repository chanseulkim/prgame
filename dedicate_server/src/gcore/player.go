package gcore

import (
	"fmt"
	"net"
)

type Player struct {
	Obj           *GObject
	NickName      string
	Addr          net.Addr
	Index_inworld int
}

func (p *Player) ColisionRadius() int {
	return p.Obj.Radius
}
func NewPlayer(usrid int, nick_name string, address net.Addr, position Vector2, colision_radius int) *Player {
	return &Player{
		NewGObject(usrid, nick_name, position, colision_radius),
		nick_name,
		address,
		0.0,
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
