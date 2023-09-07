//go:build debug_log

package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	slog "github.com/jau1jz/cornus/v2/commons/log"
	"github.com/jau1jz/cornus/v2/commons/utils"
	"net/http"
	"runtime"
	"time"
)

func Default(ctx *gin.Context) {
	value := context.WithValue(context.Background(), "trace_id", utils.GenerateUUID())
	ctx.Set("ctx", value)
	defer func() {
		if err := recover(); err != nil {
			var stacktrace string
			for i := 1; ; i++ {
				_, f, l, got := runtime.Caller(i)
				if !got {
					break
				}
				stacktrace += fmt.Sprintf("%s:%d\n", f, l)
			}
			// when stack finishes
			logMessage := fmt.Sprintf("Recovered from a route's Handler('%s')\n", ctx.HandlerName())
			logMessage += fmt.Sprintf("Trace: %s", err)
			logMessage += fmt.Sprintf("\n%s", stacktrace)
			slog.Slog.ErrorF(value, logMessage)
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
	}()
	start := time.Now()
	ctx.Next()
	addr := ctx.Request.RemoteAddr
	if ctx.GetHeader("X-Forwarded-For") != "" {
		addr = ctx.GetHeader("X-Forwarded-For")
	}
	slog.Slog.InfoF(value, "%s -> %s %s -> %dms", addr, ctx.Request.Method, ctx.Request.URL.Path, time.Now().Sub(start).Milliseconds())
	slog.Slog.PrintAll(value)
}
