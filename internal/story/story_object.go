package story

import (
	"context"

	"github.com/hadithopen-io/back/internal/story/types"
	"github.com/hadithopen-io/back/pkg/errors"
	"github.com/hadithopen-io/back/pkg/tx"
)

type Translate interface {
	Create(ctx context.Context, translates []types.Translate) (
		created types.Translates,
		err error,
	)
}

type Comment interface {
	Create(ctx context.Context, comment types.Comment) (
		id int64,
		err error,
	)
}

type Brought interface {
	Create(ctx context.Context, brought types.Brought) (
		id int64,
		err error,
	)
}

type Version interface {
	Create(ctx context.Context, version types.Version) (
		id int64,
		err error,
	)
}

type ObjectStore interface {
	Create(ctx context.Context, story types.Story) (
		id int64,
		err error,
	)

	CreateMapVersion(ctx context.Context, storyID int64, versionID ...int64) (
		err error,
	)
}

type Object struct {
	translate Translate
	comment   Comment
	brought   Brought
	version   Version
	store     ObjectStore
	tx        tx.Wrapper
}

func NewObject(
	translate Translate,
	comment Comment,
	brought Brought,
	version Version,
	store ObjectStore,
	tx tx.Wrapper,
) *Object {
	return &Object{
		translate: translate,
		comment:   comment,
		brought:   brought,
		version:   version,
		store:     store,
		tx:        tx,
	}
}

func (o Object) Create(ctx context.Context, object types.HadithObject) error {
	return o.tx.Wrap(
		ctx,
		func(ctx context.Context) error {
			return o.create(
				ctx,
				object,
			)
		},
	)
}

func (o Object) create(ctx context.Context, object types.HadithObject) error {
	storyID, err := o.createStory(
		ctx,
		object.Story,
	)
	if err != nil {
		return err
	}

	broughtID, err := o.createBrought(
		ctx,
		object.Brought,
	)
	if err != nil {
		return err
	}

	err = o.createComment(
		ctx,
		storyID,
		broughtID,
		object.Comment,
	)
	if err != nil {
		return err
	}

	versionsID, err := o.createVersions(
		ctx,
		broughtID,
		object.Versions,
	)
	if err != nil {
		return err
	}

	return errors.Wrap(
		o.store.CreateMapVersion(
			ctx,
			storyID,
			versionsID...,
		),
		"creating hadith versions map",
	)
}

func (o Object) createStory(ctx context.Context, story types.Story) (storyID int64, err error) {
	translates, err := o.translate.Create(
		ctx,
		story.Title.Values,
	)
	if err != nil {
		return 0, errors.Wrap(
			err,
			"creating hadith story title translates",
		)
	}

	story.Title = translates

	objID, err := o.store.Create(
		ctx,
		story,
	)
	return objID, errors.Wrap(
		err,
		"creating hadith story",
	)
}

func (o Object) createComment(ctx context.Context, storyID, broughtID int64, comment types.Comment) error {
	translates, err := o.translate.Create(
		ctx,
		comment.Comment.Values,
	)
	if err != nil {
		return errors.Wrap(
			err,
			"creating hadith comment translates",
		)
	}

	comment.Comment = translates
	comment.StoryID = storyID
	comment.BroughtID = broughtID

	_, err = o.comment.Create(ctx, comment)
	return errors.Wrap(
		err,
		"creating hadith comment",
	)
}

func (o Object) createBrought(ctx context.Context, brought types.Brought) (broughtID int64, err error) {
	if brought.ID != 0 {
		return brought.ID, nil
	}

	translates, err := o.translate.Create(
		ctx,
		brought.Brought.Values,
	)
	if err != nil {
		return 0, errors.Wrap(
			err,
			"creating hadith brought translates",
		)
	}

	brought.Brought = translates

	broughtID, err = o.brought.Create(
		ctx,
		brought,
	)
	return broughtID, errors.Wrap(
		err,
		"creating hadith brought",
	)
}

func (o Object) createVersions(ctx context.Context, broughtID int64, versions []types.Version) (
	versionsID []int64,
	err error,
) {
	versionsID = make([]int64, 0, len(versions))
	for _, version := range versions {
		versionTranslates, err := o.translate.Create(
			ctx,
			version.Version.Values,
		)
		if err != nil {
			return nil, errors.Wrap(err, "creating hadith version translates")
		}

		version.Version = versionTranslates
		version.BroughtID = broughtID

		id, err := o.version.Create(
			ctx,
			version,
		)
		if err != nil {
			return nil, errors.Wrap(err, "creating hadith object version")
		}

		versionsID = append(
			versionsID,
			id,
		)
	}

	return versionsID, nil
}
