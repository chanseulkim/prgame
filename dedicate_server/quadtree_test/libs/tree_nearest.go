package libs

import (
	"container/list"
	"fmt"
	"math"
)

const LEAST_BLOCKSIZE = 20

type SectorType int

const (
	NODE_TOPLEFT SectorType = iota
	NODE_TOPRIGHT
	NODE_BOTTOMLEFT
	NODE_BOTTOMRIGHT
	NODE_ROOT
)

type QuadNode struct {
	Id            int
	is_leaf       bool
	Width, Height int
	topleft_pnt   Vector2
	botright_pnt  Vector2
	objs          *list.List
	NodeSector    SectorType

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
	node_sector SectorType) *QuadNode {
	return &QuadNode{
		Id:           id,
		is_leaf:      true,
		objs:         list.New(),
		Width:        width,
		Height:       height,
		topleft_pnt:  topleft_pnt,
		botright_pnt: botright_pnt,
		Parent:       parent,
		NodeSector:   node_sector,
	}
}
func (node *QuadNode) IsLeaf() bool {
	return (node.TopLeft == nil) &&
		(node.TopRight == nil) &&
		(node.BottomLeft == nil) &&
		(node.BottomRight == nil)
}

// 자식 노드들이 가진 모든 오브젝트들 반환
// TODO: 변화가 있는 오브젝트들만 추출
func (node *QuadNode) GetAllObjects() *list.List {
	if node == nil {
		return nil
	}
	var objlist *list.List = list.New()
	objlist.PushBackList(node.objs)
	if node.TopLeft != nil {
		if !node.TopLeft.IsLeaf() {
			objlist.PushBackList(node.TopLeft.GetAllObjects())
		} else {
			objlist.PushBackList(node.TopLeft.objs)
		}
	}
	if node.TopRight != nil {
		if !node.TopRight.IsLeaf() {
			objlist.PushBackList(node.TopRight.GetAllObjects())
		} else {
			objlist.PushBackList(node.TopRight.objs)
		}
	}
	if node.BottomLeft != nil {
		if !node.BottomLeft.IsLeaf() {
			objlist.PushBackList(node.BottomLeft.GetAllObjects())
		} else {
			objlist.PushBackList(node.BottomLeft.objs)
		}
	}
	if node.BottomRight != nil {
		if !node.BottomRight.IsLeaf() {
			objlist.PushBackList(node.BottomRight.GetAllObjects())
		} else {
			objlist.PushBackList(node.BottomRight.objs)
		}
	}
	return objlist
}

// TODO: 변화가 있는 오브젝트들만 추출
func (node *QuadNode) GetAllObjectsToCh(objs_ch chan *list.List) {
	objs_ch <- node.objs
	if node.TopLeft != nil {
		if !node.TopLeft.IsLeaf() {
			objs_ch <- node.TopLeft.GetAllObjects()
		} else {
			objs_ch <- node.TopLeft.objs
		}
	}
	if node.TopRight != nil {
		if !node.TopRight.IsLeaf() {
			objs_ch <- node.TopRight.GetAllObjects()
		} else {
			objs_ch <- node.TopRight.objs
		}
	}
	if node.BottomLeft != nil {
		if !node.BottomLeft.IsLeaf() {
			objs_ch <- node.BottomLeft.GetAllObjects()
		} else {
			objs_ch <- node.BottomLeft.objs
		}
	}
	if node.BottomRight != nil {
		if !node.BottomRight.IsLeaf() {
			objs_ch <- node.BottomRight.GetAllObjects()
		} else {
			objs_ch <- node.BottomRight.objs
		}
	}
	if node.NodeSector == NODE_ROOT {
		close(objs_ch)
	}
}

