package iris

import (
	"errors"
	"fmt"
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/jau1jz/cornus/commons"
	"github.com/jau1jz/cornus/config"
	"github.com/jau1jz/cornus/middleware"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/pprof"
	"os"
)

type App struct {
	app *iris.Application
}

func (slf *App) Default() {
	slf.app = iris.New()
	//register middleware
	slf.app.UseGlobal(middleware.Default)
	//global error handling
	slf.app.OnAnyErrorCode(func(ctx iris.Context) {
		if ctx.GetStatusCode() == iris.StatusNotFound {
			_ = ctx.JSON(commons.BuildFailed(commons.HttpNotFound, commons.DefualtLanguage))
		} else {
			_ = ctx.JSON(commons.BuildFailed(commons.UnKnowError, commons.DefualtLanguage))
		}
	})
	slf.app.Logger().SetLevel(config.SC.SConfigure.LogLevel)
	slf.app.Logger().SetOutput(os.Stdout)
}

func (slf *App) New() {
	slf.app = iris.New()
	//global error handling
	slf.app.OnAnyErrorCode(func(ctx iris.Context) {
		if ctx.GetStatusCode() == iris.StatusNotFound {
			_ = ctx.JSON(commons.BuildFailed(commons.HttpNotFound, commons.DefualtLanguage))
		} else {
			_ = ctx.JSON(commons.BuildFailed(commons.UnKnowError, commons.DefualtLanguage))
		}
	})
	slf.app.Logger().SetLevel(config.SC.SConfigure.LogLevel)
	slf.app.Logger().SetOutput(os.Stdout)
}

// SetGlobalMiddleware set global middleware
func (slf *App) SetGlobalMiddleware(handlers ...context.Handler) {
	slf.app.UseGlobal(handlers...)
}

// SetMiddleware set middleware
func (slf *App) SetMiddleware(handlers ...context.Handler) {
	slf.app.Use(handlers...)
}

// GetIrisApp get Iris App
func (slf *App) GetIrisApp() *iris.Application {
	return slf.app
}

func (slf *App) Party(relativePath string, handlers ...context.Handler) {
	slf.app.Party(relativePath, handlers...)
}
func (slf *App) Post(relativePath string, handlers ...context.Handler) {
	slf.app.Post(relativePath, handlers...)
}
func (slf *App) Get(relativePath string, handlers ...context.Handler) {
	slf.app.Get(relativePath, handlers...)
}

// Start cornus server
func (slf *App) Start(params ...iris.Configurator) error {
	server := fmt.Sprintf("%s:%d", config.SC.SConfigure.Addr, config.SC.SConfigure.Port)
	if slf.app == nil {
		return errors.New("server not init")
	}
	//go slf.app.Run(iris.Addr(server))
	swaggerConfig := &swagger.Config{
		URL: fmt.Sprintf("./swagger/doc.json"), //The url pointing to API definition
	}
	slf.app.Get("/swagger/{any:path}", swagger.CustomWrapHandler(swaggerConfig, swaggerFiles.Handler))
	p := pprof.New()
	slf.app.Get("/debug/pprof", p)
	slf.app.Get("/debug/pprof/{action:path}", p)
	params = append(params, iris.WithoutStartupLog)

	return slf.app.Run(iris.Addr(server), params...)
}
