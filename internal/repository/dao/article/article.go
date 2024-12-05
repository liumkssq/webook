package article

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"time"
)

type ArticleDAO interface {
	Insert(ctx context.Context, art Article) (int64, error)
	UpdateById(ctx context.Context, article Article) error
}

type GORMArticleDAO struct {
	db *gorm.DB
}

func (dao *GORMArticleDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Utime = now
	art.Ctime = now
	err := dao.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

func (dao *GORMArticleDAO) UpdateById(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	//不用默认更新
	res := dao.db.WithContext(ctx).Model(&art).
		Where("id=? and author_id = ?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"utime":   now,
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("更新失败")
	}
	return res.Error
}

func NewGORMArticleDAO() ArticleDAO {
	return &GORMArticleDAO{}
}

type Article struct {
	Id      int64  `gorm:"primaryKey,autoIncrement"`
	Title   string `gorm:"type:varchar(128)"`
	Content string `gorm:"BLOB"`
	//Index
	AuthorId int64 `gorm:"index"`
	//AuthorId int64 	`gorm:"index=aid_ctime"`
	//Ctime    int64 	`gorm:"index=aid_ctime"`
	Ctime int64
	Utime int64
}
