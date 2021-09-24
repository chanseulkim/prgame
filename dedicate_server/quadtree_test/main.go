package main

import (
	"fmt"
	. "pr/libs"
)

func main() {
	// var grid [][]int
	grid := make([][]int, 8)
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

	quad_tree := ConstructQuadTree(grid)
	fmt.Println(quad_tree)
	// new_obj := NewGObject(1, "user", Vector2{X: 5, Y: 2}, 1)
	// quad_tree.Insert(new_obj)

	topleft_topleft := NewGObject(1, "user", Vector2{X: 1, Y: 1}, 1)
	quad_tree.Insert(topleft_topleft)

	topleft_topright := NewGObject(1, "user", Vector2{X: 2, Y: 1}, 1)
	quad_tree.Insert(topleft_topright)

	topleft_botleft := NewGObject(1, "user", Vector2{X: 1, Y: 2}, 1)
	quad_tree.Insert(topleft_botleft)

	topleft_botright := NewGObject(1, "user", Vector2{X: 2, Y: 2}, 1)
	quad_tree.Insert(topleft_botright)

}
