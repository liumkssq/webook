package article

import (
	"context"
	"gorm.io/gorm"
)

type ReaderDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
	UpdateOrInsert(ctx context.Context, entity Article) error
}

func NewReaderDAO(db *gorm.DB) ReaderDAO {
	panic("implement me")
}
