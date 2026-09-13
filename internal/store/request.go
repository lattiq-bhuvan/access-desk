package store

import (
      "context"

      "gorm.io/gorm"
      "gorm.io/gorm/clause"

      "github.com/lattiq-bhuvan/access-desk/internal/model"
)

// returns nil if request is successfully registered in the DB
func (s *Store) CreateRequest(ctx context.Context, r *model.AccessRequest) error {
	return s.db.WithContext(ctx).Create(r).Error
}

func (s *Store) GetRequest(ctx context.Context, id uint) (*model.AccessRequest, error) {
	var r model.AccessRequest
	err := s.db.WithContext(ctx).Preload("Requester").First(&r, id).Error
	if err != nil {
		return nil, err // caller checks gorm.ErrRecordNotFound
	}
	r.RequesterEmail = r.Requester.Email
	return &r, nil
}

// ListRequests: onlyRequester != 0 scopes to that user's id. Returns rows + total count.
func (s *Store) ListRequests(ctx context.Context, onlyRequester uint, status string, page, pageSize int) ([]model.AccessRequest, int64, error) {
	q := s.db.WithContext(ctx).Model(&model.AccessRequest{})
	// filter by requester
	if onlyRequester != 0 {
		q = q.Where("requester_id = ?", onlyRequester)
	}

	// filter by status(pending, approved, rejected)
	if status != "" {
		q = q.Where("status = ?", status)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []model.AccessRequest
	err := q.Preload("Requester").Order("created_at desc").
		Limit(pageSize).Offset((page - 1) * pageSize).
		Find(&rows).Error
	for i := range rows {
		rows[i].RequesterEmail = rows[i].Requester.Email
	}
	return rows, total, err
}

// DecideRequest runs the state transition + audit row in one transaction.
// checkFn receives the locked current row and returns an API error if the
// transition is illegal (terminal state, self-approval, etc.).
func (s *Store) DecideRequest(
	ctx context.Context,
	id uint,
	checkFn func(cur *model.AccessRequest) error,
	newStatus string,
	decision *model.Decision,
) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var cur model.AccessRequest
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&cur, id).Error; err != nil {
			return err
		}
		if err := checkFn(&cur); err != nil {
			return err
		}
		if err := tx.Model(&cur).Update("status", newStatus).Error; err != nil {
			return err
		}
		decision.AccessRequestID = id
		return tx.Create(decision).Error
	})
}