package ioc

import (
	mysms "github.com/liumkssq/webook/internal/service/sms"
	"github.com/liumkssq/webook/internal/service/sms/memory"
	"github.com/liumkssq/webook/internal/service/sms/ratelimit"
	"github.com/liumkssq/webook/internal/service/sms/retryable"
	"github.com/liumkssq/webook/pkg/limiter"
	"github.com/redis/go-redis/v9"
	"time"
)

//	func initAliYunSMS() *aliyun.Service {
//		secretId, ok := os.LookupEnv("ALIBABA_CLOUD_ACCESS_KEY_ID")
//		if !ok {
//			zap.L().Fatal("ALIBABA_CLOUD_ACCESS_KEY_ID not found")
//		}
//		secretKey, ok := os.LookupEnv("ALIBABA_CLOUD_ACCESS_KEY_SECRET")
//		if !ok {
//			zap.L().Fatal("ALIBABA_CLOUD_ACCESS_KEY_SECRET not found")
//		}
//
//		config := &openapi.Config{
//			AccessKeyId:     ekit.ToPtr[string](secretId),
//			AccessKeySecret: ekit.ToPtr[string](secretKey),
//		}
//		config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
//		client, err := sms.NewClient(config)
//		if err != nil {
//			zap.L().Fatal("阿里云短信服务初始化失败", zap.Error(err))
//		}
//		if client == nil {
//			panic("client is nil")
//		}
//		service := aliyun.NewService(client)
//		return service
//	}
func InitSMSService(cmd redis.Cmdable) mysms.Service {
	// 换内存，还是换别的
	svc := ratelimit.NewRatelimitSMSService(memory.NewService(),
		limiter.NewRedisSlidingWindowLimiter(cmd, time.Second, 100))
	return retryable.NewService(svc, 3)
	//return memory.NewService()
}