// Rectangle
func NewQuadTreeRoot(r Rectangle) *QuadNode {
	return NewQuadNode(0, true,
		r.BotRight.X, r.BotRight.Y,
		Vector2{X: r.TopLeft.X, Y: r.TopLeft.Y},
		Vector2{X: r.BotRight.X, Y: r.BotRight.Y},
		nil, NODE_ROOT,
	)
}
func ConstructQuadTree(grid [][]int) *QuadNode {
	var construct_task func(startr int, endr int, startc int, endc int, grid [][]int, parent *QuadNode, node_sector SectorType) *QuadNode
	construct_task = func(startr int, endr int, startc int, endc int, grid [][]int, parent *QuadNode, node_sector SectorType) *QuadNode {
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

func (node *QuadNode) append_object(new_obj *GObject) {
	node.objs.PushBack(new_obj)
}

func (node *QuadNode) Insert(new_obj *GObject) bool {
	if node == nil {
		return false
	}
	var inBoundary = func(p *Vector2) bool {
		return (p.X >= node.topleft_pnt.X && p.X <= node.botright_pnt.X+1 &&
			p.Y >= node.topleft_pnt.Y && p.Y <= node.botright_pnt.Y+1)
	}
	if !inBoundary(&new_obj.Pos) {
		fmt.Println("Error ", new_obj.Name, " is not in boundary")
		return false
	}
	tlp := node.topleft_pnt
	brp := node.botright_pnt
	if (math.Abs(float64(tlp.X-brp.X)) <= LEAST_BLOCKSIZE) &&
		math.Abs(float64(tlp.Y-brp.Y)) <= LEAST_BLOCKSIZE {
		node.append_object(new_obj)
		return true
	}
	if new_obj.Pos.X < (tlp.X+brp.X)/2 { // left
		if new_obj.Pos.Y < (tlp.Y+brp.Y)/2 { // top left
			if node.TopLeft == nil {
				node.TopLeft = NewQuadNode(new_obj.Id, true,
					node.Width/2, node.Height/2,
					Vector2{node.topleft_pnt.X, node.topleft_pnt.Y},
					Vector2{((node.topleft_pnt.X + node.botright_pnt.X + 1) / 2) - 1, // 왼쪽이기때문에 마지막 -1
						((node.topleft_pnt.Y + node.botright_pnt.Y + 1) / 2) - 1,
					},
					node, NODE_TOPLEFT,
				)
			}
			node.is_leaf = false
			node.TopLeft.Insert(new_obj)
		} else { // bottom left
			if node.BottomLeft == nil {
				node.BottomLeft = NewQuadNode(new_obj.Id, true,
					node.Width/2, node.Height/2,
					Vector2{
						node.topleft_pnt.X,
						(node.topleft_pnt.Y + node.botright_pnt.Y + 1) / 2,
					},
					Vector2{
						((node.topleft_pnt.X + node.botright_pnt.X + 1) / 2) - 1,
						node.botright_pnt.Y,
					},
					node, NODE_BOTTOMLEFT,
				)
			}
			node.is_leaf = false
			node.BottomLeft.Insert(new_obj)
		}
	} else { // right
		if new_obj.Pos.Y <= (tlp.Y+brp.Y)/2 { // top right
			if node.TopRight == nil {
				node.TopRight = NewQuadNode(new_obj.Id, true,
					node.Width/2, node.Height/2,
					Vector2{
						((node.topleft_pnt.X + node.botright_pnt.X + 1) / 2),
						node.topleft_pnt.Y,
					}, Vector2{
						node.botright_pnt.X,
						((node.topleft_pnt.Y + node.botright_pnt.Y + 1) / 2) - 1, // 위쪽이기 때문에 마지막 -1
					},
					node, NODE_TOPRIGHT,
				)
			}
			node.is_leaf = false
			node.TopRight.Insert(new_obj)
		} else { // bottom right
			if node.BottomRight == nil {
				node.BottomRight = NewQuadNode(new_obj.Id, true,
					node.Width/2, node.Height/2,
					Vector2{
						(node.topleft_pnt.X + node.botright_pnt.X + 1) / 2,
						(node.topleft_pnt.Y + node.botright_pnt.Y + 1) / 2,
					},
					Vector2{
						node.botright_pnt.X,
						node.botright_pnt.Y,
					},
					node, NODE_BOTTOMRIGHT,
				)
			}
			node.is_leaf = false
			node.BottomRight.Insert(new_obj)
		}
	}
	return true
}
func (node *QuadNode) Nearest(target_pos Vector2, sight_radius int) *list.List {
	var is_gotten map[int]bool = make(map[int]bool)
	return node.nearest_task(target_pos, sight_radius, node, &is_gotten)
}
func (node *QuadNode) nearest_task(target_pos Vector2, sight_radius int, root *QuadNode, is_gotten *map[int]bool) *list.List {
	if node == nil {
		fmt.Println("error Nearest2")
		return nil
	}
	tlp := node.topleft_pnt
	brp := node.botright_pnt
	CheckOver := func() *list.List {
		var nearests *list.List = list.New()
		if sight_radius > 0 {
			// 상하좌우
			if (target_pos.Y - sight_radius) < tlp.Y { // up
				if (node.NodeSector == NODE_TOPLEFT) || (node.NodeSector == NODE_TOPRIGHT) {
					near_pos := Vector2{X: target_pos.X, Y: (target_pos.Y - sight_radius)}
					nearests.PushBackList(root.nearest_task(near_pos, 0, root, is_gotten))
				}
			}
			if (target_pos.Y + sight_radius) >= brp.Y { // down
				if (node.NodeSector == NODE_BOTTOMLEFT) || (node.NodeSector == NODE_BOTTOMRIGHT) {
					near_pos := Vector2{X: target_pos.X, Y: (target_pos.Y + sight_radius)}
					nearests.PushBackList(root.nearest_task(near_pos, 0, root, is_gotten))
				}
			}
			if (target_pos.X - sight_radius) < tlp.X { // left
				if (node.NodeSector == NODE_TOPLEFT) || (node.NodeSector == NODE_BOTTOMLEFT) {
					near_pos := Vector2{X: target_pos.X - sight_radius, Y: target_pos.Y}
					nearests.PushBackList(root.nearest_task(near_pos, 0, root, is_gotten))
				}
			}
			if (target_pos.X + sight_radius) < brp.X { // right
				if (node.NodeSector == NODE_TOPRIGHT) || (node.NodeSector == NODE_BOTTOMRIGHT) {
					near_pos := Vector2{X: target_pos.X + sight_radius, Y: target_pos.Y}
					nearests.PushBackList(root.nearest_task(near_pos, 0, root, is_gotten))
				}
			}
			//대각

			// top left
			tl_outbnd := Vector2{
				X: (target_pos.X - (sight_radius / 2)),
				Y: (target_pos.Y - (sight_radius / 2)),
			}
			if tl_outbnd.Y < tlp.Y {
				if tl_outbnd.X < tlp.X {
					if node.NodeSector == NODE_TOPLEFT {
						nearests.PushBackList(root.nearest_task(tl_outbnd, 0, root, is_gotten))
					}
				}
			}
			// top right
			tr_outbnd := Vector2{
				X: (target_pos.X + (sight_radius / 2)),
				Y: (target_pos.Y + (sight_radius / 2)),
			}
			if tr_outbnd.Y >= tlp.Y {
				if tr_outbnd.X >= tlp.X {
					if node.NodeSector == NODE_TOPRIGHT {
						nearests.PushBackList(root.nearest_task(tr_outbnd, 0, root, is_gotten))
					}
				}
			}
			// bottom left
			bl_outbnd := Vector2{
				X: (target_pos.X - (sight_radius / 2)),
				Y: (target_pos.Y + (sight_radius / 2)),
			}
			if bl_outbnd.Y >= brp.Y {
				if bl_outbnd.X < brp.X {
					if node.NodeSector == NODE_BOTTOMLEFT {
						nearests.PushBackList(root.nearest_task(bl_outbnd, 0, root, is_gotten))
					}
				}
			}

			br_outbnd := Vector2{
				X: (target_pos.X + (sight_radius / 2)),
				Y: (target_pos.Y + (sight_radius / 2)),
			}
			if br_outbnd.Y < brp.Y { // right
				if br_outbnd.X < brp.X {
					if node.NodeSector == NODE_BOTTOMRIGHT {
						nearests.PushBackList(root.nearest_task(br_outbnd, 0, root, is_gotten))
					}
				}
			}
		}
		nearests.PushBackList(node.GetAllObjects())
		return nearests
	}

	if (math.Abs(float64(tlp.X-brp.X)) <= LEAST_BLOCKSIZE) &&
		(math.Abs(float64(tlp.Y-brp.Y)) <= LEAST_BLOCKSIZE) {
		fmt.Println("found")
		// found
		return CheckOver()
	}

	x := target_pos.X
	y := target_pos.Y
	if x < (tlp.X+brp.X)/2 { // left
		if y < (tlp.Y+brp.Y)/2 { // top left
			if node.TopLeft == nil {
				return CheckOver()
			}
			return node.TopLeft.nearest_task(target_pos, sight_radius, root, is_gotten)
		} else { // bottom left
			if node.BottomLeft == nil {
				return CheckOver()
			}
			return node.BottomLeft.nearest_task(target_pos, sight_radius, root, is_gotten)
		}
	} else { // right
		if target_pos.Y <= (tlp.Y+brp.Y)/2 { // top right
			if node.TopRight == nil {
				return CheckOver()
			}
			return node.TopRight.nearest_task(target_pos, sight_radius, root, is_gotten)
		} else { // bottom right
			if node.BottomRight == nil {
				return CheckOver()
			}
			return node.BottomRight.nearest_task(target_pos, sight_radius, root, is_gotten)
		}
	}
	return nil
}

func (node *QuadNode) near(target_pos Vector2, sight_radius int, out_objects *list.List, root *QuadNode) *QuadNode {
	var inBoundary = func(p Vector2) bool {
		if node == nil {
			return false
		}
		return (p.X >= node.topleft_pnt.X && p.X <= node.botright_pnt.X &&
			p.Y >= node.topleft_pnt.Y && p.Y <= node.botright_pnt.Y)
	}
	if !inBoundary(target_pos) {
		return node
	}
	if node.IsLeaf() {
		out_objects.PushBackList(node.Parent.GetAllObjects())
		return node
	}
	tlp := node.topleft_pnt
	brp := node.botright_pnt
	if (math.Abs(float64(tlp.X-brp.X)) <= LEAST_BLOCKSIZE) &&
		(math.Abs(float64(tlp.Y-brp.Y)) <= LEAST_BLOCKSIZE) {
		//TODO: check near sections object
		//TODO: 지금은 상하좌우만 확인중.. 대각선 확인 필요(원형), 그냥 대각선 좌표 검색하면 상하좌우에서 가져온 데이터와 겹침
		if sight_radius > 0 {
			//
		}
		out_objects.PushBackList(node.Parent.GetAllObjects())
		return node
	}

	x := target_pos.X
	y := target_pos.Y
	var near_node *QuadNode
	if x < (tlp.X+brp.X)/2 { // left
		if y < (tlp.Y+brp.Y)/2 { // top left
			near_node = node.TopLeft.near(target_pos, sight_radius, out_objects, root)
		} else { // bottom left
			near_node = node.BottomLeft.near(target_pos, sight_radius, out_objects, root)
		}
	} else { // right
		if target_pos.Y <= (tlp.Y+brp.Y)/2 { // top right
			near_node = node.TopRight.near(target_pos, sight_radius, out_objects, root)
		} else { // bottom right
			near_node = node.BottomRight.near(target_pos, sight_radius, out_objects, root)
		}
	}
	if near_node != nil {
		out_objects.PushBackList(near_node.GetAllObjects())
	}
	return node
}

func (node *QuadNode) search(target_pos Vector2) **QuadNode {
	var inBoundary = func(p Vector2) bool {
		if node == nil {
			return false
		}
		return (p.X >= node.topleft_pnt.X && p.X <= node.botright_pnt.X &&
			p.Y >= node.topleft_pnt.Y && p.Y <= node.botright_pnt.Y)
	}
	if !inBoundary(target_pos) {
		return nil
	}
	tlp := node.topleft_pnt
	brp := node.botright_pnt
	if (math.Abs(float64(tlp.X-brp.X)) <= LEAST_BLOCKSIZE) &&
		(math.Abs(float64(tlp.Y-brp.Y)) <= LEAST_BLOCKSIZE) {
		// found
		return &node
	}

	x := target_pos.X
	y := target_pos.Y
	if x < (tlp.X+brp.X)/2 { // left
		if y < (tlp.Y+brp.Y)/2 { // top left
			return node.TopLeft.search(target_pos)
		} else { // bottom left
			return node.BottomLeft.search(target_pos)
		}
	} else { // right
		if target_pos.Y <= (tlp.Y+brp.Y)/2 { // top right
			return node.TopRight.search(target_pos)
		} else { // bottom right
			return node.BottomRight.search(target_pos)
		}
	}
	return nil
}
func (node *QuadNode) Move(from_pos Vector2, from_id int, to Vector2) {
	found_node := node.search(from_pos)
	if found_node == nil {
		return
	}
	var obj_list *list.List = (*found_node).objs
	for e := obj_list.Front(); e != nil; e = e.Next() {
		var obj *GObject = e.Value.(*GObject)
		if obj.Id == from_id {
			obj_list.Remove(e)
			obj.Pos = to
			node.Insert(obj)
		}
	}
}
func (node *QuadNode) Remove(target_obj **GObject) {
	target_obj = nil
}
