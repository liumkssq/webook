package article

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/liumkssq/webook/internal/domain"
	dao "github.com/liumkssq/webook/internal/repository/dao/article"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	FindById(id int64) domain.Article
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx *gin.Context, id int64, author int64, status domain.ArticleStatus) error
}

type CachedArticleRepository struct {
	dao dao.ArticleDAO

	readerDAO dao.ReaderDAO
	authorDAO dao.AuthorDAO
}

func (c *CachedArticleRepository) SyncStatus(ctx *gin.Context, id int64, author int64, status domain.ArticleStatus) error {
	//TODO implement me
	return c.dao.SyncStatus(ctx, id, author, uint8(status))
}

func (c *CachedArticleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	artn := c.toEntity(art)
	if art.Id == 0 {
		id, err = c.authorDAO.Insert(ctx, artn)
	} else {
		err = c.authorDAO.UpdateById(ctx, artn)
	}
	if err != nil {
		return 0, err
	}
	//同步
	//Insert Or Update
	err = c.readerDAO.UpdateOrInsert(ctx, artn)
	return id, err
}

func (c *CachedArticleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return c.dao.Insert(ctx, dao.Article{
		AuthorId: art.Author.Id,
		Content:  art.Content,
		Title:    art.Title,
	})
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &CachedArticleRepository{
		dao: dao,
	}
}

func (c *CachedArticleRepository) Update(ctx context.Context, art domain.Article) error {
	return c.dao.UpdateById(ctx, dao.Article{
		Id:       art.Id,
		AuthorId: art.Author.Id,
		Content:  art.Content,
		Title:    art.Title,
	})
}

func (c *CachedArticleRepository) FindById(id int64) domain.Article {
	return domain.Article{}
}

func (c *CachedArticleRepository) toEntity(art domain.Article) dao.Article {
	return dao.Article{
		AuthorId: art.Author.Id,
		Content:  art.Content,
		Title:    art.Title,
		Id:       art.Id,
	}

}
