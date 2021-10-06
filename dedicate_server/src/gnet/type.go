package gnet

type Float float32

type Vector2 struct {
	X, Y int
}

type Vector3 struct {
	X, Y, Z Float
}

type Rectangle struct {
	/*
		TopLeft ┌─────┐
			    │     │
			    └─────┘ BotRight
	*/
	TopLeft  Vector2
	BotRight Vector2
}

const (
	ID_OBJECT = 0
	ID_USER   = 1
)

//type Square struct {}
//type Circle struct{}
//type Triangle struct{}
//type EquilTriangle struct{}
