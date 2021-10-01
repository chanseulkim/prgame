package main

import (
	"fmt"
	. "pr/libs"
)

var grid [][]int

func init_grid() {
	grid = make([][]int, 8)
	for i := 0; i < len(grid); i++ {
		grid[i] = make([]int, 8)
	}
	grid[2][6] = 1
	grid[2][7] = 1
	grid[3][6] = 1
	grid[3][7] = 1

	grid[4][2] = 1
	grid[4][3] = 1
	grid[5][2] = 1
	grid[5][3] = 1
}

var quad_tree *QuadNode

var sight_radius Float = 1

func qt_test() {

}
func test_insert() {
	topleft_topleft := NewGObject(1, "user", Vector2{X: 1, Y: 1}, sight_radius)
	quad_tree.Insert(topleft_topleft)

	topleft_topright := NewGObject(2, "user", Vector2{X: 2, Y: 1}, sight_radius)
	quad_tree.Insert(topleft_topright)

	topleft_botleft := NewGObject(3, "user", Vector2{X: 1, Y: 2}, sight_radius)
	quad_tree.Insert(topleft_botleft)

	topleft_botright := NewGObject(4, "user", Vector2{X: 2, Y: 2}, sight_radius)
	quad_tree.Insert(topleft_botright)
}
func test_search() {
	other := NewGObject(1, "other_usr", Vector2{X: 4, Y: 3}, sight_radius)
	quad_tree.Insert(other)
	other2 := NewGObject(2, "other_usr2", Vector2{X: 3, Y: 4}, sight_radius)
	quad_tree.Insert(other2)

	tln := NewGObject(3, "tln", Vector2{X: 3, Y: 3}, sight_radius)
	quad_tree.Insert(tln)

	new_obj := NewGObject(4, "user", Vector2{X: 4, Y: 4}, sight_radius)
	quad_tree.Insert(new_obj)

	objall := quad_tree.GetAllObjects()
	fmt.Println("GetAllObjects : ", len(objall))

	near := quad_tree.Nearest(new_obj.Pos, new_obj.SightRadius)
	fmt.Println(near)
}
func test_move() {
	usr := NewGObject(1, "usr1", Vector2{X: 4, Y: 4}, sight_radius)
	quad_tree.Insert(usr)

	quad_tree.Move(usr.Pos, usr.Id, Vector2{X: 3, Y: 4})
}

func main() {
	qt_test()
	return

	init_grid()
	quad_tree = ConstructQuadTree(grid)

	test_move()
	test_search()
	// fmt.Println(quad_tree)

	//test_insert()

	// new_obj := NewGObject(1, "user", Vector2{X: 5, Y: 2}, 1)
	// quad_tree.Insert(new_obj)

}
