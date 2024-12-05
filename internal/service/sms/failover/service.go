package failover

import (
	"context"
	"errors"
	"github.com/liumkssq/webook/internal/service/sms"
	"go.uber.org/zap"
	"sync/atomic"
)

type FailoverSMSService struct {
	// 多个服务
	svcs []sms.Service
	idx  uint64
}

func NewFailoverSMSService(svcs []sms.Service) sms.Service {
	return &FailoverSMSService{
		svcs: svcs,
	}
}

func (f *FailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	for _, svc := range f.svcs {
		err := svc.Send(ctx, tpl, args, numbers...)
		// 发送成功
		if err == nil {
			return nil
		}
		// 正常这边，输出日志
		// 要做好监控
		zap.L().Error("短信发送失败",
			zap.Error(err),
			zap.String("tpl", tpl),
			zap.Strings("args", args),
			zap.Strings("numbers", numbers))
	}
	return errors.New("全部服务商都失败了")
}

func (f *FailoverSMSService) SendV1(ctx context.Context, tpl string, args []string, numbers ...string) error {
	// 我取下一个节点来作为起始节点
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		svc := f.svcs[int(i%length)]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled:
			return err
		default:
			// 输出日志
		}
	}
	return errors.New("全部服务商都失败了")
}
