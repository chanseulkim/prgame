package libs

func Nearest() {

}

func NearestTopLeft(q *QuadNode, objects *[]*GObject) {
	parent := q.Parent
	if parent.NodeSector == NODE_BOTTOMLEFT {
		parent = parent.Parent.TopLeft
	} else if parent.NodeSector == NODE_BOTTOMRIGHT {
		parent = parent.Parent.TopRight
	} else { // 더 상위로..
		NearestTopLeft(parent.Parent, objects)
	}
	if parent != nil {
		if parent.BottomLeft != nil {
			*objects = append(*objects, parent.BottomLeft.objs...)
		}
	}
}
func NearestTopRight(q *QuadNode, objects *[]*GObject) {
	parent := q.Parent
	if parent.NodeSector == NODE_BOTTOMLEFT {
		parent = parent.Parent.TopLeft
	} else if parent.NodeSector == NODE_BOTTOMRIGHT {
		parent = parent.Parent.TopRight
	} else { // 더 상위로..
		NearestTopLeft(parent.Parent, objects)
	}
	if parent != nil {
		if parent.BottomLeft != nil {
			*objects = append(*objects, parent.BottomLeft.objs...)
		}
	}
}
