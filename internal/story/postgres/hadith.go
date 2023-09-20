package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"

	"github.com/hadithopen-io/back/internal/story/types"
)

type Hadith struct {
	db *pgx.Conn
}

func NewHadith(db *pgx.Conn) *Hadith {
	return &Hadith{
		db: db,
	}
}

func (h Hadith) Few(context.Context) (compacts []types.HadithCompact, err error) {
	return nil, err
}

func (h Hadith) Get(context.Context, int64) (hadith types.Hadith, err error) {
	return hadith, nil
}
