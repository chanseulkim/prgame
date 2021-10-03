package gcore

import (
	"fmt"
	"math"
)

const LEAST_BLOCKSIZE = 1
const (
	NODE_TOPLEFT     = 0
	NODE_TOPRIGHT    = 1
	NODE_BOTTOMLEFT  = 2
	NODE_BOTTOMRIGHT = 3
	NODE_ROOT        = 4
)

type QuadNode struct {
	Id            int
	is_leaf       bool
	Width, Height int
	topleft_pnt   Vector2
	botright_pnt  Vector2
	objs          []*GObject
	// objs          *list.List
	NodeSector int

	TopLeft     *QuadNode
	TopRight    *QuadNode
	BottomLeft  *QuadNode
	BottomRight *QuadNode
	Parent      *QuadNode
}

func NewQuadNode(
	id int,
	is_leaf bool,
	width, height int,
	topleft_pnt, botright_pnt Vector2,
	parent *QuadNode,
	node_sector int) *QuadNode {
	return &QuadNode{
		Id:      id,
		is_leaf: true,
		// objs:         list.New(),
		Width:        width,
		Height:       height,
		topleft_pnt:  topleft_pnt,
		botright_pnt: botright_pnt,
		Parent:       parent,
		NodeSector:   node_sector,
	}
}
func (self *QuadNode) IsLeaf() bool {
	return (self.TopLeft == nil) &&
		(self.TopRight == nil) &&
		(self.BottomLeft == nil) &&
		(self.BottomRight == nil)
}

// 자식 노드들이 가진 모든 오브젝트들 반환
func (self *QuadNode) GetAllObjects() []*GObject {
	var objs []*GObject = self.objs
	if self.TopLeft != nil {
		if !self.TopLeft.IsLeaf() {
			objs = append(objs, self.TopLeft.GetAllObjects()...)
		} else {
			objs = append(objs, self.TopLeft.objs...)
		}
	}
	if self.TopRight != nil {
		if !self.TopRight.IsLeaf() {
			objs = append(objs, self.TopRight.GetAllObjects()...)
		} else {
			objs = append(objs, self.TopRight.objs...)
		}
	}
	if self.BottomLeft != nil {
		if !self.BottomLeft.IsLeaf() {
			objs = append(objs, self.BottomLeft.GetAllObjects()...)
		} else {
			objs = append(objs, self.BottomLeft.objs...)
		}
	}
	if self.BottomRight != nil {
		if !self.BottomRight.IsLeaf() {
			objs = append(objs, self.BottomRight.GetAllObjects()...)
		} else {
			objs = append(objs, self.BottomRight.objs...)
		}
	}
	return objs
}
func NewQuadTreeRoot(x int, y int) *QuadNode {
	return NewQuadNode(0, true, int(x), int(y), Vector2{X: 0, Y: 0}, Vector2{X: x, Y: y}, nil, NODE_ROOT)
}
func ConstructQuadTree(grid [][]int) *QuadNode {
	var construct_task func(startr int, endr int, startc int, endc int, grid [][]int, parent *QuadNode, node_sector int) *QuadNode
	construct_task = func(startr int, endr int, startc int, endc int, grid [][]int, parent *QuadNode, node_sector int) *QuadNode {
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
		tlp := Vector2{X: int(startc), Y: int(startr)}
		brp := Vector2{X: int(endc - 1), Y: int(endr - 1)}
		if isleaf() {
			return NewQuadNode(val, true, endc-startc, endr-startr, tlp, brp, parent, node_sector)
		}
		new_node := NewQuadNode(val, false,
			endc-startc, endr-startr,
			tlp, brp,
			parent, node_sector)
		parent = new_node
		midr := startr + (endr-startr)/2
		midc := startc + (endc-startc)/2
		new_node.TopLeft = construct_task(startr, midr, startc, midc, grid, parent, NODE_TOPLEFT)
		new_node.TopRight = construct_task(startr, midr, midc, endc, grid, parent, NODE_TOPRIGHT)
		new_node.BottomLeft = construct_task(midr, endr, startc, midc, grid, parent, NODE_BOTTOMLEFT)
		new_node.BottomRight = construct_task(midr, endr, midc, endc, grid, parent, NODE_BOTTOMRIGHT)
		return new_node
	}
	return construct_task(0, len(grid), 0, len(grid[0]), grid, nil, NODE_ROOT)
}

func (self *QuadNode) append_object(new_obj *GObject) {
	self.objs = append(self.objs, new_obj)
}

