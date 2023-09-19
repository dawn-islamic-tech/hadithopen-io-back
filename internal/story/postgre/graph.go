package postgre

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/hadithopen-io/back/internal/story/types"
)

type Graph struct {
	db *pgx.Conn
}

func NewGraph(db *pgx.Conn) *Graph {
	return &Graph{
		db: db,
	}
}

func (g Graph) Nodes(context.Context, int64) (nodes []types.Node, err error) {
	return nil, err
}

func (g Graph) Edges(context.Context, int64) (edges types.Edges, err error) {
	return nil, err
}
