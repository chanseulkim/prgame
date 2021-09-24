package libs

import (
	"math"
)

const LEAST_BLOCKSIZE = 1
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
	topleft_pnt   Vector2
	botright_pnt  Vector2

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
		tlp := Vector2{X: Float(startc), Y: Float(startr)}
		brp := Vector2{X: Float(endc - 1), Y: Float(endr - 1)}
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
	var inBoundary = func(p Vector2) bool {
		return (p.X >= self.topleft_pnt.X && p.X <= self.botright_pnt.X &&
			p.Y >= self.topleft_pnt.Y && p.Y <= self.botright_pnt.Y)
	}
	if !inBoundary(new_obj.Pos) {
		return
	}
	tlp := self.topleft_pnt
	brp := self.botright_pnt
	if (math.Abs(float64(tlp.X-brp.X)) <= LEAST_BLOCKSIZE) &&
		math.Abs(float64(tlp.Y-brp.Y)) <= LEAST_BLOCKSIZE {
		self.append_object(new_obj)
		self.IsLeaf = true
		return
	}

	if new_obj.Pos.X < (tlp.X+brp.X)/2 { // left
		if new_obj.Pos.Y < (tlp.Y+brp.Y)/2 { // top left
			if self.TopLeft == nil {
				self.TopLeft = &QuadNode{
					topleft_pnt: Vector2{
						self.topleft_pnt.X,
						self.topleft_pnt.Y,
					},
					botright_pnt: Vector2{
						((self.topleft_pnt.X + self.botright_pnt.X + 1) / 2) - 1, // 왼쪽이기때문에 마지막 -1
						((self.topleft_pnt.Y + self.botright_pnt.Y + 1) / 2) - 1,
					},
					Width:  self.Width / 2,
					Height: self.Height / 2,
				}
			}
			self.TopLeft.Insert(new_obj)
		} else { // bottom left
			if self.BottomLeft == nil {
				self.BottomLeft = &QuadNode{
					topleft_pnt: Vector2{
						self.topleft_pnt.X,
						(self.topleft_pnt.Y + self.botright_pnt.Y + 1) / 2,
					},
					botright_pnt: Vector2{
						((self.topleft_pnt.X + self.botright_pnt.X + 1) / 2) - 1,
						self.botright_pnt.Y,
					},
					Width:  self.Width / 2,
					Height: self.Height / 2,
					IsLeaf: false,
					Id:     new_obj.Id,
				}
			}
			self.BottomLeft.Insert(new_obj)
		}
	} else { // right
		if new_obj.Pos.Y <= (tlp.Y+brp.Y)/2 { // top right
			if (self.TopRight == nil) || (self.TopRight.IsLeaf) {
				self.TopRight = &QuadNode{
					topleft_pnt: Vector2{
						((self.topleft_pnt.X + self.botright_pnt.X + 1) / 2),
						self.topleft_pnt.X,
					},
					botright_pnt: Vector2{
						self.botright_pnt.X,
						((self.topleft_pnt.Y + self.botright_pnt.Y + 1) / 2) - 1, // 위쪽이기 때문에 마지막 -1
					},
					Width:  self.Width / 2,
					Height: self.Height / 2,
				}
			}
			self.TopRight.Insert(new_obj)
		} else { // bottom right
			if (self.BottomRight == nil) || (self.BottomRight.IsLeaf) {
				self.BottomRight = &QuadNode{
					topleft_pnt: Vector2{
						(self.topleft_pnt.X + self.botright_pnt.X + 1) / 2,
						(self.topleft_pnt.Y + self.botright_pnt.Y + 1) / 2,
					},
					botright_pnt: Vector2{
						self.botright_pnt.X,
						self.botright_pnt.Y,
					},
					Width:  self.Width / 2,
					Height: self.Height / 2,
				}
			}
			self.BottomRight.Insert(new_obj)
		}
	}
}

func (self *QuadNode) GetNear(target_obj *GObject) []*GObject {
	var inBoundary = func(p Vector2) bool {
		return (p.X >= self.topleft_pnt.X && p.X <= self.botright_pnt.X &&
			p.Y >= self.topleft_pnt.Y && p.Y <= self.botright_pnt.Y)
	}
	if !inBoundary(target_obj.Pos) {
		return nil
	}
	tlp := self.topleft_pnt
	brp := self.botright_pnt
	if (math.Abs(float64(tlp.X-brp.X)) <= LEAST_BLOCKSIZE) &&
		math.Abs(float64(tlp.Y-brp.Y)) <= LEAST_BLOCKSIZE {
		return self.obj
	}

	if target_obj.Pos.X < (tlp.X+brp.X)/2 { // left
		if target_obj.Pos.Y < (tlp.Y+brp.Y)/2 { // top left
			self.TopLeft.GetNear(target_obj)
		} else { // bottom left
			self.BottomLeft.GetNear(target_obj)
		}
	} else { // right
		if target_obj.Pos.Y <= (tlp.Y+brp.Y)/2 { // top right
			self.TopRight.GetNear(target_obj)
		} else { // bottom right
			self.BottomRight.GetNear(target_obj)
		}
	}
}

func whereBlock(x int, y int, grid_r int, clen int) int {

	return -1
}
