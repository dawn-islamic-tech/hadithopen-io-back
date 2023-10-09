package pgscan

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/hadithopen-io/back/pkg/errors"
	"github.com/hadithopen-io/back/pkg/tx"
)

func Select(
	ctx context.Context,
	db sqlx.ExtContext,
	dest any,
	query string,
	arg any,
) error {
	db = useTx(
		ctx,
		db,
	)

	nq, args, err := Named(
		query,
		arg,
	)
	if err != nil {
		return err
	}

	nq = db.Rebind(
		nq,
	)

	return errors.Wrap(
		sqlx.SelectContext(
			ctx,
			db,
			dest,
			nq,
			args...,
		),
		"sqlx selecting",
	)
}

func Get(
	ctx context.Context,
	db sqlx.ExtContext,
	dest any,
	query string,
	arg any,
) error {
	db = useTx(
		ctx,
		db,
	)

	nq, args, err := Named(
		query,
		arg,
	)
	if err != nil {
		return err
	}

	nq = db.Rebind(
		nq,
	)

	return errors.Wrap(
		sqlx.GetContext(
			ctx,
			db,
			dest,
			nq,
			args...,
		),
		"sqlx getting",
	)
}

func Exec(
	ctx context.Context,
	db sqlx.ExtContext,
	query string,
	arg any,
) error {
	db = useTx(
		ctx,
		db,
	)

	nq, args, err := Named(
		query,
		arg,
	)
	if err != nil {
		return err
	}

	nq = db.Rebind(
		nq,
	)

	_, err = db.ExecContext(
		ctx,
		nq,
		args...,
	)
	return errors.Wrap(
		err,
		"sqlx executing ctx",
	)
}

func Named(
	query string,
	arg any,
) (
	q string,
	args []any,
	err error,
) {
	q, args, err = sqlx.Named(
		query,
		arg,
	)
	return q, args, errors.Wrap(
		err,
		"sqlx query named process",
	)
}

func useTx(
	ctx context.Context,
	db sqlx.ExtContext,
) sqlx.ExtContext {
	wrapped, ok := tx.GetContext(ctx)
	if !ok {
		return db
	}

	tx, _ := wrapped.(*sqlx.Tx)
	return tx
}
