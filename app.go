package main

import (
	"github.com/gin-gonic/gin"
	"github.com/liumkssq/webook/internal/events"
)

type App struct {
	server   *gin.Engine
	consumer []events.Consumer
}
