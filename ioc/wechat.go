package ioc

import (
	"github.com/liumkssq/webook/internal/service/oauth2/wechat"
	"github.com/liumkssq/webook/pkg/logger"
	"os"
)

func InitWechatService(logger logger.LoggerV1) wechat.Service {
	return nil
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("没有找到环境变量 WECHAT_APP_ID ")
	}
	appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("没有找到环境变量 WECHAT_APP_SECRET")
	}
	return wechat.NewService(appId, appKey, logger)
}
