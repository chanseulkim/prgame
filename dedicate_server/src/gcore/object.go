package gcore

type Player struct {
	Uid             string
	position        Vector3
	colision_radius Float
	index_inworld   Float
}

func (p *Player) Position() Vector3 {
	return p.position
}

type World struct {
	Players  map[string]Player
	Objects  []GObject
	position Vector3
}

var world_instance *World

func GetWorld() *World {
	if world_instance == nil {
		world_instance = &World{Players: make(map[string]Player)}
	}
	return world_instance
}

func (w *World) Position() Vector3 {
	return w.position
}
