package tx

import (
	"context"

	"github.com/hadithopen-io/back/pkg/errors"
)

type Transaction interface {
	Wrap(parent context.Context) (
		ctx context.Context,
		err error,
	)

	Commit(ctx context.Context) (
		err error,
	)

	Rollback(ctx context.Context) (
		err error,
	)

	Get(ctx context.Context) (
		tx any,
		ok bool,
	)
}

type Wrapper interface {
	Wrap(ctx context.Context, fn func(context.Context) error) (
		err error,
	)
}

type wrapper struct {
	tx Transaction
}

func NewWrapper(tx Transaction) Wrapper {
	return &wrapper{
		tx: tx,
	}
}

func (w *wrapper) Wrap(ctx context.Context, fn func(context.Context) error) (
	err error,
) {
	ctx, err = w.tx.Wrap(
		ctx,
	)
	if err != nil {
		return err
	}

	defer func() {
		if err == nil {
			err = w.tx.Commit(
				ctx,
			)

			return
		}

		err = errors.Join(
			err,
			w.tx.Rollback(
				ctx,
			),
		)
	}()

	tx, ok := w.tx.Get(
		ctx,
	)
	if !ok {
		return ErrTxNotSet
	}

	ctx = wrap(
		ctx,
		tx,
	)

	return fn(
		ctx,
	)
}

type wrapReadableCtxKey struct{}

var wrapReadableCtxKeyValue wrapReadableCtxKey

var ErrTxNotSet = errors.New("not set original tx")

func wrap(ctx context.Context, tx any) context.Context {
	return context.WithValue(
		ctx,
		wrapReadableCtxKeyValue,
		tx,
	)
}

func GetContext(ctx context.Context) (any, bool) {
	v := ctx.Value(
		wrapReadableCtxKeyValue,
	)
	if v == nil {
		return nil, false
	}

	return v, true
}
