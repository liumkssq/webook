package ratelimit

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/liumkssq/webook/pkg/ratelimit"
)

type Builder struct {
	prefix  string
	limiter ratelimit.Limiter
}

func NewBuilder(prefix string, limiter ratelimit.Limiter) *Builder {
	return &Builder{
		prefix:  prefix,
		limiter: limiter,
	}
}
func (b *Builder) Prefix(prefix string) *Builder {
	b.prefix = prefix
	return b
}

//

func (b *Builder) limit(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP())
	return b.limiter.Limit(ctx, key)
}
