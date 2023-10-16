package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/hadithopen-io/back/internal/story/types"
	"github.com/hadithopen-io/back/pkg/errors"
	"github.com/hadithopen-io/back/pkg/pgscan"
	"github.com/hadithopen-io/back/pkg/usercontext"
)

type Hadith struct{ db *sqlx.DB }

func NewHadith(db *sqlx.DB) *Hadith {
	return &Hadith{
		db: db,
	}
}

func (h Hadith) Few(ctx context.Context) (compacts []types.HadithCompact, _ error) {
	const query = `
select 
	s.id 			as id,
	t.translate  	as title,
    vt.translate 	as description
from hadith.stories s
	 inner join hadith.translates t
				on t.id = cast(
				    (s.title ->> :lang) as bigint
				)
	 inner join hadith.map_story_versions msv
				on msv.story_id = s.id
	 inner join hadith.versions v
				on v.id = msv.version_id 
					   and v.is_default
	 left join hadith.translates vt
                   on vt.id = cast(
				    (v.version ->> :lang) as bigint
				)
`

	u, err := usercontext.Get(ctx)
	if err != nil {
		return nil, err
	}

	return compacts, errors.Wrap(
		pgscan.Select(
			ctx,
			h.db,
			&compacts,
			query,
			map[string]any{
				"lang": u.Lang,
			},
		),
		"getting compacts",
	)
}

func (h Hadith) Get(ctx context.Context, storyID int64) (hadith types.Hadith, _ error) {
	const query = `
select
	s.id         as story_id,
	t.translate  as story_title,
	t.id 		 as story_translate_id,
	t.lang 		 as story_lang,
	
	v.id 			as version_id,
	v.is_default 	as version_is_default,
	v.original      as version_original,
	vt.translate 	as version,
	vt.id 			as version_translate_id,
	vt.lang 		as version_lang,
	v.brought_id 	as version_brought_id,

	
	c.id 			as comment_id,
	ct.translate 	as comment,
	ct.id 			as comment_translate_id,
	ct.lang 		as comment_lang,
	c.brought_id 	as comment_brought_id
from hadith.stories s
         inner join hadith.translates t
                    on t.id = cast ((s.title ->> :lang) as bigint)
         inner join hadith.map_story_versions msv
                    on msv.story_id = s.id
         inner join hadith.versions v
                    on v.id = msv.version_id 
                           and v.is_default
         left join hadith.translates vt
                   on vt.id = cast((v.version ->> :lang) as bigint)
         left join hadith.comments c
                   on c.story_id = s.id
         left join hadith.translates ct
                   on ct.id = cast((c.comment ->> :lang) as bigint)
         left join hadith.brought b
                   on c.brought_id = b.id
         left join hadith.translates bt
                   on bt.id = cast((b.brought ->> :lang) as bigint)
where s.id = :story_id
`

	u, err := usercontext.Get(ctx)
	if err != nil {
		return types.Hadith{}, err
	}

	return hadith, errors.Wrapf(
		pgscan.Get(
			ctx,
			h.db,
			&hadith,
			query,
			map[string]any{
				"lang":     u.Lang,
				"story_id": storyID,
			},
		),
		"getting hadith by story id - %d",
		storyID,
	)
}

func (h Hadith) Create(ctx context.Context, story types.Story) (
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
		h.db,
		&id,
		query,
		story,
	)
}

func (h Hadith) CreateMapVersion(ctx context.Context, storyID int64, versionID ...int64) (
	err error,
) {
	const query = `
insert into hadith.map_story_versions(story_id, version_id) 
values (:story_id, :version_id)
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
		h.db,
		query,
		mp,
	)
}

func (h Hadith) Update(ctx context.Context, story types.Story) (
	err error,
) {
	const query = `
update hadith.stories
	set title = :title
where id = :id
`

	return pgscan.Exec(
		ctx,
		h.db,
		query,
		story,
	)
}
