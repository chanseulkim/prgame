package gcore

import (
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

type GObject struct {
	Id             int
	Name           string
	Pos            Vector2
	ColisionRadius Rectangle
}

type Player struct {
	Uid             string
	Addr            net.Addr
	position        Vector2
	Colision_radius Float
	Index_inworld   Float
}

func NewPlayer(uid string, address net.Addr, position Vector2, colision_radius Float) *Player {
	return &Player{uid, address, position, colision_radius, 0.0}
}

func (p Player) GetPositionStr() string {
	return "(" + fmt.Sprintf("%f", p.position.X) + ", " + fmt.Sprintf("%f", p.position.Y) + ")"
}
func (p Player) GetPosition() Vector2       { return p.position }
func (p *Player) UpdatePos(new_pos Vector2) { p.position = new_pos }

type Area [][]string

func NewGameMap(x int, y int) *GameMap {
	var new_space = make(Area, x)
	for i := range new_space {
		new_space[i] = make([]string, y)
	}
	return &GameMap{
		area: new_space,
	}
}

type World struct {
	Players     map[string]*Player // addr, player
	objects     []*GObject
	screen_size Vector2
	world_map   GameMap
}

func (w *World) AddObject(object *GObject) {
	w.objects = append(w.objects, object)
	w.world_map.AddObject(object)
}
func (w *World) GetMapArea() Area { return w.world_map.GetArea() }

// func (w *World) SetScreenSize(screen_size Vector2)         { w.screen_size = screen_size }
// func (w *World) GetScreenSize(screen_size Vector2) Vector2 { return w.screen_size }
func (w *World) GetObjects() []*GObject { return w.objects }

func (w *World) Init() {
	w.world_map = NewGameMap(int(w.screen_size.X), int(w.screen_size.Y))

	w.AddObject(
		&GObject{Name: "enemy", Pos: Vector2{400, 300}},
	)

}

var world_instance *World

func GetWorld() *World {
	if world_instance == nil {
		world_instance = &World{
			Players:     make(map[string]*Player),
			screen_size: Vector2{1024, 600},
		}
		world_instance.Init()
	}
	return world_instance
}
