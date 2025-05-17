package game

type MapNode struct {
	X       int    `json:"x"`
	Y       int    `json:"y"`
	Strengh int    `json:"strengh"`
	OwnerId uint32 `json:"ownerId"`
}

type Map struct {
	Nodes [5][5]MapNode `json:"nodes"`
}

func NewMap() *Map {
	var nodes [5][5]MapNode
	for i := range 5 {
		for j := range 5 {
			nodes[i][j] = MapNode{
				X:       i,
				Y:       j,
				Strengh: 1,
				OwnerId: 0, // 0 means unowned
			}
		}
	}
	return &Map{Nodes: nodes}
}
