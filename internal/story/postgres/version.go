package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/hadithopen-io/back/internal/story/types"
	"github.com/hadithopen-io/back/pkg/pgscan"
)

type Version struct{ db *sqlx.DB }

func NewVersion(db *sqlx.DB) *Version { return &Version{db: db} }

func (v Version) Create(ctx context.Context, version types.Version) (
	id int64,
	err error,
) {
	const query = `
insert into hadith.versions(brought_id, is_default, original, version) 
values (:brought_id, :is_default, :original, :version)
returning id
`

	return id, pgscan.Get(
		ctx,
		v.db,
		&id,
		query,
		version,
	)
}
