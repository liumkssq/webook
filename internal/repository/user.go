package repository

import (
	"context"
	"database/sql"
	"github.com/liumkssq/webook/internal/domain"
	"github.com/liumkssq/webook/internal/repository/cache"
	"github.com/liumkssq/webook/internal/repository/dao"
	"go.uber.org/zap"
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
	cache      cache.UserCache
	testSignal chan struct{}
}

func NewCachedUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &CachedUserRepository{dao: dao, cache: cache}
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

func (r *CachedUserRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			String: u.Email,
			// 我确实有手机号
			Valid: u.Email != "",
		},
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Password: u.PassWord,
		//WechatOpenID: sql.NullString{
		//	String: u.WechatInfo.OpenID,
		//	Valid:  u.WechatInfo.OpenID != "",
		//},
		//WechatUnionID: sql.NullString{
		//	String: u.WechatInfo.UnionID,
		//	Valid:  u.WechatInfo.UnionID != "",
		//},
		Ctime: u.Ctime.UnixMilli(),
	}
}

func (r *CachedUserRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, nil
	}
	return r.entityToDomain(u), nil
}

func (r *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
}

func (r *CachedUserRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id)
	if err == nil {
		return u, nil
	}
	ue, err := r.dao.FindById(ctx, id)
	if err != nil {
		zap.L().Info("根据 id 查询用户失败", zap.Error(err))
		return domain.User{}, err
	}
	u = r.entityToDomain(ue)
	//异步写入缓存
	go func() {
		_ = r.cache.Set(ctx, u)
		r.testSignal <- struct{}{}
	}()
	return u, nil
}
