package web

import (
	"bytes"
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/liumkssq/webook/internal/domain"
	"github.com/liumkssq/webook/internal/repository/cache/redismocks"
	"github.com/liumkssq/webook/internal/service"
	svcmocks "github.com/liumkssq/webook/internal/service/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestNewUserHandler(t *testing.T) {
	type args struct {
		svc service.UserService
	}
	tests := []struct {
		name string
		args args
		want *UserHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserHandler(tt.args.svc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserHandler_RegisterRoutes(t *testing.T) {
	type fields struct {
		svc         service.UserService
		emailExp    *regexp.Regexp
		passwordExp *regexp.Regexp
	}
	type args struct {
		server *gin.Engine
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UserHandler{
				svc:         tt.fields.svc,
				emailExp:    tt.fields.emailExp,
				passwordExp: tt.fields.passwordExp,
			}
			u.RegisterRoutes(tt.args.server)
		})
	}
}

func TestUserHandler_SignUp(t *testing.T) {
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					PassWord: "hello#world123",
				}).Return(nil)
				// 注册成功是 return nil
				return usersvc
			},

			reqBody: `
{
	"email": "123@qq.com",
	"password": "hello#world123",
	"confirmPassword": "hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "注册成功",
		},
		{
			name: "参数不对，bind 失败",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				// 注册成功是 return nil
				return usersvc
			},

			reqBody: `
{
	"email": "123@qq.com",
	"password": "hello#world123"
`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},

			reqBody: `
{
	"email": "123@q",
	"password": "hello#world123",
	"confirmPassword": "hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "你的邮箱格式不对",
		},
		{
			name: "两次输入密码不匹配",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},

			reqBody: `
{
	"email": "123@qq.com",
	"password": "hello#world1234",
	"confirmPassword": "hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "两次输入的密码不一致",
		},
		{
			name: "密码格式不对",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				return usersvc
			},
			reqBody: `
{
	"email": "123@qq.com",
	"password": "hello123",
	"confirmPassword": "hello123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "密码必须大于8位，包含数字、特殊字符",
		},
		{
			name: "邮箱冲突",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					PassWord: "hello#world123",
				}).Return(service.ErrUserDuplicateEmail)
				// 注册成功是 return nil
				return usersvc
			},

			reqBody: `
{
	"email": "123@qq.com",
	"password": "hello#world123",
	"confirmPassword": "hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "邮箱冲突",
		},
		{
			name: "系统异常",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().SignUp(gomock.Any(), domain.User{
					Email:    "123@qq.com",
					PassWord: "hello#world123",
				}).Return(errors.New("随便一个 error"))
				// 注册成功是 return nil
				return usersvc
			},

			reqBody: `
{
	"email": "123@qq.com",
	"password": "hello#world123",
	"confirmPassword": "hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "系统异常",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			server := gin.Default()
			h := NewUserHandler(tc.mock(ctrl))
			require.NotNil(t, h, "NewUserHandler")
			h.RegisterRoutes(server)
			req, err := http.NewRequest(http.MethodPost,
				"/users/signup", bytes.NewBuffer([]byte(tc.reqBody)))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}

func TestEncrypt(t *testing.T) {
	_ = NewUserHandler(nil)
	password := "hello#world123"
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	err = bcrypt.CompareHashAndPassword(encrypted, []byte(password))
	fmt.Printf("%s", encrypted)
	assert.NoError(t, err)
}

// todo 有点难测
func TestUserHandler_LoginJWT(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tests := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) service.UserService
		reqBody  string
		wantCode int
		wantBody string
	}{
		{
			name: "登录成功",
			mock: func(ctrl *gomock.Controller) service.UserService {
				usersvc := svcmocks.NewMockUserService(ctrl)
				usersvc.EXPECT().Login(gomock.Any(), "test@example.com", "hello#world123").
					Return(domain.User{
						Id:       1,
						Email:    "test@example.com",
						PassWord: "$2a$10$FmGqYN31VG7xVumXoNL0beJqeOk1b0fPwbWTrgK6gsHKa6LJmYvl2-",
					}, nil)
				return usersvc
			},
			reqBody: `
{
    "email": "test@example.com",
    "password": "hello#world123"
}
`,
			wantCode: http.StatusOK,
			wantBody: "登录成功",
		},
		// 其他测试用例...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			usersvc := tt.mock(ctrl)
			// Mock ijwt.Handler
			jwtHandler := &svcmocks.MockHandler{}
			//jwtHandler.EXPECT().GenerateToken(int64(1)).Return("mocked-token", nil)
			jwtHandler.EXPECT().SetLoginToken(gomock.Any(), int64(1)).Return(nil)
			// Mock redis.Cmdable
			redisCmd := redismocks.NewMockCmdable(ctrl)
			redisCmd.EXPECT().Set(gomock.Any(), "user:1", "mocked-token", gomock.Any()).Return(&redis.StatusCmd{}, nil)
			handler := &UserHandler{
				svc:         usersvc,
				emailExp:    regexp.MustCompile("^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$", regexp.None),
				passwordExp: regexp.MustCompile(`^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`, regexp.None),
				Handler:     jwtHandler,
				cmd:         redisCmd,
			}

			server := gin.Default()
			server.POST("/users/login", handler.LoginJWT)

			req, _ := http.NewRequest(http.MethodPost, "/users/login", strings.NewReader(tt.reqBody))
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)

			if resp.Code != tt.wantCode {
				t.Errorf("Expected status code %d, got %d", tt.wantCode, resp.Code)
			}

			if resp.Body.String() != tt.wantBody {
				t.Errorf("Expected body %q, got %q", tt.wantBody, resp.Body.String())
			}
		})
	}
}
