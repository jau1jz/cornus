package middleware

import (
	"bytes"
	"fmt"
	slog "github.com/jau1jz/cornus/commons/log"
	"github.com/jau1jz/cornus/commons/utils"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"runtime"
	"strings"
	"time"
)

func Default(ctx iris.Context) {
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
			slog.Slog.ErrorF(logMessage)
			ctx.StatusCode(500)
			ctx.StopExecution()
		}
	}()
	ctx.Values().Set("ctx_id", utils.GenerateUUID())
	// read base information and write log
	p := ctx.Request().URL.Path
	method := ctx.Request().Method
	start := time.Now().UnixNano() / 1e6
	ip := ctx.Request().RemoteAddr
	slog.Slog.InfoF("[path]--> %s [method]--> %s [IP]-->  %s", p, method, ip)

	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		slog.Slog.InfoF("ReadAll body failed: %s", err.Error())
	} else {
		ctx.Request().Body = ioutil.NopCloser(bytes.NewBuffer(body))
		if len(body) > 0 && strings.Contains(string(body), "{}") == false {
			slog.Slog.InfoF("log http request body: %s", string(body))
		}
	}

	// calculate cost time
	ctx.Next()
	end := time.Now().UnixNano() / 1e6
	slog.Slog.InfoF("[path]--> %s [cost time]ms-->  %d", p, end-start)
}
