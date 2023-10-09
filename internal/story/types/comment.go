package types

type Comment struct {
	ID        int64      `db:"id"`
	StoryID   int64      `db:"story_id"`
	BroughtID int64      `db:"brought_id"`
	Comment   Translates `db:"comment"`
}
