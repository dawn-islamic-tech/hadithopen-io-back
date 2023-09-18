package story

import (
	"context"

	"github.com/hadithopen-io/back/internal/story/types"
	"github.com/hadithopen-io/back/pkg/errors"
)

type GraphStore interface {
	Nodes(ctx context.Context, hadithID int64) (
		nodes []types.Node,
		err error,
	)

	Edges(ctx context.Context, hadithID int64) (
		edges types.Edges,
		err error,
	)
}

type Transmitters struct {
	graph struct {
		store GraphStore
	}
}

func NewTransmitters(store GraphStore) *Transmitters {
	return &Transmitters{
		graph: struct {
			store GraphStore
		}{
			store: store,
		},
	}
}

func (t Transmitters) Get(ctx context.Context, hadithID int64) (types.Graph, error) {
	nodes, err := t.graph.store.Nodes(ctx, hadithID)
	if err != nil {
		return types.Graph{}, errors.Wrap(err, "getting nodes")
	}

	edges, err := t.graph.store.Edges(ctx, hadithID)
	if err != nil {
		return types.Graph{}, errors.Wrap(err, "getting edges")
	}

	return types.Graph{
		Nodes: nodes,
		Edges: edges,
	}, nil
}
