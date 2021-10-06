package gcore

import (
	"container/list"
	. "dedicate_server/gnet"
	"fmt"
	"log"
	"net"
	"strconv"
)

const (
// DEFAULT_SCREEN_SIZE Vector2 = Vector2{1024, 600}
)

const (
	DEFAULT_COLISION_RADIUS int = 20
	DEFAULT_SIGHT_LEN       int = 20
)

type World struct {
	Players     map[string]*Player // addr, player
	screen_size Vector2
	object_tree *QuadNode
	// object_tree *quadtree.QuadTree
}

var world_instance *World

func (w *World) Init() {
	w.TestInit()
}

func (w *World) UpdatePlayer(uid string, position Vector2) {
}

func (w *World) TestInit() {
	case2 := func() {
		enemy_num := 0
		for x := 0; x < 1024; x += 100 {
			for y := 0; y < 600; y += 60 {
				ename := "enemy_" + strconv.Itoa(enemy_num)
				e := NewGObject(enemy_num, ename, Vector2{X: x, Y: y}, DEFAULT_COLISION_RADIUS)
				if !w.object_tree.Insert(e) {
					log.Fatal("Failed to insert the point ", x, " ", y)
					return
				}
				enemy_num++
				if enemy_num == 100 {
					return
				}
			}
		}
	}
	case2()
	// test end
}

func (w *World) Nearest(player *Player) *list.List {
	founds := w.object_tree.Nearest2(player.Position(), player.Obj.Radius)
	// for _, point := range founds {
	// 	log.Printf("Found point: %s\n", point.Data().(string))
	// }
	return founds
}
func (w *World) GetAllObjects() *list.List {
	return w.object_tree.GetAllObjects()
}

func GetWorld() *World {
	if world_instance == nil {
		world_instance = &World{
			Players:     make(map[string]*Player),
			screen_size: Vector2{1024, 600},
			object_tree: NewQuadTreeRoot(1024, 600),
		}
		world_instance.Init()
	}
	return world_instance
}

func (w *World) AddObject(obj *GObject) {
	w.object_tree.Insert(obj)
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
func (p Player) GetPositionStr() string {
	return "(" + fmt.Sprintf("%f", p.Obj.Pos.X) + ", " + fmt.Sprintf("%f", p.Obj.Pos.Y) + ")"
}
func (p *Player) UpdatePos(new_pos Vector2) { p.Obj.Pos = new_pos }
