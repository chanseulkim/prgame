package gcore

import (
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
	NodeSector    int

	TopLeft     *QuadNode
	TopRight    *QuadNode
	BottomLeft  *QuadNode
	BottomRight *QuadNode
	Parent      *QuadNode
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

func ConstructQuadTree(grid [][]int) *QuadNode {
	var construct_task func(startr int, endr int, startc int, endc int, grid [][]int, prev *QuadNode, node_sector int) *QuadNode
	construct_task = func(startr int, endr int, startc int, endc int, grid [][]int, prev *QuadNode, node_sector int) *QuadNode {
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
				is_leaf:      true,
				Width:        endc - startc,
				Height:       endr - startr,
				topleft_pnt:  tlp,
				botright_pnt: brp,
				Parent:       prev,
				NodeSector:   node_sector,
			}
		}
		new_node := &QuadNode{
			Id:           val,
			is_leaf:      false,
			Width:        endc - startc,
			Height:       endr - startr,
			topleft_pnt:  tlp,
			botright_pnt: brp,
			Parent:       prev,
			NodeSector:   node_sector,
		}
		prev = new_node
		midr := startr + (endr-startr)/2
		midc := startc + (endc-startc)/2
		new_node.TopLeft = construct_task(startr, midr, startc, midc, grid, prev, NODE_TOPLEFT)
		new_node.TopRight = construct_task(startr, midr, midc, endc, grid, prev, NODE_TOPRIGHT)
		new_node.BottomLeft = construct_task(midr, endr, startc, midc, grid, prev, NODE_BOTTOMLEFT)
		new_node.BottomRight = construct_task(midr, endr, midc, endc, grid, prev, NODE_BOTTOMRIGHT)
		return new_node
	}
	return construct_task(0, len(grid), 0, len(grid[0]), grid, nil, NODE_ROOT)
}

func (self *QuadNode) append_object(new_obj *GObject) {
	self.objs = append(self.objs, new_obj)
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
					Width:      self.Width / 2,
					Height:     self.Height / 2,
					is_leaf:    true,
					Id:         new_obj.Id,
					NodeSector: NODE_TOPLEFT,
					Parent:     self,
				}
			}
			self.is_leaf = false
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
					Width:      self.Width / 2,
					Height:     self.Height / 2,
					is_leaf:    true,
					Id:         new_obj.Id,
					NodeSector: NODE_BOTTOMLEFT,
					Parent:     self,
				}
			}
			self.is_leaf = false
			self.BottomLeft.Insert(new_obj)
		}
	} else { // right
		if new_obj.Pos.Y <= (tlp.Y+brp.Y)/2 { // top right
			if self.TopRight == nil {
				self.TopRight = &QuadNode{
					topleft_pnt: Vector2{
						((self.topleft_pnt.X + self.botright_pnt.X + 1) / 2),
						self.topleft_pnt.X,
					},
					botright_pnt: Vector2{
						self.botright_pnt.X,
						((self.topleft_pnt.Y + self.botright_pnt.Y + 1) / 2) - 1, // 위쪽이기 때문에 마지막 -1
					},
					Width:      self.Width / 2,
					Height:     self.Height / 2,
					is_leaf:    true,
					Id:         new_obj.Id,
					NodeSector: NODE_TOPRIGHT,
					Parent:     self,
				}
			}
			self.is_leaf = false
			self.TopRight.Insert(new_obj)
		} else { // bottom right
			if self.BottomRight == nil {
				self.BottomRight = &QuadNode{
					topleft_pnt: Vector2{
						(self.topleft_pnt.X + self.botright_pnt.X + 1) / 2,
						(self.topleft_pnt.Y + self.botright_pnt.Y + 1) / 2,
					},
					botright_pnt: Vector2{
						self.botright_pnt.X,
						self.botright_pnt.Y,
					},
					Width:      self.Width / 2,
					Height:     self.Height / 2,
					is_leaf:    true,
					Id:         new_obj.Id,
					NodeSector: NODE_BOTTOMRIGHT,
					Parent:     self,
				}
			}
			self.is_leaf = false
			self.BottomRight.Insert(new_obj)
		}
	}
}

