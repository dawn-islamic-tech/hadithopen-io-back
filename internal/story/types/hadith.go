package types

type Hadith struct {
	HadithStory
	HadithVersion
	HadithComment
}

type HadithStory struct {
	StoryID          int64  `db:"story_id"`
	Story            string `db:"story_title"`
	StoryTranslateID int64  `db:"story_translate_id"`
	StoryLang        string `db:"story_lang"`
}

type HadithVersion struct {
	VersionID          int64  `db:"version_id"`
	Original           string `db:"version_original"`
	IsDefault          bool   `db:"version_is_default"`
	Version            string `db:"version"`
	VersionTranslateID int64  `db:"version_translate_id"`
	VersionLang        string `db:"version_lang"`
	VersionBroughtID   int64  `db:"version_brought_id"`
}

type HadithComment struct {
	CommentID          int64  `db:"comment_id"`
	Comment            string `db:"comment"`
	CommentTranslateID int64  `db:"comment_translate_id"`
	CommentLang        string `db:"comment_lang"`
	CommentBroughtID   int64  `db:"comment_brought_id"`
}

type HadithObject struct {
	Story    Story
	Versions []Version
	Comment  Comment
}

type HadithCompact struct {
	ID    int64  `db:"id"`
	Title string `db:"title"`
	Desc  string `db:"description"`
}
