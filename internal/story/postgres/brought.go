package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/hadithopen-io/back/internal/story/types"
	"github.com/hadithopen-io/back/pkg/pgscan"
)

type Brought struct{ db *sqlx.DB }

func NewBrought(db *sqlx.DB) *Brought { return &Brought{db: db} }

func (b Brought) Create(ctx context.Context, brought types.Brought) (
	id int64,
	err error,
) {
	const query = `
insert into hadith.brought(brought)
values (:brought)
returning id
`

	return id, pgscan.Get(
		ctx,
		b.db,
		&id,
		query,
		brought,
	)
}
