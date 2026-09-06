package store

import (
	"context"

	"github.com/lattiq/foundry/database"
	"github.com/lattiq/foundry/database/sql"
	"gorm.io/gorm"

	"github.com/lattiq-bhuvan/access-desk/internal/model"
)

type Store struct {
	db *gorm.DB // db handler
}

// New creates a Store that can be used by other packages.
func New(ctx context.Context, cfg *sql.Config) (*Store, error) {
	db, err := database.GormDB(ctx, cfg)
	if err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}

	if err := seed(db); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(&model.Dataset{}, &model.AccessRequest{}, &model.Decision{})
}

func seed(db *gorm.DB) error {
	datasets := []model.Dataset{
		{Slug: "credit-bureau", Name: "Credit Bureau Data", Description: "..."},
		{Slug: "txn-history", Name: "Transaction History", Description: "..."},
		{Slug: "kyc-records", Name: "KYC Records", Description: "..."},
	}
	for _, d := range datasets {
		// FirstOrCreate keyed on slug → re-running boot doesn't duplicate rows
		if err := db.Where(model.Dataset{Slug: d.Slug}).
			Attrs(d).
			FirstOrCreate(&model.Dataset{}).Error; err != nil {
			return err
		}
	}
	return nil
}
