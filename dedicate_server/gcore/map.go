package gcore

// // TODO: 좌표값은 Float이고 배열은 정수 index인데, 어떻게 배열로 처리해야할지..
// type GameMap struct {
// 	area Area
// 	// grid        [][]int
// 	object_tree *QuadNode
// }

// func NewGameMap(x int, y int) *GameMap {
// 	var new_space = make(Area, x)
// 	for i := range new_space {
// 		new_space[i] = make([]string, y)
// 	}
// 	return &GameMap{
// 		area: new_space,
// 	}
// }

// func (g *GameMap) Init() {
// 	g.object_tree = ConstructQuadTree(g.grid)
// }

// // TODO: collision shape를 등록해야함
// func (m *GameMap) AddObject(obj *GObject) {
// 	m.area[int(obj.Pos.X)][int(obj.Pos.Y)] = obj.Name
// }
// func (m *GameMap) GetArea() Area { return m.area }
