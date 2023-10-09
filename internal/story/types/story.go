package types

type Story struct {
	ID    int64      `db:"id"`
	Title Translates `db:"title"`
}
