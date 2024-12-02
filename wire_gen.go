// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/liumkssq/webook/internal/repository"
	"github.com/liumkssq/webook/internal/repository/dao"
	"github.com/liumkssq/webook/internal/service"
	"github.com/liumkssq/webook/internal/web"
	"github.com/liumkssq/webook/ioc"
)

// Injectors from wire.go:

func InitWebServer() *gin.Engine {
	v := ioc.InitMiddlewares()
	db := ioc.InitDB()
	userDAO := dao.NewGORMUserDAO(db)
	userRepository := repository.NewCachedUserRepositoryV1(userDAO)
	userService := service.NewUserService(userRepository)
	userHandler := web.NewUserHandler(userService)
	engine := ioc.InitWebServer(v, userHandler)
	return engine
}
