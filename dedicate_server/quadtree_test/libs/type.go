package libs

type Float float32

type Vector2 struct {
	X, Y Float
}

type Vector3 struct {
	X, Y, Z Float
}

/*
	TopLeft ┌─────┐
		    │     │
		    └─────┘ BotRight
*/
type Rectangle struct {
	TopLeft  Vector2
	BotRight Vector2
}

type GObject struct {
	Id            int
	Name          string
	Pos           Vector2
	CollisionArea Rectangle
}

func NewGObject(id int, name string, pos Vector2, radius Float) *GObject {
	return &GObject{
		Pos: Vector2{pos.X, pos.Y},
		CollisionArea: Rectangle{
			TopLeft:  Vector2{X: pos.X - radius, Y: pos.Y - radius},
			BotRight: Vector2{X: pos.X + radius, Y: pos.Y + radius},
		},
	}
}
