package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/liumkssq/webook/internal/web"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc,
	userHdl *web.UserHandler) *gin.Engine {

	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)

	return server
}

func InitMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
	}
}

func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowHeaders: []string{"Content-Type", "Authorization"},
		// 允许跨域的域名
		ExposeHeaders: []string{"x-jwt-token", "x-refresh-token"},

		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "yourdomain.com")
		},
		// 缓存时间
		MaxAge: 12 * time.Hour,
	})
}
