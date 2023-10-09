package story

import (
	"context"

	"github.com/hadithopen-io/back/internal/story/types"
)

type HadithStore interface {
	Few(ctx context.Context) (
		compacts []types.HadithCompact,
		err error,
	)

	HadithGetter
}

type HadithCache interface {
	HadithGetter
}

type HadithGetter interface {
	Get(ctx context.Context, id int64) (
		hadith types.Hadith,
		err error,
	)
}

type Story struct {
	hadith struct {
		store HadithStore
		cache HadithCache
	}
}

func NewStory(
	store HadithStore,
	cache HadithCache,
) *Story {
	return &Story{
		hadith: struct {
			store HadithStore
			cache HadithCache
		}{
			store: store,
			cache: cache,
		},
	}
}

func (s Story) Few(ctx context.Context) (
	compacts []types.HadithCompact,
	err error,
) {
	return s.hadith.store.Few(
		ctx,
	)
}

func (s Story) Get(ctx context.Context, id int64) (
	hadith types.Hadith,
	err error,
) {
	return s.hadith.store.Get(
		ctx,
		id,
	)
}
