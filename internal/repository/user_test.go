package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/liumkssq/webook/internal/domain"
	"github.com/liumkssq/webook/internal/repository/cache"
	cachemocks "github.com/liumkssq/webook/internal/repository/cache/mocks"
	"github.com/liumkssq/webook/internal/repository/dao"
	daomocks "github.com/liumkssq/webook/internal/repository/dao/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	now := time.Now()
	now = time.UnixMilli(now.UnixMilli())
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)

		id       int64
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "缓存命中,查询成功",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				//
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)
				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(dao.User{
						Id: 123,
						Email: sql.NullString{
							String: "123@qq.com",
							Valid:  true,
						},
						Password: "123456",
						Phone: sql.NullString{
							String: "15222222222",
							Valid:  true,
						},
						Ctime: now.UnixMilli(),
						Utime: now.UnixMilli(),
					}, nil)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					PassWord: "123456",
					Phone:    "15222222222",
					Ctime:    now,
				}).Return(nil)
				return d, c
			},
			id: 123,
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				PassWord: "123456",
				Phone:    "15222222222",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				//
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{
						Id:       123,
						Email:    "123@qq.com",
						PassWord: "123456",
						Phone:    "15222222222",
						Ctime:    now,
					}, nil)
				d := daomocks.NewMockUserDAO(ctrl)
				return d, c
			},
			id: 123,
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				PassWord: "123456",
				Phone:    "15222222222",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "查询失败",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				//
				c := cachemocks.NewMockUserCache(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).
					Return(domain.User{}, cache.ErrKeyNotExist)
				d := daomocks.NewMockUserDAO(ctrl)
				d.EXPECT().FindById(gomock.Any(), int64(123)).
					Return(dao.User{}, errors.New("数据库查询失败"))
				return d, c
			},
			id:       123,
			wantUser: domain.User{},
			wantErr:  errors.New("数据库查询失败"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ud, uc := tc.mock(ctrl)
			repo := NewCachedUserRepository(ud, uc)
			u, err := repo.FindById(context.Background(), tc.id)
			assert.Equal(t, tc.wantUser, u)
			assert.Equal(t, tc.wantErr, err)
			time.Sleep(1 * time.Second)
		})
	}
}
