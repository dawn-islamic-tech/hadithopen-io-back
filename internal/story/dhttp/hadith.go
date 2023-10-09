package dhttp

import (
	"context"
	"net/http"
	"strconv"

	"github.com/ogen-go/ogen/middleware"

	"github.com/hadithopen-io/back/internal/story/dhttp/hadithgen"
	"github.com/hadithopen-io/back/internal/story/types"
	"github.com/hadithopen-io/back/pkg/errors"
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

	hadithgen.UnimplementedHandler
}

func NewStoryHandler(hadith Hadith, transmitters Transmitters, object HadithObject) *StoryHandler {
	return &StoryHandler{
		hadith:       hadith,
		transmitters: transmitters,
		object:       object,
	}
}

const path = "/api/hadith"

func (s *StoryHandler) Path() string { return path }

func (s *StoryHandler) Handler(m ...middleware.Middleware) (http.Handler, error) {
	handler, err := hadithgen.NewServer(
		s,
		hadithgen.WithMiddleware(
			m...,
		),
		hadithgen.WithPathPrefix(
			path,
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
		ID: hadithgen.NewOptInt64(
			hadith.ID,
		),
		Origin: hadithgen.NewOptString(
			hadith.Origin,
		),
		Translate: hadithgen.NewOptString(
			hadith.Translate,
		),
		Interpretation: hadithgen.NewOptString(
			hadith.Interpretation,
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

func (s *StoryHandler) GetTransmitters(ctx context.Context, params hadithgen.GetTransmittersParams) (
	*hadithgen.TransmittersResponse,
	error,
) {
	graph, err := s.transmitters.Get(ctx, params.ID)
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

func (s *StoryHandler) CreateHadith(ctx context.Context, params *hadithgen.HadithObjectRequest) error {
	obj := types.HadithObject{
		Story: types.Story{
			Title: genToCommonTranslates(
				params.Story.Title,
			),
		},
		Brought: types.Brought{
			Brought: genToCommonTranslates(
				params.Brought,
			),
		},
		Comment: types.Comment{
			Comment: genToCommonTranslates(
				params.Comment,
			),
		},
	}

	for _, v := range params.Versions {
		obj.Versions = append(
			obj.Versions,
			types.Version{
				Original:  v.Original.Value,
				IsDefault: v.IsDefault.Value,
				Version: genToCommonTranslates(
					v.Version,
				),
			},
		)
	}

	return s.object.Create(
		ctx,
		obj,
	)
}

func genToCommonTranslates(gen []hadithgen.ObjectTranslate) (ret types.Translates) {
	ret.Values = make([]types.Translate, 0, len(gen))
	for _, t := range gen {
		ret.Values = append(
			ret.Values,
			types.Translate{
				Lang:      t.Lang.Value,
				Translate: t.Translate.Value,
			},
		)
	}

	return ret
}
