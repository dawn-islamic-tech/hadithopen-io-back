package story

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"

	"github.com/hadithopen-io/back/internal/story/mocks"
	"github.com/hadithopen-io/back/internal/story/types"
)

//go:generate mockgen -destination mocks/object_story.go -source=story_object.go -package mocks Translate, Comment, Brought, Version, ObjectStore
func Test_Object(t *testing.T) {
	suite.Run(
		t,
		new(
			TestObjectStore,
		),
	)
}

type TestObjectStore struct {
	suite.Suite

	controller *gomock.Controller

	translate *mocks.MockTranslate
	comment   *mocks.MockComment
	brought   *mocks.MockBrought
	version   *mocks.MockVersion
	object    *mocks.MockObjectStore
}

func (t *TestObjectStore) SetupTest() {
	t.controller = gomock.NewController(
		t.T(),
	)

	t.translate = mocks.NewMockTranslate(
		t.controller,
	)

	t.comment = mocks.NewMockComment(
		t.controller,
	)

	t.brought = mocks.NewMockBrought(
		t.controller,
	)

	t.version = mocks.NewMockVersion(
		t.controller,
	)

	t.object = mocks.NewMockObjectStore(
		t.controller,
	)
}

type mockWrapper struct{}

func (mockWrapper) Wrap(ctx context.Context, fn func(context.Context) error) error {
	return fn(
		ctx,
	)
}

func (t *TestObjectStore) TestCreate() {
	defer t.controller.Finish()

	t.translate.
		EXPECT().
		Create(
			context.Background(),
			[]types.Translate{
				{
					Lang:      "en",
					Translate: "Hello",
				},
			},
		).
		Return(
			types.Translates{
				Values: []types.Translate{
					{
						ID:        1,
						Lang:      "en",
						Translate: "Hello",
					},
				},
			},
			nil,
		)

	t.object.
		EXPECT().
		Create(
			context.Background(),
			types.Story{
				Title: types.Translates{
					Values: []types.Translate{
						{
							ID:        1,
							Lang:      "en",
							Translate: "Hello",
						},
					},
				},
			},
		).
		Return(
			int64(1),
			nil,
		)

	t.translate.
		EXPECT().
		Create(
			context.Background(),
			[]types.Translate{
				{
					Lang:      "en",
					Translate: "Brought",
				},
			},
		).
		Return(
			types.Translates{
				Values: []types.Translate{
					{
						ID:        2,
						Lang:      "en",
						Translate: "Brought",
					},
				},
			},
			nil,
		)

	t.brought.
		EXPECT().
		Create(
			context.Background(),
			types.Brought{
				Brought: types.Translates{
					Values: []types.Translate{
						{
							ID:        2,
							Lang:      "en",
							Translate: "Brought",
						},
					},
				},
			},
		).
		Return(
			int64(1),
			nil,
		)

	t.translate.
		EXPECT().
		Create(
			context.Background(),
			[]types.Translate{
				{
					Lang:      "en",
					Translate: "Comment",
				},
			},
		).
		Return(
			types.Translates{
				Values: []types.Translate{
					{
						ID:        3,
						Lang:      "en",
						Translate: "Comment",
					},
				},
			},
			nil,
		)

	t.comment.
		EXPECT().
		Create(
			context.Background(),
			types.Comment{
				StoryID:   1,
				BroughtID: 1,
				Comment: types.Translates{
					Values: []types.Translate{
						{
							ID:        3,
							Lang:      "en",
							Translate: "Comment",
						},
					},
				},
			},
		).
		Return(
			int64(1),
			nil,
		)

	t.translate.
		EXPECT().
		Create(
			context.Background(),
			[]types.Translate{
				{
					Lang:      "en",
					Translate: "Version",
				},
			},
		).
		Return(
			types.Translates{
				Values: []types.Translate{
					{
						ID:        4,
						Lang:      "en",
						Translate: "Version",
					},
				},
			},
			nil,
		)

	t.version.
		EXPECT().
		Create(
			context.Background(),
			types.Version{
				Original:  "Original",
				BroughtID: 1,
				IsDefault: true,
				Version: types.Translates{
					Values: []types.Translate{
						{
							ID:        4,
							Lang:      "en",
							Translate: "Version",
						},
					},
				},
			},
		).
		Return(
			int64(1),
			nil,
		)

	t.object.
		EXPECT().
		CreateMapVersion(
			context.Background(),
			int64(1),
			int64(1),
		).
		Return(
			nil,
		)

	obj := NewObject(
		t.translate,
		t.comment,
		t.brought,
		t.version,
		t.object,
		mockWrapper{},
	)

	err := obj.Create(
		context.Background(),
		types.HadithObject{
			Story: types.Story{
				Title: types.Translates{
					Values: []types.Translate{
						{
							Lang:      "en",
							Translate: "Hello",
						},
					},
				},
			},
			Versions: []types.Version{
				{
					Original:  "Original",
					IsDefault: true,
					Version: types.Translates{
						Values: []types.Translate{
							{
								Lang:      "en",
								Translate: "Version",
							},
						},
					},
				},
			},
			Brought: types.Brought{
				Brought: types.Translates{
					Values: []types.Translate{
						{
							Lang:      "en",
							Translate: "Brought",
						},
					},
				},
			},
			Comment: types.Comment{
				Comment: types.Translates{
					Values: []types.Translate{
						{
							Lang:      "en",
							Translate: "Comment",
						},
					},
				},
			},
		},
	)
	require.NoError(
		t.T(),
		err,
	)
}
