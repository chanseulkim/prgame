package gcore

import (
	"container/list"
	"fmt"
	"net"
)

const (
// DEFAULT_SCREEN_SIZE Vector2 = Vector2{1024, 600}
)

const (
	DEFAULT_COLISION_RADIUS Float = 20.0
	DEFAULT_SIGHT_LEN       Float = 20.0
)

type Player struct {
	Uid             string
	Addr            net.Addr
	position        Vector2
	Colision_radius Float
	Index_inworld   Float
}

func (p Player) GetPositionStr() string {
	x := fmt.Sprintf("%f", p.position.X)
	y := fmt.Sprintf("%f", p.position.Y)
	return "(" + x + ", " + y + ")"
}
func (p Player) GetPosition() Vector2       { return p.position }
func (p *Player) UpdatePos(new_pos Vector2) { p.position = new_pos }

func NewPlayer(uid string, address net.Addr, position Vector2, colision_radius Float) *Player {
	return &Player{uid, address, position, colision_radius, 0.0}
}

type World struct {
	Players     map[string]*Player // addr, player
	objects     *list.List
	position    Vector2
	screen_size Vector2
}

func (w *World) Init() {

}

func (w *World) SetScreenSize(screen_size Vector2)         { w.screen_size = screen_size }
func (w *World) GetScreenSize(screen_size Vector2) Vector2 { return w.screen_size }

func (w *World) AddObject(object *GObject) {
	w.objects.PushBack(object)
}

var world_instance *World

func GetWorld() *World {
	if world_instance == nil {
		world_instance = &World{
			Players: make(map[string]*Player),
			objects: list.New(),
		}
		world_instance.Init()
	}
	return world_instance
}
