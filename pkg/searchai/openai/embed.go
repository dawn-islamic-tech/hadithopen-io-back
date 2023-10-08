package openai

import (
	"context"

	"github.com/tmc/langchaingo/embeddings"

	"github.com/hadithopen-io/back/pkg/errors"
	"github.com/hadithopen-io/back/pkg/searchai"
)

type Embed struct {
	embeder embeddings.Embedder
}

func NewEmbed(embeder embeddings.Embedder) *Embed {
	return &Embed{
		embeder: embeder,
	}
}

func (e Embed) Query(ctx context.Context, query string) (
	vec searchai.Vector,
	err error,
) {
	pre, err := e.embeder.EmbedQuery(ctx, query)
	if err != nil {
		return nil, errors.Wrap(err, "openai embed query")
	}

	return formatSlice[searchai.Vector](pre), nil
}

func (e Embed) Embed(ctx context.Context, documents []searchai.Document) (
	vec searchai.Vectors,
	err error,
) {
	texts := make([]string, 0, len(documents))
	for _, doc := range documents {
		texts = append(texts, doc.Text)
	}

	pre, err := e.embeder.EmbedDocuments(ctx, texts)
	if err != nil {
		return nil, errors.Wrap(err, "embedding documents")
	}

	vecs := make(searchai.Vectors, 0, len(pre))
	for _, p := range pre {
		vecs = append(vecs, formatSlice[searchai.Vector](p))
	}

	return vecs, nil
}

func formatSlice[T ~[]E, K ~[]Z, E ~float64 | float32, Z ~float64 | float32](a K) T {
	to := make(T, 0, len(a))
	for _, t := range a {
		to = append(to, E(t))
	}

	return to
}
