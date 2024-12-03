//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liumkssq/webook/internal/repository"
	"github.com/liumkssq/webook/internal/repository/cache"
	"github.com/liumkssq/webook/internal/repository/dao"
	"github.com/liumkssq/webook/internal/service"
	"github.com/liumkssq/webook/internal/web"
	ijwt "github.com/liumkssq/webook/internal/web/jwt"
	"github.com/liumkssq/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,
		//dao
		dao.NewGORMUserDAO,
		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,

		repository.NewCachedUserRepository,
		repository.NewCachedCodeRepository,

		//svc
		service.NewUserService,
		service.NewCodeService,

		ioc.InitSMSService,
		ioc.InitWechatService,

		//Handle
		web.NewUserHandler,
		web.NewOAuth2WechatHandler,
		ijwt.NewRedisJWTHandler,

		//todo
		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
