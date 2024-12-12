//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/liumkssq/webook/interactive/events"
	"github.com/liumkssq/webook/interactive/grpc"
	"github.com/liumkssq/webook/interactive/ioc"
	repository2 "github.com/liumkssq/webook/interactive/repository"
	cache2 "github.com/liumkssq/webook/interactive/repository/cache"
	dao2 "github.com/liumkssq/webook/interactive/repository/dao"
	service2 "github.com/liumkssq/webook/interactive/service"
)

var thirdPartySet = wire.NewSet(ioc.InitSrcDB,
	ioc.InitDstDB,
	ioc.InitDoubleWritePool,
	ioc.InitBizDB,
	ioc.InitLogger,
	ioc.InitSaramaClient,
	ioc.InitSaramaSyncProducer,
	ioc.InitRedis)

var interactiveSvcSet = wire.NewSet(dao2.NewGORMInteractiveDAO,
	cache2.NewInteractiveRedisCache,
	repository2.NewCachedInteractiveRepository,
	events.NewInteractiveProducer,
	service2.NewInteractiveService,
)

func InitApp() *App {
	wire.Build(thirdPartySet,
		interactiveSvcSet,
		grpc.NewInteractiveServiceServer,
		events.NewInteractiveReadEventConsumer,
		ioc.InitInteractiveProducer,
		ioc.InitFixerConsumer,
		ioc.InitConsumers,
		ioc.NewGrpcxServer,
		ioc.InitGinxServer,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
