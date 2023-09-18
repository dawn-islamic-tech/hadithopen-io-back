package types

type Node struct {
	ID   int64
	Name string
}

type Edges map[int64][]int64

type Graph struct {
	Nodes []Node
	Edges Edges
}
