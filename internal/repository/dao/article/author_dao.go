package article

import (
	"context"
	"gorm.io/gorm"
)

type AuthorDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
	UpdateOrInsert(ctx context.Context, entity Article) error
}

func NewAuthorDAO(db *gorm.DB) AuthorDAO {
	panic("implement me")
}
