package integration

import (
	"github.com/gin-gonic/gin"
	article2 "github.com/liumkssq/webook/internal/repository/article"
	"github.com/liumkssq/webook/internal/repository/dao/article"
	"github.com/liumkssq/webook/internal/service"
	"github.com/liumkssq/webook/internal/web"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ArticleTestSuite struct {
	suite.Suite
	server *gin.Engine
}

func (s *ArticleTestSuite) SetupSuite() {
	s.server = gin.Default()
	//todo
	dao := article.NewGORMArticleDAO()
	repo := article2.NewArticleRepository(dao)
	svc := service.NewArticleService(repo)
	artHdl := web.NewArticleHandler(svc, nil)
	artHdl.RegisterRoutes(s.server)
}

func (s *ArticleTestSuite) TestABC() {
	testcases := []struct {
		name     string
		before   func(t *testing.T)
		after    func(t *testing.T)
		art      Article
		wantCode int
		wantRes  Result[int64]
	}{
		{},
	}
	for _, tc := range testcases {
		s.Run(tc.name, func() {
			s.Equal(tc.wantCode, tc.wantRes.Code)
		})
	}
}

func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}

type Article struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type Result[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}
