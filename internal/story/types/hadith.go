package types

type Hadith struct {
	ID             int64  `db:"id"`
	Title          string `db:"title"`
	Origin         string `db:"original"`
	Translate      string `db:"version"`
	Interpretation string `db:"comment"`
	Lang           string `db:"lang"`
}

type HadithObject struct {
	ID       int64 `db:"id"`
	Story    Story
	Versions []Version
	Brought  Brought
	Comment  Comment
}

type HadithCompact struct {
	ID    int64  `db:"id"`
	Title string `db:"title"`
	Desc  string `db:"description"`
}
