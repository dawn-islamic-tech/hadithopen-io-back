package dhttp

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ogen-go/ogen/middleware"

	"github.com/hadithopen-io/back/internal/story/dhttp/hadithgen"
	"github.com/hadithopen-io/back/internal/story/types"
	"github.com/hadithopen-io/back/pkg/errors"
	"github.com/hadithopen-io/back/pkg/jwt"
)

type Hadith interface {
	Few(ctx context.Context) (
		compacts []types.HadithCompact,
		err error,
	)

	Get(ctx context.Context, id int64) (
		hadith types.Hadith,
		err error,
	)
}

type HadithObject interface {
	Create(ctx context.Context, object types.HadithObject) (
		err error,
	)

	Update(ctx context.Context, object types.HadithObject) (
		err error,
	)

	MarkDelete(ctx context.Context, id int64) (
		err error,
	)
}

type Transmitters interface {
	Get(ctx context.Context, hadithID int64) (
		graph types.Graph,
		err error,
	)
}

type StoryHandler struct {
	hadith       Hadith
	transmitters Transmitters
	object       HadithObject
	authWrapper  jwt.Wrapper

	hadithgen.UnimplementedHandler
}

func NewStoryHandler(
	hadith Hadith,
	transmitters Transmitters,
	object HadithObject,
	authWrapper jwt.Wrapper,
) *StoryHandler {
	return &StoryHandler{
		hadith:       hadith,
		transmitters: transmitters,
		object:       object,
		authWrapper:  authWrapper,
	}
}

const path = "/api/hadith"

func (s *StoryHandler) Path() string { return path }

func (s *StoryHandler) HandleCookieAuth(ctx context.Context, _ string, token hadithgen.CookieAuth) (
	context.Context,
	error,
) {
	return s.authWrapper.WrapUser(ctx, token.APIKey)
}

func (s *StoryHandler) Handler(m ...middleware.Middleware) (http.Handler, error) {
	handler, err := hadithgen.NewServer(
		s, // server implementation
		s, // secure middleware implementation
		hadithgen.WithMiddleware(
			m..., // global middlewares
		),
		hadithgen.WithPathPrefix(
			path, // basic prefix
		),
	)

	return handler, errors.Wrap(
		err,
		"make http hadith handler",
	)
}

func (s *StoryHandler) GetHadith(ctx context.Context) (hadithgen.HadithCardsResponse, error) {
	few, err := s.hadith.Few(ctx)
	if err != nil {
		return nil, err
	}

	res := make(hadithgen.HadithCardsResponse, 0, len(few))
	for _, v := range few {
		res = append(res, hadithgen.HadithCard{
			ID: hadithgen.NewOptInt64(
				v.ID,
			),
			Title: hadithgen.NewOptString(
				v.Title,
			),
			Desc: hadithgen.NewOptString(
				v.Desc,
			),
		})
	}

	return res, nil
}

func (s *StoryHandler) GetHadithByID(ctx context.Context, params hadithgen.GetHadithByIDParams) (
	*hadithgen.HadithResponse,
	error,
) {
	hadith, err := s.hadith.Get(ctx, params.ID)
	if err != nil {
		return nil, err
	}

	return &hadithgen.HadithResponse{
		Story: hadithgen.NewOptGetHadithStory(
			hadithgen.GetHadithStory{
				ID: hadithgen.NewOptInt64(
					hadith.StoryID,
				),
				Title: hadithgen.NewOptObjectTranslate(
					hadithgen.ObjectTranslate{
						ID: hadithgen.NewOptInt64(
							hadith.StoryTranslateID,
						),
						Lang: hadithgen.NewOptString(
							hadith.StoryLang,
						),
						Translate: hadithgen.NewOptString(
							hadith.Story,
						),
					},
				),
			},
		),
		Version: hadithgen.NewOptGetHadithVersion(
			hadithgen.GetHadithVersion{
				ID: hadithgen.NewOptInt64(
					hadith.VersionID,
				),
				Original: hadithgen.NewOptString(
					hadith.Original,
				),
				IsDefault: hadithgen.NewOptBool(
					hadith.IsDefault,
				),
				BroughtId: hadithgen.NewOptInt64(
					hadith.VersionBroughtID,
				),
				Version: hadithgen.NewOptObjectTranslate(
					hadithgen.ObjectTranslate{
						ID: hadithgen.NewOptInt64(
							hadith.VersionTranslateID,
						),
						Lang: hadithgen.NewOptString(
							hadith.VersionLang,
						),
						Translate: hadithgen.NewOptString(
							hadith.Version,
						),
					},
				),
			},
		),
		Comment: hadithgen.NewOptGetHadithComment(
			hadithgen.GetHadithComment{
				ID: hadithgen.NewOptInt64(
					hadith.CommentID,
				),
				BroughtId: hadithgen.NewOptInt64(
					hadith.CommentBroughtID,
				),
				Comment: hadithgen.NewOptObjectTranslate(
					hadithgen.ObjectTranslate{
						ID: hadithgen.NewOptInt64(
							hadith.CommentTranslateID,
						),
						Lang: hadithgen.NewOptString(
							hadith.CommentLang,
						),
						Translate: hadithgen.NewOptString(
							hadith.Comment,
						),
					},
				),
			},
		),
	}, nil
}

