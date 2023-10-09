package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/hadithopen-io/back/internal/story/types"
)

type Graph struct {
	db *sqlx.DB
}

func NewGraph(db *sqlx.DB) *Graph { return &Graph{db: db} }

func (g Graph) Nodes(context.Context, int64) (nodes []types.Node, err error) {
	return nil, err
}

func (g Graph) Edges(context.Context, int64) (edges types.Edges, err error) {
	return nil, err
}
