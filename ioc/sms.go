package ioc

import (
	"github.com/liumkssq/webook/internal/service/sms"
	"github.com/liumkssq/webook/internal/service/sms/memory"
	"github.com/liumkssq/webook/internal/service/sms/ratelimit"
	"github.com/liumkssq/webook/internal/service/sms/retryable"
	"github.com/liumkssq/webook/pkg/limiter"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitSMSService(cmd redis.Cmdable) sms.Service {
	// 换内存，还是换别的
	svc := ratelimit.NewRatelimitSMSService(memory.NewService(),
		limiter.NewRedisSlidingWindowLimiter(cmd, time.Second, 100))
	return retryable.NewService(svc, 3)
	//return memory.NewService()
}
