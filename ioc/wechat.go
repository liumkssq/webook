package ioc

import (
	"github.com/liumkssq/webook/internal/service/oauth2/wechat"
	"github.com/liumkssq/webook/pkg/logger"
	"go.uber.org/zap"
	"os"
)

func InitWechatService(l logger.LoggerV1) wechat.Service {
	//todo
	//return nil
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		zap.L().Fatal("没有找到环境变量 WECHAT_APP_ID")
		panic("没有找到环境变量 WECHAT_APP_ID ")
	}
	appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		zap.L().Fatal("没有找到环境变量 WECHAT_APP_SECRET")
		panic("没有找到环境变量 WECHAT_APP_SECRET")
	}
	// 692jdHsogrsYqxaUK9fgxw
	return wechat.NewService(appId, appKey, nil)
}
