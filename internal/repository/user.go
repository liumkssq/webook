package repository

import (
	"context"
	"github.com/liumkssq/webook/internal/domain"
	"github.com/liumkssq/webook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicate = dao.ErrUserDuplicate
	ErrUserNotFound  = dao.ErrUserNotFound
)

// UserRepository 是核心，它有不同实现。
// 但是 Factory 本身如果只是初始化一下，那么它不是你的核心
type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
	Create(ctx context.Context, u domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
	//todo
	//FindByWechat(ctx context.Context, openID string) (domain.User, error)
}

type CachedUserRepository struct {
	dao dao.UserDAO
	//todo
	//cache cache.UserCache
	//testSignal chan struct{}
}

func (r *CachedUserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, nil
	}
	return r.entityToDomain(u), nil
}

// utils
func (r *CachedUserRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		PassWord: u.Password,
		Phone:    u.Phone.String,
		//todo
		//WechatInfo: domain.WechatInfo{
		//	UnionID: u.WechatUnionID.String,
		//	OpenID:  u.WechatOpenID.String,
		//},
		Ctime: time.UnixMilli(u.Ctime),
	}
}

func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	//TODO implement me
	panic("implement me")
}

func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	//TODO implement me
	panic("implement me")
}

func NewCachedUserRepositoryV1(dao dao.UserDAO) UserRepository {
	return &CachedUserRepository{dao: dao}
}