func (self *QuadNode) Insert(new_obj *GObject) bool {
	if self == nil {
		return false
	}
	var inBoundary = func(p *Vector2) bool {
		return (p.X >= self.topleft_pnt.X && p.X <= self.botright_pnt.X+1 &&
			p.Y >= self.topleft_pnt.Y && p.Y <= self.botright_pnt.Y+1)
	}
	if !inBoundary(&new_obj.Pos) {
		fmt.Println("Error ", new_obj.Name, " is not in boundary")
		return false
	}
	tlp := self.topleft_pnt
	brp := self.botright_pnt
	if new_obj.Pos.X < (tlp.X+brp.X)/2 { // left
		if new_obj.Pos.Y < (tlp.Y+brp.Y)/2 { // top left
			if self.TopLeft == nil {
				self.TopLeft = NewQuadNode(new_obj.Id, true,
					self.Width/2, self.Height/2,
					Vector2{self.topleft_pnt.X, self.topleft_pnt.Y},
					Vector2{((self.topleft_pnt.X + self.botright_pnt.X + 1) / 2) - 1, // 왼쪽이기때문에 마지막 -1
						((self.topleft_pnt.Y + self.botright_pnt.Y + 1) / 2) - 1,
					},
					self, NODE_TOPLEFT,
				)
			} else {
				self.TopLeft.Insert(new_obj)
			}
			self.is_leaf = false
			self.append_object(new_obj)
		} else { // bottom left
			if self.BottomLeft == nil {
				self.BottomLeft = NewQuadNode(new_obj.Id, true,
					self.Width/2, self.Height/2,
					Vector2{
						self.topleft_pnt.X,
						(self.topleft_pnt.Y + self.botright_pnt.Y + 1) / 2,
					},
					Vector2{
						((self.topleft_pnt.X + self.botright_pnt.X + 1) / 2) - 1,
						self.botright_pnt.Y,
					},
					self, NODE_BOTTOMLEFT,
				)
			} else {
				self.BottomLeft.Insert(new_obj)
			}
			self.is_leaf = false
			self.append_object(new_obj)
		}
	} else { // right
		if new_obj.Pos.Y <= (tlp.Y+brp.Y)/2 { // top right
			if self.TopRight == nil {
				self.TopRight = NewQuadNode(new_obj.Id, true,
					self.Width/2, self.Height/2,
					Vector2{
						((self.topleft_pnt.X + self.botright_pnt.X + 1) / 2),
						self.topleft_pnt.Y,
					}, Vector2{
						self.botright_pnt.X,
						((self.topleft_pnt.Y + self.botright_pnt.Y + 1) / 2) - 1, // 위쪽이기 때문에 마지막 -1
					},
					self, NODE_TOPRIGHT,
				)
			} else {
				self.TopRight.Insert(new_obj)
			}
			self.is_leaf = false
			self.append_object(new_obj)
		} else { // bottom right
			if self.BottomRight == nil {
				self.BottomRight = NewQuadNode(new_obj.Id, true,
					self.Width/2, self.Height/2,
					Vector2{
						(self.topleft_pnt.X + self.botright_pnt.X + 1) / 2,
						(self.topleft_pnt.Y + self.botright_pnt.Y + 1) / 2,
					},
					Vector2{
						self.botright_pnt.X,
						self.botright_pnt.Y,
					},
					self, NODE_BOTTOMRIGHT,
				)
			} else {
				self.BottomRight.Insert(new_obj)
			}
			self.is_leaf = false
			self.append_object(new_obj)
		}
	}
	return true
}

func (self *QuadNode) Nearest(target_pos Vector2, sight_radius int) []*GObject {
	var objects []*GObject
	targets_node := self.near(target_pos, sight_radius, &objects, self /*root*/)

	if sight_radius > 0 {
		if (target_pos.X - sight_radius) <= targets_node.topleft_pnt.X {
			// 왼
			left := Vector2{X: target_pos.X - sight_radius, Y: target_pos.Y}
			self.near(left, 0, &objects, self /*root*/)
		}
		if (target_pos.X + sight_radius) >= targets_node.botright_pnt.X {
			// 오른
			right := Vector2{X: target_pos.X + sight_radius, Y: target_pos.Y}
			self.near(right, 0, &objects, self /*root*/)
		}
		if (target_pos.Y - sight_radius) <= targets_node.topleft_pnt.Y {
			// 위
			top := Vector2{X: target_pos.X, Y: target_pos.Y - sight_radius}
			self.near(top, 0, &objects, self /*root*/)
		}
		if (target_pos.Y + sight_radius) >= targets_node.botright_pnt.Y {
			// 아래
			bot := Vector2{X: target_pos.X, Y: target_pos.Y + sight_radius}
			self.near(bot, 0, &objects, self /*root*/)
		}
	}
	return objects
}

