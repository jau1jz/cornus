package log

import (
	"context"
	"fmt"
	"github.com/jau1jz/cornus/commons"
	"github.com/jau1jz/cornus/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/gorm/logger"
	"io"
	"os"
	"time"
)

var ZapLogger *zap.Logger
var Log *zap.SugaredLogger
var Slog Logger

type Logger struct {
	LogLevel int
}

func init() {
	encoder := getEncoder()
	Slog = Logger{LogLevel: commons.LogLevel[config.SC.SConfigure.LogLevel]}
	writeSyncer := getLogWriter(fmt.Sprintf("%s/%s.log", config.SC.SConfigure.LogPath, config.SC.SConfigure.LogName))
	core := zapcore.NewCore(encoder, writeSyncer, commons.ZapLogLevel[config.SC.SConfigure.LogLevel])

	// zap.AddCaller()  添加将调用函数信息记录到日志中的功能。
	logger := zap.New(core, zap.AddCaller())
	Log = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.LineEnding = zapcore.DefaultLineEnding
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(logPath string) zapcore.WriteSyncer {
	return zapcore.AddSync(io.MultiWriter(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     1, //days
		LocalTime:  true,
		Compress:   true, // disabled by default
	}, os.Stdout))
}

func (Logger) InfoF(template string, args ...interface{}) {
	Log.Infof(template, args...)
}

func (Logger) DebugF(template string, args ...interface{}) {
	Log.Debugf(template, args...)
}

func (Logger) ErrorF(template string, args ...interface{}) {
	Log.Errorf(template, args...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	Log.Infof(format, v...)
}
func (l *Logger) Print(v ...interface{}) {
	Log.Info(v...)
}
func (l *Logger) LogMode(level logger.LogLevel) logger.Interface {
	l.LogLevel = int(level)
	return l
}
func (l *Logger) Info(ctx context.Context, template string, args ...interface{}) {
	Log.Infof(template, args...)
}
func (l *Logger) Warn(ctx context.Context, template string, args ...interface{}) {
	Log.Warnf(template, args...)
}
func (l *Logger) Error(ctx context.Context, template string, args ...interface{}) {
	Log.Warnf(template, args...)
}

func (l *Logger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	//elapsed := time.Since(begin)
	//if l.LogLevel <= commons.Debug {
	//	sql, rows := fc()
	//	Log.Infof("Sql : %s , Affected : %d , time: %d ms", sql, rows, elapsed.Milliseconds())
	//} else if elapsed > time.Second && l.LogLevel >= commons.Warn {
	//	sql, rows := fc()
	//	Log.Warnf("SLOW SQL : %s , Affected :%s , Excute Time: %d ms", sql, rows, elapsed.Milliseconds())
	//}
	panic("test")
}
