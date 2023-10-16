package types

type Version struct {
	ID        int64      `db:"id"`
	Original  string     `db:"original"`
	BroughtID int64      `db:"brought_id"`
	IsDefault bool       `db:"is_default"`
	Version   Translates `db:"version"`
}
