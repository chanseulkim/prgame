package gcore

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/asim/quadtree"
)

const (
// DEFAULT_SCREEN_SIZE Vector2 = Vector2{1024, 600}
)

const (
	DEFAULT_COLISION_RADIUS int = 200
	DEFAULT_SIGHT_LEN       int = 200
)

type GObject struct {
	Id            int
	Name          string
	Pos           Vector2
	SightRadius   int
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

type Player struct {
	Uid            string
	Addr           net.Addr
	position       Vector2
	ColisionRadius int
	Index_inworld  int
}

func NewPlayer(uid string, address net.Addr, position Vector2, colision_radius int) *Player {
	return &Player{uid, address, position, colision_radius, 0.0}
}

func (p Player) GetPositionStr() string {
	return "(" + fmt.Sprintf("%f", p.position.X) + ", " + fmt.Sprintf("%f", p.position.Y) + ")"
}
func (p Player) GetPosition() Vector2       { return p.position }
func (p *Player) UpdatePos(new_pos Vector2) { p.position = new_pos }

type World struct {
	Players     map[string]*Player // addr, player
	screen_size Vector2
	//object_tree *QuadNode
	object_tree *quadtree.QuadTree
}

var world_instance *World

func (w *World) Init() {
	centerPoint := quadtree.NewPoint(0.0, 0.0, nil)
	halfPoint := quadtree.NewPoint(1024, 600, nil)
	boundingBox := quadtree.NewAABB(centerPoint, halfPoint)
	w.object_tree = quadtree.New(boundingBox, 0, nil)

	enemy_num := 1
	for x := 0; x < 1024; x += 60 {
		for y := 0; y < 600; y += 30 {
			ename := "enemy_" + strconv.Itoa(enemy_num)
			e := quadtree.NewPoint(float64(x), float64(y), ename)
			if !w.object_tree.Insert(e) {
				log.Fatal("Failed to insert the point")
				return
			}
			enemy_num++
		}
	}
}

func (w *World) Nearest(player *Player) []*quadtree.Point {
	center := quadtree.NewPoint(float64(player.position.X), float64(player.position.Y), nil)
	//distance := 10000000 // 100단위
	//distance := 8000000 // 80단위
	distance := 5000000 // 50단위
	// distance := 1000000 // 10단위
	bounds := quadtree.NewAABB(center, center.HalfPoint(float64(distance)))
	maxPoints := 10
	founds := w.object_tree.KNearest(bounds, maxPoints, nil)
	// for _, point := range founds {
	// 	log.Printf("Found point: %s\n", point.Data().(string))
	// }
	return founds
}

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
func (w *World) AddObject(obj *GObject) {
	quadtree.NewPoint(float64(obj.Pos.X), float64(obj.Pos.Y), obj)
}
