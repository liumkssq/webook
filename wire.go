//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/liumkssq/webook/internal/repository"
	"github.com/liumkssq/webook/internal/repository/dao"
	"github.com/liumkssq/webook/internal/service"
	"github.com/liumkssq/webook/internal/web"
	"github.com/liumkssq/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		ioc.InitDB,
		//dao
		dao.NewGORMUserDAO,

		repository.NewCachedUserRepositoryV1,

		//svc
		service.NewUserService,

		//Handle
		web.NewUserHandler,

		//todo
		ioc.InitWebServer,
		ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
