package postgres

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/hadithopen-io/back/pkg/errors"
)

type ObjectTX struct{ db *sqlx.DB }

func NewObjectTX(db *sqlx.DB) *ObjectTX { return &ObjectTX{db: db} }

type objectTxKey struct{}

var objectTxKeyValue objectTxKey

func (o *ObjectTX) Wrap(parent context.Context) (ctx context.Context, err error) {
	tx, err := o.db.BeginTxx(
		parent,
		&sql.TxOptions{
			Isolation: sql.LevelReadCommitted,
			ReadOnly:  false,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "beginning postgresql transaction")
	}

	return context.WithValue(
		parent,
		objectTxKeyValue,
		tx,
	), nil
}

var (
	ErrTxNotWrapped   = errors.New("from context not found transaction")
	ErrTxNotSupported = errors.New("tx not supported")
)

func (o *ObjectTX) Commit(ctx context.Context) (err error) {
	tx, err := wrappedTx(ctx)
	if err != nil {
		return err
	}

	return errors.Wrap(
		tx.Commit(),
		"transaction wrapper committing",
	)
}

func (o *ObjectTX) Rollback(ctx context.Context) (err error) {
	tx, err := wrappedTx(ctx)
	if err != nil {
		return err
	}

	return errors.Wrap(
		tx.Rollback(),
		"transaction wrapper committing",
	)
}

func (o *ObjectTX) Get(ctx context.Context) (any, bool) {
	tx, err := wrappedTx(ctx)
	if err != nil {
		return nil, false
	}

	return tx, true
}

func wrappedTx(ctx context.Context) (*sqlx.Tx, error) {
	tx := ctx.Value(objectTxKeyValue)
	if tx == nil {
		return nil, ErrTxNotWrapped
	}

	ttx, ok := tx.(*sqlx.Tx)
	if !ok || ttx == nil {
		return nil, ErrTxNotSupported
	}

	return ttx, nil
}
