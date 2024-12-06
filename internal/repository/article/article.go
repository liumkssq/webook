package article

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	"github.com/liumkssq/webook/internal/domain"
	"github.com/liumkssq/webook/internal/repository"
	"github.com/liumkssq/webook/internal/repository/cache"
	dao "github.com/liumkssq/webook/internal/repository/dao/article"
	"github.com/liumkssq/webook/pkg/logger"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	FindById(id int64) domain.Article
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx *gin.Context, id int64, author int64, status domain.ArticleStatus) error
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetByID(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id int64) (domain.Article, error)
}

type CachedArticleRepository struct {
	dao      dao.ArticleDAO
	userRepo repository.UserRepository

	readerDAO dao.ReaderDAO
	authorDAO dao.AuthorDAO

	cache cache.ArticleCache
	l     logger.LoggerV1
}

func (c *CachedArticleRepository) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	//缓存
	if offset == 0 && limit <= 100 {
		data, err := c.cache.GetFirstPage(ctx, uid)
		if err == nil {
			go func() {
				c.preCache(ctx, data)
			}()
			return data, err
		}
	}
	res, err := c.dao.GetByAuthor(ctx, uid, offset, limit)
	if err != nil {
		return nil, err
	}
	//利用go 语言
	data := slice.Map[dao.Article, domain.Article](res, func(idx int, src dao.Article) domain.Article {
		return c.toDomain(src)
	})

	go func() {
		err := c.cache.SetFirstPage(ctx, uid, data)
		c.l.Error("缓存失败", logger.Error(err))
		c.preCache(ctx, data)
	}()
	return data, err
}

func (c *CachedArticleRepository) preCache(ctx context.Context, data []domain.Article) {
	if len(data) > 0 && len(data[0].Content) < 1024*1024 {
		err := c.cache.Set(ctx, data[0])
		if err != nil {
			c.l.Error("提前预加载缓存失败", logger.Error(err))
		}
	}
}

func (c *CachedArticleRepository) GetByID(ctx context.Context, id int64) (domain.Article, error) {
	data, err := c.dao.GetById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return c.toDomain(data), err
}

func (c *CachedArticleRepository) GetPublishedById(ctx context.Context, id int64) (domain.Article, error) {
	art, err := c.dao.GetPubById(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	user, err := c.userRepo.FindById(ctx, art.AuthorId)
	res := domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id:   art.AuthorId,
			Name: user.NickName,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
	return res, nil
}

func (c *CachedArticleRepository) SyncStatus(ctx *gin.Context, id int64, author int64, status domain.ArticleStatus) error {
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
func (repo *CachedArticleRepository) toDomain(art dao.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
}
