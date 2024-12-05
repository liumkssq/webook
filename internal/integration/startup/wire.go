//go:build wireinject
// +build wireinject

package startup

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
