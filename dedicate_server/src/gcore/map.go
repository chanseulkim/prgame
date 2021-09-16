package gcore

const (
	NODE_TOPLEFT int = 0
	NODE_TOPRIGHT
	NODE_BOTTOMLEFT
	NODE_BOTTOMRIGHT
)

type QuadTreeNode struct {
	Id            int
	IsLeaf        bool
	Width, Height int
	Objects       []GObject
	TopLeft       *QuadTreeNode
	TopRight      *QuadTreeNode
	BottomLeft    *QuadTreeNode
	BottomRight   *QuadTreeNode
}

func ConstructQuadTree(grid [][]int) *QuadTreeNode {
	var construct_task func(startr int, endr int, startc int, endc int, grid [][]int) *QuadTreeNode
	construct_task = func(startr int, endr int, startc int, endc int, grid [][]int) *QuadTreeNode {
		val := grid[startr][startc]
		loopend := false
		var r, c int = 0, 0
		for r = startr; r < endr; r++ {
			for c = startc; c < endc; c++ {
				if grid[r][c] != val {
					loopend = true
					break
				}
			}
			if loopend {
				break
			}
		}
		if (r == endr) && (c == endc) {
			return &QuadTreeNode{Id: val, IsLeaf: true, Width: endc - startc, Height: endr - startr}
		}
		new_node := &QuadTreeNode{Id: val, IsLeaf: false, Width: endc - startc, Height: endr - startr}
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

func (self *QuadTreeNode) Search(r int, c int, rlen int, clen int, node *QuadTreeNode) {
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
	quadTree *QuadTreeNode
}

func (g *GameMap) Init() {
	g.quadTree = ConstructQuadTree(g.grid)
}

// TODO: collision shape를 등록해야함
func (m *GameMap) AddObject(obj *GObject) {
	m.area[int(obj.Position.X)][int(obj.Position.Y)] = obj.Name
}
func (m *GameMap) GetArea() Area { return m.area }
