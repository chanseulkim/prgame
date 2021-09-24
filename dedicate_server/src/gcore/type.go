package gcore

type Float float32

type Vector2 struct {
	X, Y Float
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

//type Square struct {}
//type Circle struct{}
//type Triangle struct{}
//type EquilTriangle struct{}
