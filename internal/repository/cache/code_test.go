package cache

import (
	"context"
	"errors"
	"github.com/liumkssq/webook/internal/repository/cache/redismocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		ctx     context.Context
		biz     string
		phone   string
		code    string
		wantErr error
	}{
		{
			name: "设置成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_code:login:12345678901"},
					[]any{"123456"}).
					Return(res)
				return cmd
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "12345678901",
			code:    "123456",
			wantErr: nil,
		},

		{
			name: "redis错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				//res.SetVal(int64(0))
				res.SetErr(errors.New("系统错误"))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_code:login:12345678901"},
					[]any{"123456"}).
					Return(res)
				return cmd
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "12345678901",
			code:    "123456",
			wantErr: errors.New("系统错误"),
		},

		{
			name: "发送次数过多",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-1))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_code:login:12345678901"},
					[]any{"123456"}).
					Return(res)
				return cmd
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "12345678901",
			code:    "123456",
			wantErr: ErrCodeSendTooMany,
		},

		{
			name: "未知错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				cmd := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-10))
				cmd.EXPECT().Eval(gomock.Any(), luaSetCode,
					[]string{"phone_code:login:12345678901"},
					[]any{"123456"}).
					Return(res)
				return cmd
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "12345678901",
			code:    "123456",
			wantErr: errors.New("系统错误"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewRedisCodeCache(tc.mock(ctrl))
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}
