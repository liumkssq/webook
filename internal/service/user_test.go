package service

import (
	"context"
	"errors"
	"github.com/liumkssq/webook/internal/domain"
	"github.com/liumkssq/webook/internal/repository"
	repomocks "github.com/liumkssq/webook/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func Test_userService_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) repository.UserRepository

		//ctx      context.Context
		email    string
		password string

		wantUser domain.User
		wantErr  error
	}{
		{
			name:     "登录成功",
			email:    "test@example.com",
			password: "hello#world123",

			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "test@example.com").
					Return(domain.User{
						Email:    "test@example.com",
						Phone:    "12345678901",
						PassWord: "$2a$10$dIJMY9AmPd0qVegZgZ6YIebtkKICRnq.kOhVff.DY0iG5cb59oqJG",
						Ctime:    now,
					}, nil)
				return repo
			},
			wantUser: domain.User{
				Email:    "test@example.com",
				Phone:    "12345678901",
				PassWord: "$2a$10$dIJMY9AmPd0qVegZgZ6YIebtkKICRnq.kOhVff.DY0iG5cb59oqJG",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name:     "用户不存在",
			email:    "test@example.com",
			password: "hello#world123",

			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "test@example.com").
					Return(domain.User{}, repository.ErrUserNotFound)
				return repo
			},
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
		{
			name:     "系统错误",
			email:    "test@example.com",
			password: "hello#world123",

			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "test@example.com").
					Return(domain.User{}, errors.New("数据库错误"))
				return repo
			},
			wantUser: domain.User{},
			wantErr:  errors.New("数据库错误"),
		},
		{
			name:     "密码错误",
			email:    "test@example.com",
			password: "hello12#world123",

			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "test@example.com").
					Return(domain.User{
						Email:    "test@example.com",
						Phone:    "12345678901",
						PassWord: "$2a$10$dIJMY9AmPd0qVegZgZ6YIebtkKICRnq.kOhVff.DY0iG5cb59oqJG",
						Ctime:    now,
					}, nil)
				return repo
			},
			wantUser: domain.User{},
			wantErr:  ErrInvalidUserOrPassword,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			svc := NewUserService(tc.mock(ctrl), nil)
			u, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantUser, u)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestEncryptPassword(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("hello#world123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(res))
}
