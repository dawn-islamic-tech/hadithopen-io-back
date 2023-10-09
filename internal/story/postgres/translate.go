package postgres

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/hadithopen-io/back/internal/story/types"
	"github.com/hadithopen-io/back/pkg/pgscan"
)

type Translate struct{ db *sqlx.DB }

func NewTranslate(db *sqlx.DB) *Translate { return &Translate{db: db} }

func (t Translate) Create(ctx context.Context, translates []types.Translate) (
	ret types.Translates,
	err error,
) {
	const query = `
insert into hadith.translates(lang, translate)
values (:lang, :translate)
returning id, lang
	`

	return ret, pgscan.Select(
		ctx,
		t.db,
		&ret.Values,
		query,
		translates,
	)
}
