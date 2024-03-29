package cornus

import (
	"context"
	"errors"
	"fmt"
	"github.com/jau1jz/cornus/v2/middleware"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	redisV8 "github.com/go-redis/redis/v8"
	"github.com/jau1jz/cornus/v2/commons"
	slog "github.com/jau1jz/cornus/v2/commons/log"
	"github.com/jau1jz/cornus/v2/config"
	"github.com/jau1jz/cornus/v2/cornusdb"
	"github.com/jau1jz/cornus/v2/redis"
)

// Instance we need create the single object but thread safe
var Instance *Server

type Server struct {
	app        *gin.Engine
	redis      []redis.Redis
	db         []cornusdb.CornusDB
	ctx        context.Context
	httpServer *http.Server
}
type ServerOption int

const (
	DatabaseService = iota + 1
	RedisService
	OssService
	HttpService
)

func init() {
	Instance = &Server{}
}

// GetCornusInstance create the single object
func GetCornusInstance() *Server {
	return Instance
}
func (slf *Server) SetMysqlLogCallerSkip(skip int) {
	slog.GormSkip = skip
	slog.ReInit()
}
func (slf *Server) RegisterErrorCodeAndMsg(language string, arr map[commons.ResponseCode]string) {
	commons.RegisterCodeAndMsg(language, arr)
}

func (slf *Server) WaitClose() {
	defer func(ZapLog *zap.SugaredLogger) {
		_ = ZapLog.Sync()
	}(slog.ZapLog)
	//创建HTTP服务器
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.SC.SConfigure.Port),
		Handler: slf.app,
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch,
		// kill -SIGINT XXL 或 Ctrl+c
		os.Interrupt,
		syscall.SIGINT, // register that too, it should be ok
		// os.Kill等同于syscall.Kill
		os.Kill,
		syscall.SIGKILL, // register that too, it should be ok
		// kill -SIGTERM XXE
		//^
		syscall.SIGTERM,
	)
	select {
	case <-ch:
		slog.Slog.InfoF(context.Background(), "wait for close server")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		for _, db := range slf.db {
			_ = db.StopDb()
		}
		for {
			if atomic.LoadInt64(&commons.ActiveRequests) == 0 {
				break
			}
			time.Sleep(time.Second)
		}

		err := server.Shutdown(ctx)
		if err != nil {
			slog.Slog.ErrorF(context.Background(), err.Error())
		}
	}
}

// App return app
func (slf *Server) App() *gin.Engine {
	return slf.app
}
func (slf *Server) FeatureDB(name string) *cornusdb.CornusDB {
	for _, v := range slf.db {
		if v.Name() == name {
			return &v
		}
	}
	return nil
}
func (slf *Server) Redis(name string) *redisV8.Client {
	for _, v := range slf.redis {
		if v.Name() == name {
			return v.Redis()
		}
	}
	return nil
}
func (slf *Server) LoadCustomizeConfig(slfConfig interface{}) {
	err := config.LoadCustomizeConfig(slfConfig)
	if err != nil {
		panic(err)
	}
}
func (slf *Server) http() {
	//设置模式
	if config.SC.SConfigure.Profile == "prod" {
		gin.SetMode(gin.ReleaseMode)
	} else if config.SC.SConfigure.Profile == "test" {
		gin.SetMode(gin.TestMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	gin.ForceConsoleColor()
	slf.app = gin.New()
	//插入中间件
	slf.app.Use(middleware.Default)

	slf.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.SC.SConfigure.Port),
		Handler: slf.App(),
	}
	go func() {
		if err := slf.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Slog.ErrorF(context.Background(), err.Error())
		}
	}()
}

// StartServer need call this function after Option, if Dependent service is not started return panic.
func (slf *Server) StartServer(opt ...ServerOption) {
	var err error
	for _, v := range opt {
		switch v {
		case HttpService:
			slf.http()
		case DatabaseService:
			slf.db = make([]cornusdb.CornusDB, 0)
			for _, v := range config.Configs.DataBase {
				if v.Type == "sqlite" {
					db := cornusdb.CornusDB{}
					err = db.StartSqlite(v)
					if err != nil {
						panic(err)
					}
					slf.db = append(slf.db, db)
				} else if v.Type == "mysql" {
					db := cornusdb.CornusDB{}
					err = db.StartMysql(v)
					if err != nil {
						panic(err)
					}
					slf.db = append(slf.db, db)
				} else if v.Type == "pgsql" {
					db := cornusdb.CornusDB{}
					err = db.StartPgsql(v)
					if err != nil {
						panic(err)
					}
					slf.db = append(slf.db, db)
				} else {
					continue
				}
			}
		case RedisService:
			slf.redis = make([]redis.Redis, len(config.Configs.Redis))
			for i, v := range config.Configs.Redis {
				err = slf.redis[i].StartRedis(v)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
