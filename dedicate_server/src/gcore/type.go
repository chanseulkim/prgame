package gcore

type Float float32

type GObject struct {
	Name     string
	Position Vector2
}

type Vector2 struct {
	X, Y Float
}
type Vector3 struct {
	X, Y, Z Float
}

type Rectangle struct {
	StartX, StartY Float
	EndX, EndY     Float
}

//type Square struct {}
//type Circle struct{}
//type Triangle struct{}
//type EquilTriangle struct{}
