//go:build wireinject

package main

import (
	"github.com/google/wire"
	article2 "github.com/liumkssq/webook/internal/events/article"
	"github.com/liumkssq/webook/internal/repository"
	"github.com/liumkssq/webook/internal/repository/cache"
	"github.com/liumkssq/webook/internal/repository/dao"
	"github.com/liumkssq/webook/internal/repository/dao/article"
	"github.com/liumkssq/webook/internal/service"
	"github.com/liumkssq/webook/internal/web"
	ijwt "github.com/liumkssq/webook/internal/web/jwt"
	"github.com/liumkssq/webook/ioc"
)

func InitApp() *App {
	wire.Build(
		ioc.InitRedis, ioc.InitDB,
		ioc.InitLogger,
		ioc.InitKafka,
		ioc.NewSyncProducer,

		// DAO 部分
		dao.NewGORMUserDAO,
		dao.NewGORMInteractiveDAO,
		article.NewGORMArticleDAO,

		// Cache 部分
		cache.NewRedisUserCache,
		cache.NewRedisCodeCache,
		cache.NewRedisArticleCache,
		cache.NewRedisInteractiveCache,

		// repository 部分
		repository.NewCachedUserRepository,
		repository.NewCachedCodeRepository,
		repository.NewArticleRepository,
		repository.NewCachedInteractiveRepository,

		// events 部分
		article2.NewSaramaSyncProducer,
		article2.NewInteractiveReadEventConsumer,
		ioc.NewConsumers,

		// service 部分
		ioc.InitSmsService,
		ioc.InitWechatService,
		service.NewSMSCodeService,
		service.NewUserService,
		service.NewArticleService,
		service.NewInteractiveService,

		// handler 部分
		ijwt.NewRedisHandler,
		web.NewUserHandler,
		web.NewArticleHandler,
		web.NewOAuth2WechatHandler,

		// gin 的中间件
		ioc.GinMiddlewares,

		// Web 服务器
		ioc.InitWebServer,

		wire.Struct(new(App), "*"),
	)
	// 随便返回一个
	return new(App)
}
