package notes

import (
	"context"
	"errors"
	"personalKnowledgeSearchEngine/internal/es"
	"personalKnowledgeSearchEngine/internal/models"
)

type Service struct {
	ctx context.Context
	es  *es.ESClient
}

func NewService(ctx context.Context, es *es.ESClient) *Service {
	return &Service{
		ctx: ctx,
		es:  es,
	}
}

func (s *Service) Create(note *models.Note) error {
	if note == nil {
		return errors.New("invalid note")
	}

	err := s.es.IndexNote(s.ctx, note)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) SearchNotes(query string) ([]models.Note, error) {
	if query == "" {
		return nil, errors.New("invalid query")
	}

	searchQuery := &models.SearchQuery{
		Query: models.MultiMatchContainer{
			MultiMatch: models.MultiMatchQuery{
				Query:  query,
				Fields: []string{"title", "content"},
			},
		},
	}
	notes, err := s.es.SearchNotes(searchQuery)
	if err != nil {
		return nil, err
	}

	return notes, nil
}
