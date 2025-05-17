package game

type MapNode struct {
	X       int
	Y       int
	Strengh int
	OwnerId uint32
}

type Map struct {
	Nodes [5][5]MapNode
}

func NewMap() *Map {
	var nodes [5][5]MapNode
	for i := range 5 {
		for j := range 5 {
			nodes[i][j] = MapNode{
				X:       i,
				Y:       j,
				Strengh: 1,
				OwnerId: 1,
			}
		}
	}

	return &Map{
		Nodes: nodes,
	}
}
