//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/liumkssq/webook/internal/events/article"
	"github.com/liumkssq/webook/internal/repository"
	article2 "github.com/liumkssq/webook/internal/repository/article"
	"github.com/liumkssq/webook/internal/repository/cache"
	"github.com/liumkssq/webook/internal/repository/dao"
	"github.com/liumkssq/webook/internal/service"
	"github.com/liumkssq/webook/internal/web"
	ijwt "github.com/liumkssq/webook/internal/web/jwt"
	"github.com/liumkssq/webook/ioc"
)

func InitWebServer() *App {
	wire.Build(
		ioc.InitLogger,
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitKafka,
		ioc.NewConsumers,
		ioc.NewSyncProducer,

		article.NewKafkaConsumer,
		//dao
		dao.NewGORMInteractiveDAO,
		//dao.NewGORMInteractiveDAO,
		dao.NewUserDAO,
		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,
		cache.NewRedisInteractiveCache,

		article2.NewArticleRepository,
		repository.NewCachedUserRepository,
		repository.NewCachedCodeRepository,
		repository.NewCachedInteractiveRepository,

		//svc
		service.NewArticleServiceV1,
		service.NewInteractiveService,
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
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
