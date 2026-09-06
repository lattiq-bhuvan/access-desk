package store

import (
	"context"

	"github.com/lattiq-bhuvan/access-desk/internal/model"
)

func (s *Store) ListDatasets(ctx context.Context) ([]model.Dataset, error) {
	var datasets []model.Dataset
	err := s.db.WithContext(ctx).Order("name").Find(&datasets).Error
	return datasets, err
}
