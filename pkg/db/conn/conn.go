package conn

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/hadithopen-io/back/pkg/errors"
)

func NewConn(ctx context.Context, conn string) (*sqlx.DB, error) {
	db, err := sqlx.ConnectContext(ctx, "pgx", conn)
	return db, errors.Wrap(err, "connection with pgx")
}
