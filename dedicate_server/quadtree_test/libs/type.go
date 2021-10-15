package libs

type Float float32

type Vector2 struct {
	X, Y int
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
	SightRadius   int
	CollisionArea Rectangle
}

func NewGObject(id int, name string, pos Vector2, radius int) *GObject {
	return &GObject{
		Id:          id,
		Name:        name,
		Pos:         Vector2{pos.X, pos.Y},
		SightRadius: radius,
		CollisionArea: Rectangle{
			TopLeft:  Vector2{X: pos.X - radius, Y: pos.Y - radius},
			BotRight: Vector2{X: pos.X + radius, Y: pos.Y + radius},
		},
	}
}
