package searchai

import (
	"context"

	"github.com/go-faster/errors"
)

type SearchAI struct {
	embedding Embedding
	index     IndexSearcher
}

func New(embedding Embedding, index IndexSearcher) *SearchAI {
	return &SearchAI{
		embedding: embedding,
		index:     index,
	}
}

type Document struct {
	ID   string
	Text string
}

type SearchVector struct {
	Limit,
	Offset int64

	Vector
}

type Vector []float64

type Vectors []Vector

type Embedding interface {
	Query(ctx context.Context, query string) (
		vec Vector,
		err error,
	)

	Embed(ctx context.Context, document []Document) (
		vec Vectors,
		err error,
	)
}

type VectorizedDocument struct {
	Vector
	Document
}

type IndexSearcher interface {
	Indexer
	Searcher
}

type Indexer interface {
	Index(ctx context.Context, docs []VectorizedDocument) (
		err error,
	)
}

type Searcher interface {
	Search(ctx context.Context, vec SearchVector) (
		docs []Document,
		err error,
	)
}

func (s *SearchAI) Index(ctx context.Context, documents []Document) error {
	vecs, err := s.embedding.Embed(ctx, documents)
	if err != nil {
		return err
	}

	vd := make([]VectorizedDocument, 0, len(vecs))
	for i, doc := range documents {
		vd = append(
			vd,
			VectorizedDocument{
				Vector: vecs[i],
				Document: Document{
					ID:   doc.ID,
					Text: doc.Text,
				},
			},
		)
	}

	return s.index.Index(
		ctx,
		vd,
	)
}

type SearchParams struct {
	Limit,
	Offset int64

	Query string
}

func (s *SearchAI) Search(ctx context.Context, params SearchParams) (
	docs []Document,
	err error,
) {
	vec, err := s.embedding.Query(
		ctx,
		params.Query,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "vectorize query %s", params.Query)
	}

	return s.index.Search(
		ctx,
		SearchVector{
			Limit:  params.Limit,
			Offset: params.Offset,
			Vector: vec,
		},
	)
}
