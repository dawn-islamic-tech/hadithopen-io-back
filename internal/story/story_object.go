package story

import (
	"context"

	"github.com/hadithopen-io/back/internal/story/types"
	"github.com/hadithopen-io/back/pkg/errors"
	"github.com/hadithopen-io/back/pkg/tx"
)

type Translate interface {
	Add(ctx context.Context, translates ...types.Translate) (
		changed types.Translates,
		err error,
	)
}

type Comment interface {
	Create(ctx context.Context, comment types.Comment) (
		id int64,
		err error,
	)

	Update(ctx context.Context, comment types.Comment) (
		err error,
	)

	Get(ctx context.Context, id int64) (
		comment types.Comment,
		err error,
	)
}

type Brought interface {
	Create(ctx context.Context, brought types.Brought) (
		id int64,
		err error,
	)

	Update(ctx context.Context, brought types.Brought) (
		err error,
	)
}

type Version interface {
	Create(ctx context.Context, version types.Version) (
		id int64,
		err error,
	)

	Update(ctx context.Context, version types.Version) (
		err error,
	)

	Get(ctx context.Context, id int64) (
		version types.Version,
		err error,
	)
}

type ObjectStore interface {
	ObjectStoreCreater

	Update(ctx context.Context, story types.Story) (
		err error,
	)

	Get(ctx context.Context, id int64) (
		story types.Story,
		err error,
	)
}

type ObjectStoreCreater interface {
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

	err = o.createComment(
		ctx,
		storyID,
		object.Comment,
	)
	if err != nil {
		return err
	}

	versionsID, err := o.createVersions(
		ctx,
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

func (o Object) createStory(ctx context.Context, story types.Story) (id int64, err error) {
	translates, err := o.translate.Add(
		ctx,
		story.Title.Values...,
	)
	if err != nil {
		return 0, errors.Wrap(
			err,
			"creating hadith story title translates",
		)
	}

	id, err = o.store.Create(
		ctx,
		types.Story{
			Title: translates,
		},
	)
	return id, errors.Wrap(
		err,
		"creating hadith story",
	)
}

func (o Object) createComment(ctx context.Context, storyID int64, comment types.Comment) error {
	translates, err := o.translate.Add(
		ctx,
		comment.Comment.Values...,
	)
	if err != nil {
		return errors.Wrap(
			err,
			"creating hadith comment translates",
		)
	}

	comment.Comment = translates
	comment.StoryID = storyID

	_, err = o.comment.Create(
		ctx,
		comment,
	)
	return errors.Wrap(
		err,
		"creating hadith comment",
	)
}

// nolint:unused
func (o Object) createBrought(ctx context.Context, brought types.Brought) (broughtID int64, err error) {
	if brought.ID != 0 {
		return brought.ID, nil
	}

	translates, err := o.translate.Add(
		ctx,
		brought.Brought.Values...,
	)
	if err != nil {
		return 0, errors.Wrap(
			err,
			"creating hadith brought translates",
		)
	}

	broughtID, err = o.brought.Create(
		ctx,
		types.Brought{
			Brought: translates,
		},
	)
	return broughtID, errors.Wrap(
		err,
		"creating hadith brought",
	)
}

func (o Object) createVersions(ctx context.Context, versions []types.Version) (
	versionsID []int64,
	err error,
) {
	versionsID = make([]int64, 0, len(versions))
	for _, version := range versions {
		id, err := o.createVersion(
			ctx,
			version,
		)
		if err != nil {
			return nil, err
		}

		versionsID = append(
			versionsID,
			id,
		)
	}

	return versionsID, nil
}

func (o Object) createVersion(ctx context.Context, version types.Version) (id int64, err error) {
	version.Version, err = o.translate.Add(
		ctx,
		version.Version.Values...,
	)
	if err != nil {
		return 0, errors.Wrap(
			err,
			"creating hadith version translates",
		)
	}

	id, err = o.version.Create(
		ctx,
		version,
	)
	return id, errors.Wrap(
		err,
		"creating hadith object version",
	)
}

func (o Object) Update(ctx context.Context, object types.HadithObject) error {
	return o.tx.Wrap(
		ctx,
		func(ctx context.Context) error {
			return o.update(
				ctx,
				object,
			)
		},
	)
}

func (o Object) update(ctx context.Context, object types.HadithObject) (err error) {
	err = o.updateStory(
		ctx,
		object.Story,
	)
	if err != nil {
		return err
	}

	err = o.updateComment(
		ctx,
		object.Comment,
	)
	if err != nil {
		return err
	}

	ids, err := o.updateVersions(
		ctx,
		object.Versions,
	)
	if err != nil {
		return err
	}

	return errors.Wrap(
		o.store.CreateMapVersion(
			ctx,
			object.Story.ID,
			ids...,
		),
		"updating hadith versions map",
	)
}

func (o Object) updateStory(ctx context.Context, story types.Story) error {
	old, err := o.store.Get(
		ctx,
		story.ID,
	)
	if err != nil {
		return err
	}

	changed, err := o.translate.Add(
		ctx,
		story.Title.Values...,
	)
	if err != nil {
		return err
	}

	return o.store.Update(
		ctx,
		types.Story{
			ID: story.ID,
			Title: o.combineTranslate(
				old.Title,
				changed,
			),
		},
	)
}

func (o Object) updateComment(ctx context.Context, comment types.Comment) error {
	if comment.ID == 0 { // need check the state is available
		return nil
	}

	c, err := o.comment.Get(ctx, comment.ID)
	if err != nil {
		return err
	}

	changed, err := o.translate.Add(
		ctx,
		comment.Comment.Values...,
	)
	if err != nil {
		return err
	}

	c.Comment = o.combineTranslate(
		c.Comment,
		changed,
	)

	return o.comment.Update(
		ctx,
		c,
	)
}

func (o Object) updateVersions(ctx context.Context, versions []types.Version) (ids []int64, _ error) {
	for _, version := range versions {
		if version.ID == 0 {
			id, err := o.createVersion(
				ctx,
				version,
			)
			if err != nil {
				return nil, err
			}

			ids = append(
				ids,
				id,
			)
			continue
		}

		if err := o.updateVersion(ctx, version); err != nil {
			return nil, err
		}
	}

	return ids, nil
}

func (o Object) updateVersion(ctx context.Context, version types.Version) error {
	v, err := o.version.Get(ctx, version.ID)
	if err != nil {
		return err
	}

	changed, err := o.translate.Add(
		ctx,
		version.Version.Values...,
	)
	if err != nil {
		return err
	}

	version.Version = o.combineTranslate(
		v.Version,
		changed,
	)

	return o.version.Update(
		ctx,
		version,
	)
}

func (Object) combineTranslate(old, new types.Translates) types.Translates {
	return types.Translates{
		Values: append(
			old.Values,
			new.Values...,
		),
	}
}
