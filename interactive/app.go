package main

import (
	"github.com/liumkssq/webook/internal/events"
	"github.com/liumkssq/webook/pkg/ginx"
	"github.com/liumkssq/webook/pkg/grpcx"
)

type App struct {
	consumers   []events.Consumer
	server      *grpcx.Server
	adminServer *ginx.Server
}
