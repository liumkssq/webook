package service

import (
	"context"
	"errors"
	"github.com/liumkssq/webook/internal/domain"
	"github.com/liumkssq/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserDuplicateEmail = repository.ErrUserDuplicate
var ErrInvalidUserOrPassword = errors.New("账号/邮箱或密码不对")

type UserService interface {
	Login(ctx context.Context, email, password string) (domain.User, error)
	SignUp(ctx context.Context, u domain.User) error
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	//todo
	//FindOrCreateByWechat(ctx context.Context, wechatInfo domain.WechatInfo) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
}

type userService struct {
	repo repository.UserRepository
	//	todo log
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
	//todo
	panic("")
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}