func (s *StoryHandler) GetSearchedHadith(context.Context, hadithgen.GetSearchedHadithParams) (
	hadithgen.HadithCardsResponse,
	error,
) {
	return hadithgen.HadithCardsResponse{}, nil
}

func (s *StoryHandler) GetSearchedTags(context.Context, hadithgen.GetSearchedTagsParams) (
	hadithgen.HadithTagsResponse,
	error,
) {
	return hadithgen.HadithTagsResponse{}, nil
}

func (s *StoryHandler) GetTransmitters(ctx context.Context, req hadithgen.GetTransmittersParams) (
	*hadithgen.TransmittersResponse,
	error,
) {
	graph, err := s.transmitters.Get(ctx, req.ID)
	if err != nil {
		return nil, errors.Wrap(err, "transmitters getting")
	}

	seqs := make([]hadithgen.Seq, 0, len(graph.Nodes))
	for _, v := range graph.Nodes {
		seqs = append(seqs, hadithgen.Seq{
			ID: hadithgen.NewOptInt64(
				v.ID,
			),
			Name: hadithgen.NewOptString(
				v.Name,
			),
		})
	}

	lines := make(hadithgen.TransmittersResponseLines, len(graph.Edges))
	for k, v := range graph.Edges {
		lines[strconv.Itoa(int(k))] = v
	}

	return &hadithgen.TransmittersResponse{
		Seqs: seqs,
		Lines: hadithgen.NewOptTransmittersResponseLines(
			lines,
		),
	}, nil
}

func (s *StoryHandler) CreateHadith(ctx context.Context, req *hadithgen.HadithObjectRequest) error {
	obj := types.HadithObject{
		Story: types.Story{
			Title: genTranslates(
				req.Title,
			).Common(),
		},
		Comment: types.Comment{
			BroughtID: req.Comment.BroughtId.Value,
			Comment: genTranslates(
				req.Comment.Translates,
			).Common(),
		},
		Versions: genVersions(
			req.Versions,
		).Common(),
	}

	return s.object.Create(
		ctx,
		obj,
	)
}

func (s *StoryHandler) UpdateHadithByID(
	ctx context.Context,
	req *hadithgen.HadithObjectRequest,
	params hadithgen.UpdateHadithByIDParams,
) error {
	obj := types.HadithObject{
		Story: types.Story{
			ID: params.ID,
			Title: genTranslates(
				req.Title,
			).Common(),
		},
		Comment: types.Comment{
			ID:        req.Comment.ID.Value,
			BroughtID: req.Comment.BroughtId.Value,
			Comment: genTranslates(
				req.Comment.Translates,
			).Common(),
		},
		Versions: genVersions(
			req.Versions,
		).Common(),
	}

	return s.object.Update(
		ctx,
		obj,
	)
}

func (s *StoryHandler) MarkDeleteHadithByID(
	ctx context.Context,
	params hadithgen.MarkDeleteHadithByIDParams,
) error {
	return s.object.MarkDelete(
		ctx,
		params.ID,
	)
}

type genTranslates []hadithgen.ObjectTranslate

func (g genTranslates) Common() (ret types.Translates) {
	ret.Values = make([]types.Translate, 0, len(g))
	for _, t := range g {
		ret.Values = append(
			ret.Values,
			types.Translate{
				ID:        t.ID.Value,
				Lang:      t.Lang.Value,
				Translate: t.Translate.Value,
			},
		)
	}

	return ret
}

type genVersions []hadithgen.HadithVersion

func (g genVersions) Common() []types.Version {
	res := make([]types.Version, 0, len(g))
	for _, v := range g {
		res = append(
			res,
			types.Version{
				ID:        v.ID.Value,
				Original:  v.Original,
				BroughtID: v.BroughtId.Value,
				IsDefault: v.IsDefault,
				Version: genTranslates(
					v.Version,
				).Common(),
			},
		)
	}

	return res
}
