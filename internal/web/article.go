package web

import (
	"github.com/gin-gonic/gin"
	"github.com/liumkssq/webook/internal/domain"
	"github.com/liumkssq/webook/internal/service"
	ijwt "github.com/liumkssq/webook/internal/web/jwt"
	"github.com/liumkssq/webook/pkg/logger"
	"net/http"
)

type ArticleHandler struct {
	svc service.ArticleService
	l   logger.LoggerV1
}

func NewArticleHandler(svc service.ArticleService, l logger.LoggerV1) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
		l:   l,
	}
}

func (h *ArticleHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/articles")

	g.POST("/edit", h.Edit)
	g.POST("/withdraw", h.Withdraw)
	g.POST("/publish", h.Publish)
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("获取用户信息出错")
		return
	}
	id, err := h.svc.Save(ctx, req.toDomain(claims.Id))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("保存出错", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})
}

func (h *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("获取用户信息出错")
		return
	}
	id, err := h.svc.Publish(ctx, req.toDomain(claims.Id))
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("发布出错", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: id,
	})
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	c := ctx.MustGet("claims")
	claims, ok := c.(*ijwt.UserClaims)
	if !ok {
		// 你可以考虑监控住这里
		//ctx.AbortWithStatus(http.StatusUnauthorized)
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		h.l.Error("未发现用户的 session 信息")
		return
	}

	err := h.svc.Withdraw(ctx, domain.Article{
		Id: req.Id,
		Author: domain.Author{
			Id: claims.Uid,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 打日志？
		h.l.Error("保存帖子失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

type ArticleReq struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func (receiver ArticleReq) toDomain(uid int64) domain.Article {
	return domain.Article{
		Title:   receiver.Title,
		Content: receiver.Content,
		Author: domain.Author{
			Id: uid,
		},
	}
}