func (self *QuadNode) near(target_pos Vector2, sight_radius int, out_objects *[]*GObject, root *QuadNode) *QuadNode {
	var inBoundary = func(p Vector2) bool {
		if self == nil {
			return false
		}
		return (p.X >= self.topleft_pnt.X && p.X <= self.botright_pnt.X &&
			p.Y >= self.topleft_pnt.Y && p.Y <= self.botright_pnt.Y)
	}
	if !inBoundary(target_pos) {
		return self
	}
	if self.IsLeaf() {
		*out_objects = append(*out_objects, self.Parent.GetAllObjects()...)
		return self
	}
	tlp := self.topleft_pnt
	brp := self.botright_pnt
	if (math.Abs(float64(tlp.X-brp.X)) <= LEAST_BLOCKSIZE) &&
		(math.Abs(float64(tlp.Y-brp.Y)) <= LEAST_BLOCKSIZE) {
		//TODO: check near sections object
		//TODO: 지금은 상하좌우만 확인중.. 대각선 확인 필요(원형), 그냥 대각선 좌표 검색하면 상하좌우에서 가져온 데이터와 겹침
		if sight_radius > 0 {
		}
		*out_objects = append(*out_objects, self.Parent.GetAllObjects()...)
		return self
	}

	x := target_pos.X
	y := target_pos.Y
	var near []*GObject
	if x < (tlp.X+brp.X)/2 { // left
		if y < (tlp.Y+brp.Y)/2 { // top left
			self.TopLeft.near(target_pos, sight_radius, out_objects, root)
		} else { // bottom left
			self.BottomLeft.near(target_pos, sight_radius, out_objects, root)
		}
	} else { // right
		if target_pos.Y <= (tlp.Y+brp.Y)/2 { // top right
			self.TopRight.near(target_pos, sight_radius, out_objects, root)
		} else { // bottom right
			self.BottomRight.near(target_pos, sight_radius, out_objects, root)
		}
	}
	if near != nil {
		*out_objects = append(*out_objects, near...)
	}
	return self
}

func (self *QuadNode) search(target_pos Vector2) **QuadNode {
	var inBoundary = func(p Vector2) bool {
		if self == nil {
			return false
		}
		return (p.X >= self.topleft_pnt.X && p.X <= self.botright_pnt.X &&
			p.Y >= self.topleft_pnt.Y && p.Y <= self.botright_pnt.Y)
	}
	if !inBoundary(target_pos) {
		return nil
	}
	tlp := self.topleft_pnt
	brp := self.botright_pnt
	if (math.Abs(float64(tlp.X-brp.X)) <= LEAST_BLOCKSIZE) &&
		(math.Abs(float64(tlp.Y-brp.Y)) <= LEAST_BLOCKSIZE) {
		// found
		return &self
	}

	x := target_pos.X
	y := target_pos.Y
	if x < (tlp.X+brp.X)/2 { // left
		if y < (tlp.Y+brp.Y)/2 { // top left
			return self.TopLeft.search(target_pos)
		} else { // bottom left
			return self.BottomLeft.search(target_pos)
		}
	} else { // right
		if target_pos.Y <= (tlp.Y+brp.Y)/2 { // top right
			return self.TopRight.search(target_pos)
		} else { // bottom right
			return self.BottomRight.search(target_pos)
		}
	}
	return nil
}
func (self *QuadNode) MoveObject(target *GObject, to Vector2) {
	s := self.search(target.Pos)
	if s == nil {
		return
	}
	for i, obj := range (*s).objs {
		if obj.Id == target.Id {
			*&obj.Pos = to
			self.Insert(obj)
			(*s).objs[i] = nil
		}
	}
	after := self.search(to)
	if after == nil {
		fmt.Println("s")
	}
}
func (self *QuadNode) Move(from_pos Vector2, from_id int, to Vector2) {
	s := self.search(from_pos)
	if s == nil {
		return
	}
	for i, obj := range (*s).objs {
		if obj.Id == from_id {
			*&obj.Pos = to
			self.Insert(obj)
			(*s).objs[i] = nil
		}
	}
	after := self.search(to)
	if after == nil {
		fmt.Println("s")
	}
}
func (self *QuadNode) Remove(target_obj **GObject) {
	target_obj = nil
}
