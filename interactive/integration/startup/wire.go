//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"github.com/liumkssq/webook/interactive/grpc"
	repository2 "github.com/liumkssq/webook/interactive/repository"
	cache2 "github.com/liumkssq/webook/interactive/repository/cache"
	dao2 "github.com/liumkssq/webook/interactive/repository/dao"
	service2 "github.com/liumkssq/webook/interactive/service"
)

var thirdPartySet = wire.NewSet( // 第三方依赖
	InitRedis, InitDB,
	//InitSaramaClient,
	//InitSyncProducer,
	InitLogger,
)

var interactiveSvcSet = wire.NewSet(dao2.NewGORMInteractiveDAO,
	//events2.NewInteractiveProducer,
	cache2.NewInteractiveRedisCache,
	repository2.NewCachedInteractiveRepository,
	service2.NewInteractiveService,
)

func InitInteractiveService() *grpc.InteractiveServiceServer {
	wire.Build(thirdPartySet, interactiveSvcSet, grpc.NewInteractiveServiceServer)
	return new(grpc.InteractiveServiceServer)
}