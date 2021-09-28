package libs

//2D
type NodeMatrix [][]*QuadNode

func NewNodeMatrix(x int, y int) NodeMatrix {
	grid := make([][]*QuadNode, y)
	for i := 0; i < y; i++ {
		grid[i] = make([]*QuadNode, x)
	}
	return grid
}
