package gcore

import (
	"container/list"
	"log"
	"strconv"

	"github.com/asim/quadtree"
	// "github.com/asim/quadtree"
)

const (
// DEFAULT_SCREEN_SIZE Vector2 = Vector2{1024, 600}
)

const (
	DEFAULT_COLISION_RADIUS int = 20
	DEFAULT_SIGHT_LEN       int = 20
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

func (w *World) TestInit() {
	// test start
	// case1 := func() {
	// 	enemy_num := 0
	// 	for x := 0; x < 1024; x += 60 {
	// 		for y := 0; y < 600; y += 30 {
	// 			ename := "enemy_" + strconv.Itoa(enemy_num)
	// 			e := NewGObject(enemy_num, ename, Vector2{X: x, Y: y}, DEFAULT_COLISION_RADIUS)
	// 			if !w.object_tree.Insert(e) {
	// 				log.Fatal("Failed to insert the point ", x, " ", y)
	// 				return
	// 			}
	// 			enemy_num++
	// 		}
	// 	}
	// }()
	// case1()
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
	quadtree.NewPoint(float64(obj.Pos.X), float64(obj.Pos.Y), obj)
}
