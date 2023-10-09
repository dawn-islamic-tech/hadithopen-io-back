package types

type Brought struct {
	ID      int64      `db:"id"`
	Brought Translates `db:"brought"`
}
