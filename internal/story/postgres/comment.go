package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/hadithopen-io/back/internal/story/types"
	"github.com/hadithopen-io/back/pkg/pgscan"
)

type Comment struct{ db *sqlx.DB }

func NewComment(db *sqlx.DB) *Comment { return &Comment{db: db} }

func (c Comment) Create(ctx context.Context, comment types.Comment) (
	id int64,
	err error,
) {
	const query = `
insert into hadith.comments(story_id, brought_id, comment) 
values (:story_id, :brought_id, :comment)
returning id
`

	return id, pgscan.Get(
		ctx,
		c.db,
		&id,
		query,
		comment,
	)
}

func (c Comment) Update(ctx context.Context, comment types.Comment) (
	err error,
) {
	const query = `
update hadith.comments
	set comment = :comment
where id = :id
`

	return pgscan.Exec(
		ctx,
		c.db,
		query,
		comment,
	)
}

func (c Comment) Get(ctx context.Context, id int64) (
	comment types.Comment,
	err error,
) {
	const query = `
select 
    id, 
    brought_id,
    comment
from hadith.comments
where id = :id
`

	return comment, pgscan.Get(
		ctx,
		c.db,
		&comment,
		query,
		map[string]any{
			"id": id,
		},
	)
}
