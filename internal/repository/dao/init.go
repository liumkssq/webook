package dao

import (
	"github.com/liumkssq/webook/internal/repository/dao/article"
	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &article.Article{},
		&article.PublishedArticle{},
		&article.PublishedArticleV1{},
		&AsyncSms{},
		&Interactive{},
		&UserLikeBiz{},
		&Collection{},
		&UserCollectionBiz{},
	)
}
