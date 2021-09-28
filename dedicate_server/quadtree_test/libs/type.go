package libs

type Float float32

type Vector2 struct {
	X, Y Float
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
	SightRadius   Float
	CollisionArea Rectangle
}

func NewGObject(id int, name string, pos Vector2, radius Float) *GObject {
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
