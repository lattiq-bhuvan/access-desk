package store

import (
	"context"

	"github.com/lattiq-bhuvan/access-desk/internal/model"
)

func (s *Store) UserByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (s *Store) UserByID(ctx context.Context, id string) (*model.User, error) {
	var u model.User
	if err := s.db.WithContext(ctx).First(&u, id).Error; err != nil {
			return nil, err
	}
	return &u, nil
}