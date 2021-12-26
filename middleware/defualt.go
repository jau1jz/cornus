package middleware

import (
	"context"
	"fmt"
	slog "github.com/jau1jz/cornus/commons/log"
	"github.com/jau1jz/cornus/commons/utils"
	"github.com/kataras/iris/v12"
	"runtime"
	"time"
)

func Default(ctx iris.Context) {
	value := context.WithValue(context.Background(), "trace_id", utils.GenerateUUID())
	ctx.Values().Set("ctx", value)
	defer func() {
		if err := recover(); err != nil {
			if ctx.IsStopped() {
				return
			}

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
			ctx.StatusCode(500)
			ctx.StopExecution()
		}
	}()

	start := time.Now()
	ctx.Next()
	slog.Slog.InfoF(value, "%s -> %s -> %dms", ctx.Request().RemoteAddr, ctx.Request().URL.Path, time.Now().Sub(start).Milliseconds())
}
