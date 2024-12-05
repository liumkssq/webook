package service

import (
	"context"
	"github.com/liumkssq/webook/internal/domain"
	"github.com/liumkssq/webook/internal/repository/article"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
}
type articleService struct {
	repo article.ArticleRepository

	//v1
	author article.ArticleAuthorRepository
	reader article.ArticleReaderRepository
}

func (a *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	id, err := a.repo.Create(ctx, art)

}

func (a *articleService) PublishV1(ctx context.Context, art domain.Article) (int64, error) {
	id, err := a.repo.Create(ctx, art)

}

func NewArticleService(repo article.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
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
