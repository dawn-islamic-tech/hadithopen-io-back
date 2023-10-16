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

func (v Version) Update(ctx context.Context, version types.Version) (
	err error,
) {
	const query = `
update hadith.versions
	set  
	    is_default = :is_default,
	    original = :original,
	    version = :version
where id = :id 
`

	return pgscan.Exec(
		ctx,
		v.db,
		query,
		version,
	)
}

func (v Version) Get(ctx context.Context, id int64) (version types.Version, err error) {
	const query = `
select 
    id,
	brought_id, 
	is_default, 
	original,
	version
from hadith.versions
where id = :id 
`

	return version, pgscan.Get(
		ctx,
		v.db,
		&version,
		query,
		map[string]any{
			"id": id,
		},
	)
}
