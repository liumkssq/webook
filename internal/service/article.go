package service

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/liumkssq/webook/internal/domain"
	"github.com/liumkssq/webook/internal/repository/article"
	"github.com/liumkssq/webook/pkg/logger"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx *gin.Context, d domain.Article) error
	List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error)
	GetById(ctx context.Context, id int64) (domain.Article, error)
	GetPublishedById(ctx context.Context, id int64) (domain.Article, error)
}
type articleService struct {
	repo article.ArticleRepository

	//v1
	author article.ArticleAuthorRepository
	reader article.ArticleReaderRepository
	l      logger.LoggerV1
}

func (a *articleService) List(ctx context.Context, uid int64, offset int, limit int) ([]domain.Article, error) {
	return a.repo.List(ctx, uid, offset, limit)
}

func (a *articleService) GetById(ctx context.Context, id int64) (domain.Article, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *articleService) GetPublishedById(ctx context.Context, id int64) (domain.Article, error) {
	return a.repo.GetPublishedById(ctx, id)
}

func (a *articleService) Withdraw(ctx *gin.Context, art domain.Article) error {
	return a.repo.SyncStatus(ctx, art.Id, art.Author.Id, domain.ArticleStatusPrivate)
}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	id, err := a.repo.Sync(ctx, art)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (a *articleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	//todo
	return 0, nil
}

func NewArticleServiceV1(repo article.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func NewArticleService(author article.ArticleAuthorRepository,
	reader article.ArticleReaderRepository,
	l logger.LoggerV1) ArticleService {
	return &articleService{
		author: author,
		reader: reader,
		l:      l,
	}
}

func (a *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.Id > 0 {
		err := a.repo.Update(ctx, art)
		return art.Id, err
	}
	return a.repo.Create(ctx, art)
}

// Update
func (a *articleService) Update(ctx context.Context, art domain.Article) error {
	//artInDB := a.repo.FindById(art.Id)
	//if art.Author.Id != artInDB.Author.Id {
	//	return errors.New("无权限修改")
	//}
	return a.repo.Update(ctx, art)
}
