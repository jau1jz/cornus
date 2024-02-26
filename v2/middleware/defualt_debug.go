//go:build debug_log

package middleware

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jau1jz/cornus/v2/commons"
	slog "github.com/jau1jz/cornus/v2/commons/log"
	"github.com/jau1jz/cornus/v2/commons/utils"
	"io"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"
)

func Default(ctx *gin.Context) {
	value := context.WithValue(context.Background(), "trace_id", utils.GenerateUUID())
	ctx.Set("ctx", value)
	atomic.AddInt64(&commons.ActiveRequests, 1)
	defer atomic.AddInt64(&commons.ActiveRequests, -1)
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
	if _, ok := ignoreRequestMap.Load(ctx.Request.URL.Path); !ok {
		if ctx.Request.Method == http.MethodPost {
			all, err := io.ReadAll(ctx.Request.Body)
			if err != nil {
				slog.Slog.ErrorF(value, "ReadAll %s", err)
			} else if len(all) > 0 {
				slog.Slog.InfoF(value, "Body \n%s", string(all))
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(all))
			}
		}
		start := time.Now()
		ctx.Next()
		path := ctx.Request.URL.Path
		if ctx.Request.URL.RawQuery != "" {
			path += "?" + ctx.Request.URL.RawQuery
		}
		ip := ctx.ClientIP()
		if ctx.Request.Header.Get("X-Forwarded-For") != "" {
			ip = ctx.Request.Header.Get("X-Forwarded-For")
		}
		slog.Slog.InfoF(value, "[response code:%d] [%s] [%dms] [%s:%s]", ctx.Writer.Status(), ip, time.Now().Sub(start).Milliseconds(), ctx.Request.Method, path)
	}

	slog.Slog.PrintAll(value)
}
