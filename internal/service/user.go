package service

import (
	"context"
	"errors"
	"github.com/liumkssq/webook/internal/domain"
	"github.com/liumkssq/webook/internal/repository"
	"github.com/liumkssq/webook/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicate
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

type UserService interface {
	Login(ctx context.Context, email, password string) (domain.User, error)
	SignUp(ctx context.Context, u domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
	l    logger.LoggerV1
}

func (svc *userService) FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error) {
	u, err := svc.repo.FindByWechat(ctx, wechatInfo.OpenID)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	u = domain.User{
		WechatInfo: wechatInfo,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return u, err
	}
	return svc.repo.FindByWechat(ctx, wechatInfo.OpenID)
}

func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PassWord), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, nil
}

func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.PassWord), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PassWord = string(hash)
	// 然后就是，存起来
	return svc.repo.Create(ctx, u)
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFound {
		return u, err
	}
	svc.l.Info("用户不存在，自动创建", logger.String("phone", phone))
	u = domain.User{
		Phone: phone,
	}
	err = svc.repo.Create(ctx, u)
	if err != nil && err != repository.ErrUserDuplicate {
		return domain.User{}, err
	}
	//主从延迟
	return svc.repo.FindByPhone(ctx, phone)
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}

func NewUserService(repo repository.UserRepository, l logger.LoggerV1) UserService {
	return &userService{
		repo: repo,
		l:    l,
	}
}
