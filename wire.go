//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liumkssq/webook/internal/repository"
	article2 "github.com/liumkssq/webook/internal/repository/article"
	"github.com/liumkssq/webook/internal/repository/cache"
	"github.com/liumkssq/webook/internal/repository/dao"
	artdao "github.com/liumkssq/webook/internal/repository/dao/article"
	"github.com/liumkssq/webook/internal/service"
	"github.com/liumkssq/webook/internal/web"
	ijwt "github.com/liumkssq/webook/internal/web/jwt"
	"github.com/liumkssq/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitLogger,
		ioc.InitDB,
		ioc.InitRedis,
		//dao
		artdao.NewGORMArticleDAO,
		dao.NewGORMUserDAO,
		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,

		article2.NewArticleRepository,
		repository.NewCachedUserRepository,
		repository.NewCachedCodeRepository,

		//svc
		service.NewArticleServiceV1,
		service.NewUserService,
		service.NewCodeService,

		ioc.InitSMSService,
		ioc.InitWechatService,

		//Handle
		web.NewArticleHandler,
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		ijwt.NewRedisJWTHandler,

		//todo
		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
