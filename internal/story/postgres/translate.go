package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/hadithopen-io/back/internal/story/types"
	"github.com/hadithopen-io/back/pkg/errors"
	"github.com/hadithopen-io/back/pkg/pgscan"
)

type Translate struct{ db *sqlx.DB }

func NewTranslate(db *sqlx.DB) *Translate { return &Translate{db: db} }

func (t Translate) Create(ctx context.Context, translates []types.Translate) (
	created []types.Translate,
	err error,
) {
	if len(translates) == 0 {
		return nil, nil
	}

	const query = `
insert into hadith.translates(lang, translate)
values (:lang, :translate)
returning id, lang
	`

	return created, pgscan.Select(
		ctx,
		t.db,
		&created,
		query,
		translates,
	)
}

func (t Translate) Update(ctx context.Context, translates []types.Translate) (
	updated []types.Translate,
	err error,
) {
	if len(translates) == 0 {
		return nil, nil
	}

	const query = `
update hadith.translates
	set translate = :translate 
where id = :id
returning id, lang
`

	// TODO: Use the pgx batch
	for _, tt := range translates {
		var obj types.Translate
		if err := pgscan.Get(ctx, t.db, &obj, query, tt); err != nil {
			return nil, errors.Wrapf(
				err,
				"update translate by id %d",
				tt.ID,
			)
		}

		updated = append(
			updated,
			obj,
		)
	}

	return updated, nil
}
