package cornus

import (
	"context"
	"github.com/jau1jz/cornus/oss"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"

	redisV8 "github.com/go-redis/redis/v8"
	"github.com/jau1jz/cornus/commons"
	slog "github.com/jau1jz/cornus/commons/log"
	"github.com/jau1jz/cornus/config"
	"github.com/jau1jz/cornus/cornusdb"
	"github.com/jau1jz/cornus/iris"
	"github.com/jau1jz/cornus/kafka"
	"github.com/jau1jz/cornus/redis"
	irisV12 "github.com/kataras/iris/v12"
)

// Instance we need create the single object but thread safe
var Instance *Server

type Server struct {
	app   iris.App
	redis []redis.Redis
	db    []cornusdb.CornusDB
	kafka kafka.Kafka
	oss   oss.Client
}
type ServerOption int

const (
	DatabaseService = iota + 1
	RedisService
	OssService
)

func init() {
	Instance = &Server{}
}

// GetCornusInstance create the single object
func GetCornusInstance() *Server {
	return Instance
}

func (slf *Server) Default() {
	slf.app.Default()
}

func GetOSS() oss.Client {
	return Instance.oss
}

func (slf *Server) RegisterController(f func(app *irisV12.Application)) {
	f(slf.app.GetIrisApp())
}

func (slf *Server) RegisterErrorCodeAndMsg(language string, arr map[commons.ResponseCode]string) {
	if len(arr) == 0 {
		return
	}
	for k, v := range arr {
		commons.CodeMsg[language][k] = v
	}
}

func (slf *Server) WaitClose(params ...irisV12.Configurator) {
	defer func(ZapLog *zap.SugaredLogger) {
		_ = ZapLog.Sync()
	}(slog.ZapLog)
	go func() {
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
			_ = slf.app.GetIrisApp().Shutdown(ctx)
		}
	}()
	err := slf.app.Start(params...)
	if err != nil {
		panic(err)
	}
}
func (slf *Server) New() {
	slf.app.New()
}

// App return app
func (slf *Server) App() *iris.App {
	return &slf.app
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

// StartServer need call this function after Option, if Dependent service is not started return panic.
func (slf *Server) StartServer(opt ...ServerOption) {
	var err error
	for _, v := range opt {
		switch v {
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
		case OssService:
			slf.oss = oss.ClientInstance(config.Configs.Oss.OssBucket, config.Configs.Oss.AccessKeyID, config.Configs.Oss.AccessKeySecret, config.Configs.Oss.OssEndPoint)
		}
	}
}

func (slf *Server) KafkaService(ctx context.Context, topic string, callBackChan chan []byte) {
	slf.kafka.KafkaReceiver(ctx, topic, callBackChan)
}
