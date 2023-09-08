//go:build debug_log

package log

import (
	"context"
	"fmt"
	"github.com/jau1jz/cornus/v2/commons"
	"github.com/jau1jz/cornus/v2/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"math/rand"
	"os"
	"runtime"
	"sync"
)

var Slog Logger
var Gorm GormLogger
var ZapLog *zap.SugaredLogger
var GormLog *zap.SugaredLogger

var logMap sync.Map

type Logger struct {
}

type logStack struct {
	F        func(ctx context.Context, template string, args ...interface{})
	Caller   string
	Ctx      context.Context
	Template string
	Args     []interface{}
}

func init() {
	encoder := getEncoder()
	Slog = Logger{}
	Gorm = GormLogger{
		LogLevel:                  commons.LogLevel[config.SC.SConfigure.LogLevel],
		IgnoreRecordNotFoundError: true,
	}
	writeSyncer := getLogWriter(fmt.Sprintf("%s/%s.log", config.SC.SConfigure.LogPath, config.SC.SConfigure.LogName))
	core := zapcore.NewCore(encoder, writeSyncer, commons.ZapLogLevel[config.SC.SConfigure.LogLevel])
	// zap.AddCaller()  添加将调用函数信息记录到日志中的功能。
	ZapLog = zap.New(core).Sugar()
	GormLog = zap.New(core).Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
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
func getTraceId(ctx context.Context) string {
	if traceId, ok := ctx.Value("trace_id").(string); ok {
		return fmt.Sprintf("trace_id: %s", traceId)
	} else {
		return ""
	}
}

func color(traceId string) string {
	return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m", 0x1B, 0, rand.Intn(36-31)+31, 40, traceId, 0x1B)
}
func pressLog(f func(ctx context.Context, template string, args ...interface{}), ctx context.Context, template string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(2)
	caller := fmt.Sprintf("%s:%d", file, line)
	newLog := logStack{
		F:        f,
		Caller:   caller,
		Ctx:      ctx,
		Template: template,
		Args:     args,
	}

	if actual, loaded := logMap.LoadOrStore(getTraceId(ctx), []logStack{
		newLog,
	}); loaded {
		logStacks := actual.([]logStack)
		logStacks = append(logStacks, newLog)
		logMap.Store(getTraceId(ctx), logStacks)
	}
}
func pressGormLog(f func(ctx context.Context, template string, args ...interface{}), ctx context.Context, template string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(6)
	caller := fmt.Sprintf("%s:%d", file, line)
	newLog := logStack{
		F:        f,
		Caller:   caller,
		Ctx:      ctx,
		Template: template,
		Args:     args,
	}

	if actual, loaded := logMap.LoadOrStore(getTraceId(ctx), []logStack{
		newLog,
	}); loaded {
		logStacks := actual.([]logStack)
		logStacks = append(logStacks, newLog)
		logMap.Store(getTraceId(ctx), logStacks)
	}
}
func (l *Logger) InfoF(ctx context.Context, template string, args ...interface{}) {
	pressLog(l.infoF, ctx, template, args...)
}
func (l *Logger) infoF(_ context.Context, template string, args ...interface{}) {
	ZapLog.Infof(template, args...)
}
func (l *Logger) DebugF(ctx context.Context, template string, args ...interface{}) {
	pressLog(l.debugF, ctx, template, args...)
}
func (l *Logger) debugF(_ context.Context, template string, args ...interface{}) {
	ZapLog.Debugf(template, args...)
}
func (l *Logger) ErrorF(ctx context.Context, template string, args ...interface{}) {
	pressLog(l.errorF, ctx, template, args...)
}
func (l *Logger) errorF(_ context.Context, template string, args ...interface{}) {
	ZapLog.Errorf(template, args...)
}
func (l *Logger) WarnF(ctx context.Context, template string, args ...interface{}) {
	pressLog(l.warnF, ctx, template, args...)
}
func (l *Logger) warnF(_ context.Context, template string, args ...interface{}) {
	ZapLog.Warnf(template, args...)
}
func (l *Logger) PrintAll(ctx context.Context) {
	value, ok := logMap.Load(getTraceId(ctx))
	if ok {
		logStacks := value.([]logStack)
		traceID := color(getTraceId(logStacks[0].Ctx)) + " "
		for _, logStack := range logStacks {
			logStack.F(context.Background(), fmt.Sprintf("%s %s %s", logStack.Caller, traceID, logStack.Template), logStack.Args...)
		}
	}
}
func (l *Logger) Printf(format string, v ...interface{}) {
	ZapLog.Infof(format, v...)
}
func (l *Logger) Print(v ...interface{}) {
	ZapLog.Info(v...)
}
