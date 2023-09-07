package utils

import (
	"context"
	"github.com/gin-gonic/gin"
)

func GetGinCtx(ctx *gin.Context) context.Context {
	return ctx.MustGet("ctx").(context.Context)
}
