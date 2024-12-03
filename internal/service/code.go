package service

import (
	"context"
	"fmt"
	"github.com/liumkssq/webook/internal/repository"
	"github.com/liumkssq/webook/internal/service/sms"
	"go.uber.org/atomic"
	"golang.org/x/exp/rand"
)

// 短信模板
var codeTplId atomic.String = atomic.String{}

var (
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
)

type CodeService interface {
	Send(ctx context.Context,
		// 区别业务场景
		biz string, phone string) error
	Verify(ctx context.Context, biz string,
		phone string, inputCode string) (bool, error)
}
type codeService struct {
	repo   repository.CodeRepository
	smsSvc sms.Service
	//tplId string
}

func (svc *codeService) Send(ctx context.Context, biz string, phone string) error {
	code := svc.generateCode()
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	err = svc.smsSvc.Send(ctx, codeTplId.Load(), []string{code}, phone)
	if err != nil {
		err = fmt.Errorf("短信服务不可用: %w", err)
	}
	return err
}

func (svc *codeService) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func (svc *codeService) generateCode() string {
	// 六位数，num 在 0, 999999 之间，包含 0 和 999999
	num := rand.Intn(1000000)
	// 不够六位的，加上前导 0
	// 000001
	return fmt.Sprintf("%06d", num)
}

func NewCodeService(repo repository.CodeRepository, smsSvc sms.Service) CodeService {
	// 类似Java thread-local
	codeTplId.Store("191919")
	return &codeService{repo: repo, smsSvc: smsSvc}
}
