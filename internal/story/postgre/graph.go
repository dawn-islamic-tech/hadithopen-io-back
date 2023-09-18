package postgre

import (
	"context"

	"github.com/hadithopen-io/back/internal/story/types"
)

type Graph struct {
}

func NewGraph() *Graph {
	return &Graph{}
}

func (g Graph) Nodes(ctx context.Context, hadithID int64) (nodes []types.Node, err error) {
	//TODO implement me
	panic("implement me")
}

func (g Graph) Edges(ctx context.Context, hadithID int64) (edges types.Edges, err error) {
	//TODO implement me
	panic("implement me")
}
