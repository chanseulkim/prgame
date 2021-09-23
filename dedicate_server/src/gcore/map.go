package gcore

import (
	"math"
)

const (
	NODE_TOPLEFT int = 0
	NODE_TOPRIGHT
	NODE_BOTTOMLEFT
	NODE_BOTTOMRIGHT
)

type QuadNode struct {
	Id            int
	IsLeaf        bool
	Width, Height int
	topleft_pnt   Point
	botright_pnt  Point

	obj []*GObject

	TopLeft     *QuadNode
	TopRight    *QuadNode
	BottomLeft  *QuadNode
	BottomRight *QuadNode
}

func ConstructQuadTree(grid [][]int) *QuadNode {
	var construct_task func(startr int, endr int, startc int, endc int, grid [][]int) *QuadNode
	construct_task = func(startr int, endr int, startc int, endc int, grid [][]int) *QuadNode {
		val := grid[startr][startc]
		var isleaf = func() bool {
			for r := startr; r < endr; r++ {
				for c := startc; c < endc; c++ {
					if grid[r][c] != val {
						return false
					}
				}
			}
			return true
		}
		tlp := Point{X: Float(startc), Y: Float(startr)}
		brp := Point{X: Float(endc), Y: Float(endr)}
		if isleaf() {
			return &QuadNode{
				Id:           val,
				IsLeaf:       true,
				Width:        endc - startc,
				Height:       endr - startr,
				topleft_pnt:  tlp,
				botright_pnt: brp,
			}
		}
		new_node := &QuadNode{
			Id:           val,
			IsLeaf:       false,
			Width:        endc - startc,
			Height:       endr - startr,
			topleft_pnt:  tlp,
			botright_pnt: brp,
		}
		midr := startr + (endr-startr)/2
		midc := startc + (endc-startc)/2
		new_node.TopLeft = construct_task(startr, midr, startc, midc, grid)
		new_node.TopRight = construct_task(startr, midr, midc, endc, grid)
		new_node.BottomLeft = construct_task(midr, endr, startc, midc, grid)
		new_node.BottomRight = construct_task(midr, endr, midc, endc, grid)
		return new_node
	}
	return construct_task(0, len(grid), 0, len(grid[0]), grid)
}
func (self *QuadNode) append_object(new_obj *GObject) {
	self.obj = append(self.obj, new_obj)
}
func (self *QuadNode) Insert(new_obj *GObject) {
	var inBoundary = func(obj *GObject) bool {
		var p = obj.Pos
		return (p.X >= self.topleft_pnt.X && p.X <= self.botright_pnt.X &&
			p.Y >= self.topleft_pnt.Y && p.Y <= self.botright_pnt.Y)
	}
	if !inBoundary(new_obj) {
		return
	}
	tlp := self.topleft_pnt
	brp := self.botright_pnt
	if (math.Abs(float64(tlp.X-brp.X)) <= 1) &&
		math.Abs(float64(tlp.Y-brp.Y)) <= 1 {
		if self.obj == nil {
			self.append_object(new_obj)
		}
		return
	}

	if new_obj.Pos.X <= (tlp.X+brp.X)/2 { // left
		if new_obj.Pos.Y <= (tlp.Y+brp.Y)/2 { // top left
			if self.TopLeft == nil {
				self.TopLeft = &QuadNode{
					topleft_pnt: Point{
						self.topleft_pnt.X,
						self.topleft_pnt.Y,
					},
					botright_pnt: Point{
						(self.topleft_pnt.X + self.botright_pnt.X) / 2,
						(self.topleft_pnt.Y + self.botright_pnt.Y) / 2,
					},
				}
			}
			self.TopLeft.append_object(new_obj)
		} else { // bottom left
			if self.BottomLeft == nil {
				self.BottomLeft = &QuadNode{
					topleft_pnt: Point{
						self.topleft_pnt.X,
						(self.topleft_pnt.Y + self.botright_pnt.Y) / 2,
					},
					botright_pnt: Point{
						(self.topleft_pnt.X + self.botright_pnt.X) / 2,
						self.botright_pnt.Y,
					},
				}
			}
			self.BottomLeft.append_object(new_obj)
		}
	}

}

func (self *QuadNode) Search(r int, c int, rlen int, clen int, node *QuadNode) {
	if !node.IsLeaf {
		whereBlock(r, c, rlen, clen)
	}
}
func whereBlock(r int, c int, rlen int, clen int) int {
	if (r < (rlen / 2)) && (c < (clen / 2)) {
		return NODE_TOPLEFT
	} else if (r < (rlen / 2)) && (c >= (clen / 2)) {
		return NODE_TOPRIGHT
	} else if (r >= (rlen / 2)) && (c < (clen / 2)) {
		return NODE_BOTTOMLEFT
	} else if (r >= (rlen / 2)) && (c >= (clen / 2)) {
		return NODE_TOPRIGHT
	}
	return -1
}

// TODO: 좌표값은 Float이고 배열은 정수 index인데, 어떻게 배열로 처리해야할지..
type GameMap struct {
	area     Area
	grid     [][]int
	quadTree *QuadNode
}

func (g *GameMap) Init() {
	g.quadTree = ConstructQuadTree(g.grid)
}

// TODO: collision shape를 등록해야함
func (m *GameMap) AddObject(obj *GObject) {
	m.area[int(obj.Pos.X)][int(obj.Pos.Y)] = obj.Name
}
func (m *GameMap) GetArea() Area { return m.area }
