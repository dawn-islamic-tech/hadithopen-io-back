package types

type Hadith struct {
	ID             int64
	Origin         string
	Translate      string
	Interpretation string
}

type HadithCompact struct {
	ID    int64
	Title string
	Desc  string
}
