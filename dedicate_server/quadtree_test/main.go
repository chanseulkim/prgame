package main

import (
	"fmt"
	. "quadtree/quadtree"
	"strconv"
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

var sight_radius int = 20
var quad_tree *QuadNode = NewQuadTreeRoot(Rectangle{
	TopLeft:  Vector2{X: 0, Y: 0},
	BotRight: Vector2{X: 1024, Y: 600},
})

func test_insert() {
	tree := NewQuadTreeRoot(Rectangle{
		TopLeft:  Vector2{X: 0, Y: 0},
		BotRight: Vector2{X: 1024, Y: 600},
	})
	topleft_topleft := NewGObject(1, "user1", Vector2{X: 1, Y: 1}, sight_radius)
	tree.Insert(topleft_topleft)

	topleft_topright := NewGObject(2, "user2", Vector2{X: 2, Y: 1}, sight_radius)
	tree.Insert(topleft_topright)

	topleft_botleft := NewGObject(3, "user3", Vector2{X: 1, Y: 2}, sight_radius)
	tree.Insert(topleft_botleft)

	topleft_botright := NewGObject(4, "user4", Vector2{X: 2, Y: 2}, sight_radius)
	tree.Insert(topleft_botright)
}
func test_search() {
	var tree *QuadNode = NewQuadTreeRoot(Rectangle{
		TopLeft:  Vector2{X: 0, Y: 0},
		BotRight: Vector2{X: 32, Y: 32},
	})
	other := NewGObject(1, "other_usr", Vector2{X: 15, Y: 15}, sight_radius)
	tree.Insert(other)
	other2 := NewGObject(2, "other_usr2", Vector2{X: 15, Y: 16}, sight_radius)
	tree.Insert(other2)

	tln := NewGObject(3, "tln", Vector2{X: 16, Y: 15}, sight_radius)
	tree.Insert(tln)

	new_obj := NewGObject(4, "user", Vector2{X: 18, Y: 18}, sight_radius)
	tree.Insert(new_obj)

	objall := tree.GetAllObjects()
	if objall != nil {
		fmt.Println("GetAllObjects : ", objall.Len())
	}
	near := tree.Nearest(new_obj.Pos, new_obj.SightRadius)
	fmt.Println(near)
}
func test_full_search() {
	var tree *QuadNode = NewQuadTreeRoot(Rectangle{
		TopLeft:  Vector2{X: 0, Y: 0},
		BotRight: Vector2{X: 32, Y: 32},
	})
	for x := 0; x < 32; x++ {
		for y := 0; y < 32; y++ {
			newnode := NewGObject(x+y, "usr_"+strconv.Itoa(x+y), Vector2{X: x, Y: y}, sight_radius)
			tree.Insert(newnode)
		}
	}
	objall := tree.GetAllObjects()
	if objall != nil {
		fmt.Println("GetAllObjects : ", objall.Len())
	}
	center := Vector2{X: 24, Y: 16}
	nears := tree.Nearest(center, 20)
	cnt := 0
	prev := 0
	for i := nears.Front(); i != nil; i = i.Next() {
		list := i.Value.(*GObject)
		if prev != list.Pos.X {
			fmt.Println()
			prev = list.Pos.X
		}
		fmt.Print(list.Pos, "    ")
		if list.Pos.X < 10 {
			fmt.Print(" ")
		}
		if list.Pos.Y < 10 {
			fmt.Print(" ")
		}
		cnt++
		if cnt > 15 {
			fmt.Println()
			cnt = 0
		}
	}
	fmt.Println()
}

func test_move() {
	var quad_tree *QuadNode
	usr := NewGObject(1, "usr1", Vector2{X: 4, Y: 4}, sight_radius)
	quad_tree.Insert(usr)

	quad_tree.Move(usr.Pos, usr.Id, Vector2{X: 3, Y: 4})
}

func main() {
	// test_insert()
	test_full_search()
	return

	init_grid()
	quad_tree = ConstructQuadTree(grid)

	test_move()
	// fmt.Println(quad_tree)

	//test_insert()

	// new_obj := NewGObject(1, "user", Vector2{X: 5, Y: 2}, 1)
	// quad_tree.Insert(new_obj)

}
