package translate

import (
	"context"

	"github.com/hadithopen-io/back/internal/story/types"
	"github.com/hadithopen-io/back/pkg/errors"
)

type Translator interface {
	Create(ctx context.Context, translates []types.Translate) (
		created []types.Translate,
		err error,
	)

	Update(ctx context.Context, translates []types.Translate) (
		updated []types.Translate,
		err error,
	)
}

type Translate struct {
	tr Translator
}

// Option
// TODO: use after add search engine
type Option func(*Translate)

func NewTranslate(tr Translator) *Translate {
	return &Translate{
		tr: tr,
	}
}

func (t Translate) Add(ctx context.Context, translates ...types.Translate) (_ types.Translates, err error) {
	create := make([]types.Translate, 0)
	update := make([]types.Translate, 0)
	for _, v := range translates {
		if v.ID == 0 {
			create = append(
				create,
				v,
			)

			continue
		}

		update = append(
			update,
			v,
		)
	}

	created, err := t.tr.Create(
		ctx,
		create,
	)
	if err != nil {
		return types.Translates{}, errors.Wrap(err, "create translation")
	}

	updated, err := t.tr.Update(
		ctx,
		update,
	)
	if err != nil {
		return types.Translates{}, errors.Wrap(err, "update translation")
	}

	return types.Translates{
		Values: append(
			created,
			updated...,
		),
	}, nil
}
