package postgre

import (
	"context"

	"github.com/hadithopen-io/back/internal/story/types"
)

type Hadith struct {
}

func NewHadith() *Hadith {
	return &Hadith{}
}

func (h Hadith) Few(ctx context.Context) (compacts []types.HadithCompact, err error) {
	//TODO implement me
	panic("implement me")
}

func (h Hadith) Get(ctx context.Context, id int64) (hadith types.Hadith, err error) {
	//TODO implement me
	panic("implement me")
}
