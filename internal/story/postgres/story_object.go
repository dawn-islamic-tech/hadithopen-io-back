package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/hadithopen-io/back/internal/story/types"
	"github.com/hadithopen-io/back/pkg/pgscan"
)

type StoryObject struct{ db *sqlx.DB }

func NewStoryObject(db *sqlx.DB) *StoryObject { return &StoryObject{db: db} }

func (s StoryObject) Get(ctx context.Context, id int64) (story types.Story, err error) {
	const query = `
select 
    id, 
    title
from hadith.stories
where id = :id
`

	return story, pgscan.Get(
		ctx,
		s.db,
		&story,
		query,
		map[string]any{
			"id": id,
		},
	)
}

func (s StoryObject) Create(ctx context.Context, story types.Story) (
	id int64,
	err error,
) {
	const query = `
insert into hadith.stories(title)
values (:title)
returning id
`

	return id, pgscan.Get(
		ctx,
		s.db,
		&id,
		query,
		story,
	)
}

func (s StoryObject) CreateMapVersion(ctx context.Context, storyID int64, versionID ...int64) (
	err error,
) {
	if len(versionID) == 0 {
		return nil
	}

	const query = `
insert into hadith.map_story_versions(story_id, version_id) 
values (:story_id, :version_id)
on conflict do nothing
`

	mp := make([]map[string]any, 0, len(versionID))
	for _, id := range versionID {
		mp = append(
			mp,
			map[string]any{
				"story_id":   storyID,
				"version_id": id,
			},
		)
	}

	return pgscan.Exec(
		ctx,
		s.db,
		query,
		mp,
	)
}

func (s StoryObject) Update(ctx context.Context, story types.Story) (
	err error,
) {
	const query = `
update hadith.stories
	set title = :title
where id = :id
`

	return pgscan.Exec(
		ctx,
		s.db,
		query,
		story,
	)
}

func (s StoryObject) MarkDelete(ctx context.Context, id int64) error {
	const query = `
update hadith.stories
	set mark_delete = true
where id = :id
`

	return pgscan.Exec(
		ctx,
		s.db,
		query,
		map[string]any{
			"id": id,
		},
	)
}
