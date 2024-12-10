package startup

import (
	"github.com/liumkssq/webook/internal/service/oauth2/wechat"
	"github.com/liumkssq/webook/pkg/logger"
)

// InitPhantomWechatService 没啥用的虚拟的 wechatService
func InitPhantomWechatService(l logger.LoggerV1) wechat.Service {
	return wechat.NewService("", "", l)
}
