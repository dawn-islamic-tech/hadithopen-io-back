package meilisearch

import (
	"context"

	"github.com/meilisearch/meilisearch-go"

	"github.com/hadithopen-io/back/pkg/errors"
	"github.com/hadithopen-io/back/pkg/searchai"
)

type Search struct {
	index string

	client *meilisearch.Client
}

func NewSearch(index string, client *meilisearch.Client) *Search {
	return &Search{
		index:  index,
		client: client,
	}
}

func (s Search) Index(_ context.Context, docs []searchai.VectorizedDocument) error {
	type doc struct {
		ID   string `json:"id"`
		Text string `json:"text"`

		// Experimental https://github.com/meilisearch/product/discussions/621
		Vector searchai.Vector `json:"_vector"`
	}

	indexdocs := make([]doc, 0, len(docs))
	for _, d := range docs {
		indexdocs = append(indexdocs, doc{
			ID:     d.ID,
			Text:   d.Text,
			Vector: d.Vector,
		})
	}

	info, err := s.client.Index(s.index).AddDocuments(indexdocs)
	if err != nil {
		return errors.Wrap(err, "adding documents")
	}
	if info.Status == meilisearch.TaskStatusFailed {
		return errors.New("adding document status failed")
	}

	return nil
}

func (s Search) Search(_ context.Context, sv searchai.SearchVector) (
	docs []searchai.Document,
	_ error,
) {
	info, err := s.client.Index(s.index).Search("", &meilisearch.SearchRequest{
		Offset: sv.Offset,
		Limit:  sv.Limit,
		Vector: sv.Vector,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "meilisearch search by index %s", s.index)
	}

	const (
		fieldID   = "id"
		fieldText = "text"
	)

	docs = make([]searchai.Document, 0, len(info.Hits))
	for _, initial := range info.Hits {
		hit, ok := initial.(map[string]any)
		if !ok {
			return nil, errors.New("meilisearch searched object don't matched")
		}

		id, _ := hit[fieldID].(string)
		text, _ := hit[fieldText].(string)

		docs = append(
			docs,
			searchai.Document{
				ID:   id,
				Text: text,
			},
		)
	}

	return docs, nil
}