//func (self *QuadNode) SearchSector(target_obj *GObject) []*GObject {
func (self *QuadNode) NearObjs(obj *GObject) []*GObject {
	var objects []*GObject
	self.search_task(obj.Pos, obj.SightRadius, &objects, self)
	return objects
}
func (self *QuadNode) NearPosition(target_pos Vector2, sight_radius Float) []*GObject {
	var objects []*GObject
	self.search_task(target_pos, sight_radius, &objects, self /*root*/)
	return objects
}

func (self *QuadNode) search_task(target_pos Vector2, sight_radius Float, out_objects *[]*GObject, root *QuadNode) {
	var inBoundary = func(p Vector2) bool {
		if self == nil {
			return false
		}
		return (p.X >= self.topleft_pnt.X && p.X <= self.botright_pnt.X &&
			p.Y >= self.topleft_pnt.Y && p.Y <= self.botright_pnt.Y)
	}
	if !inBoundary(target_pos) {
		return
	}
	tlp := self.topleft_pnt
	brp := self.botright_pnt
	if (math.Abs(float64(tlp.X-brp.X)) <= LEAST_BLOCKSIZE) &&
		(math.Abs(float64(tlp.Y-brp.Y)) <= LEAST_BLOCKSIZE) {
		//TODO: check near sections object
		//TODO: 지금은 상하좌우만 확인중.. 대각선 확인 필요(원형), 그냥 대각선 좌표 검색하면 상하좌우에서 가져온 데이터와 겹침
		//TODO: 인접노드 즉시가져오기 필요. 지금은 새로운 위치를 다시검색
		if sight_radius > 0 {
			if (target_pos.Y - sight_radius) < self.topleft_pnt.Y {
				// 위 노드로 넘어감
				*out_objects = append(*out_objects, root.NearPosition(Vector2{target_pos.X, (target_pos.Y - sight_radius)}, 0)...)
			}
			if (target_pos.Y + sight_radius) > self.botright_pnt.Y {
				// 아래 노드로 넘어감
				*out_objects = append(*out_objects, root.NearPosition(Vector2{target_pos.X, (target_pos.Y + sight_radius)}, 0)...)
			}
			if (target_pos.X - sight_radius) < (self.topleft_pnt.X) {
				// 왼쪽 노드로 넘어감
				*out_objects = append(*out_objects, root.NearPosition(Vector2{(target_pos.X - sight_radius), target_pos.Y}, 0)...)
			}
			if (target_pos.X + sight_radius) > self.botright_pnt.X {
				// 오른쪽 노드로 넘어감
				*out_objects = append(*out_objects, root.NearPosition(Vector2{(target_pos.X + sight_radius), target_pos.Y}, 0)...)
			}
		}
		*out_objects = append(*out_objects, self.objs...)
	}

	x := target_pos.X
	y := target_pos.Y
	var near []*GObject
	if x < (tlp.X+brp.X)/2 { // left
		if y < (tlp.Y+brp.Y)/2 { // top left
			self.TopLeft.search_task(target_pos, sight_radius, out_objects, root)
		} else { // bottom left
			self.BottomLeft.search_task(target_pos, sight_radius, out_objects, root)
		}
	} else { // right
		if target_pos.Y <= (tlp.Y+brp.Y)/2 { // top right
			self.TopRight.search_task(target_pos, sight_radius, out_objects, root)
		} else { // bottom right
			self.BottomRight.search_task(target_pos, sight_radius, out_objects, root)
		}
	}
	if near != nil {
		*out_objects = append(*out_objects, near...)
	}
	return
}
func (self *QuadNode) Move(from Vector2, to Vector2) {

}
func (self *QuadNode) Remove(target_obj *GObject) {

}
