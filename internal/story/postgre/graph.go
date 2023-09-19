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

func (g Graph) Nodes(ctx context.Context, hadithID int64) (nodes []types.Node, err error) {
	//TODO implement me
	panic("implement me")
}

func (g Graph) Edges(ctx context.Context, hadithID int64) (edges types.Edges, err error) {
	//TODO implement me
	panic("implement me")
}
