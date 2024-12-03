package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/liumkssq/webook/internal/domain"
	"github.com/liumkssq/webook/internal/service"
	ijwt "github.com/liumkssq/webook/internal/web/jwt"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net/http"
)

const biz = "login"

// UserHandler 我准备在它上面定义跟用户有关的路由
type UserHandler struct {
	svc         service.UserService
	codeSvc     service.CodeService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
	phoneExp    *regexp.Regexp
	ijwt.Handler
	cmd redis.Cmdable
}

func NewUserHandler(svc service.UserService,
	codeSvc service.CodeService,
	jwtHdl ijwt.Handler) *UserHandler {
	const (
		emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
		phoneRegexPattern    = `^1[3456789]\d{9}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	phoneExp := regexp.MustCompile(phoneRegexPattern, regexp.None)
	return &UserHandler{
		svc:         svc,
		emailExp:    emailExp,
		passwordExp: passwordExp,
		Handler:     jwtHdl,
		phoneExp:    phoneExp,
		codeSvc:     codeSvc,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.GET("/profile", u.ProfileJWT)
	ug.POST("/signup", u.SignUp)
	ug.POST("/login", u.LoginJWT)
	ug.POST("/logout", u.LogoutJWT)
	ug.POST("/edit", u.Edit)
	ug.POST("/login_sms/code/send", u.SendLoginSMSCode)
	ug.POST("/login_sms", u.LoginSMS)
	ug.POST("/refresh_token", u.RefreshToken)
}

func (u *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if err = u.SetLoginToken(ctx, user.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "登录成功")
	return
}

// todo 可以跳转到登录页面
func (u *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		ConfirmPassword string `json:"confirmPassword"`
		Password        string `json:"password"`
	}
	var req SignUpReq
	// Bind 方法会根据 Content-Type 来解析你的数据到 req 里面
	// 解析错了，就会直接写回一个 400 的错误
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		zap.L().Error("邮箱匹配错误",
			zap.Error(err))
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		zap.L().Error("邮箱格式错误",
			zap.String("email", req.Email))
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}
	if req.ConfirmPassword != req.Password {
		zap.L().Error("密码不匹配",
			zap.String("email", req.Email))
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}
	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		// 记录日志
		zap.L().Error("密码正则匹配错误",
			zap.Error(err))
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字、特殊字符")
		return
	}

	// 调用一下 svc 的方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		PassWord: req.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
}

func (u *UserHandler) ProfileJWT(ctx *gin.Context) {
	c, _ := ctx.Get("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		zap.L().Error("claims 类型断言失败")
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	user, err := u.svc.Profile(ctx, claims.Id)
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: user,
	})
}

func (u *UserHandler) Edit(ctx *gin.Context) {
	//todo
	//type Req struct {
	//	Id      int64
	//	NickName string `json:"nickName"`
	//}
	//var req Req
	//if err := ctx.Bind(&req); err != nil {
	//	return
	//}
	//uc := ctx.MustGet("user").(ijwt.UserClaims)
	//err := u.svc.Edit(ctx, domain.User{
	//	Id:       req.Id,
	//	NickName: req.NickName,
	//})
	//if err != nil {
	//	ctx.JSON(http.StatusOK, Result{
	//		Code: 5,
	//		Msg: "编辑失败",
	//	})
	//	return
	//}
}

func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ok, err := u.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		zap.L().Error("校验验证码出错", zap.Error(err),
			// 不能这样打，因为手机号码是敏感数据，你不能达到日志里面
			// 打印加密后的串
			// 脱敏，152****1234
			zap.String("手机号码", req.Phone))
		// 最多最多就这样
		zap.L().Debug("", zap.String("手机号码", req.Phone))
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码有误",
		})
		return
	}
	// 我这个手机号，会不会是一个新用户呢？
	// 这样子
	user, err := u.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	if err = u.SetLoginToken(ctx, user.Id); err != nil {
		// 记录日志
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}

	ctx.JSON(http.StatusOK, Result{
		Msg: "验证码校验通过",
	})

}

func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 是不是一个合法的手机号码
	// 考虑正则表达式
	ok, err := u.phoneExp.MatchString(req.Phone)
	if err != nil {
		zap.L().Error("系统错误", zap.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "手机号码格式不对",
		})
		return
	}
	err = u.codeSvc.Send(ctx, biz, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		zap.L().Warn("短信发送太频繁",
			zap.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送太频繁，请稍后再试",
		})
	default:
		zap.L().Error("短信发送失败",
			zap.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}

// Refresh Token
func (u *UserHandler) RefreshToken(ctx *gin.Context) {
	refreshToken := u.ExtractToken(ctx)
	var rc ijwt.RefreshClaims
	token, err := jwt.ParseWithClaims(refreshToken, &rc, func(token *jwt.Token) (interface{}, error) {
		return ijwt.RtKey, nil
	})
	if err != nil || !token.Valid {
		zap.L().Error("系统错误", zap.Error(err))
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = u.CheckSession(ctx, rc.Ssid)
	if err != nil {
		zap.L().Error("系统错误", zap.Error(err))
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	err = u.SetJWTToken(ctx, rc.Uid, rc.Ssid)
	if err != nil {
		zap.L().Error("系统错误", zap.Error(err))
		ctx.AbortWithStatus(http.StatusUnauthorized)
		zap.L().Error("JWT 续约失败",
			zap.Error(err),
			zap.String("token", refreshToken))
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (u *UserHandler) LogoutJWT(ctx *gin.Context) {
	err := u.ClearToken(ctx)
	if err != nil {
		zap.L().Error("系统错误", zap.Error(err))
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "退出登陆失败",
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "退出登陆成功",
	})
}
