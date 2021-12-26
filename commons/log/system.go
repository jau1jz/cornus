package log

import (
	"context"
	"fmt"
	"github.com/jau1jz/cornus/commons"
	"github.com/jau1jz/cornus/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

var Slog Logger
var Gorm GormLogger
var ZapLog *zap.SugaredLogger

type Logger struct {
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
	logger := zap.New(core, zap.AddCaller())
	ZapLog = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05")
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
func getTraceId(ctx context.Context) string {
	if trace_id, ok := ctx.Value("trace_id").(string); ok {
		return fmt.Sprintf("trace_id: %s ", trace_id)
	} else {
		return ""
	}

}
func (l *Logger) InfoF(ctx context.Context, template string, args ...interface{}) {
	ZapLog.Infof(getTraceId(ctx)+template, args...)
}

func (l *Logger) DebugF(ctx context.Context, template string, args ...interface{}) {
	ZapLog.Debugf(getTraceId(ctx)+template, args...)
}

func (l *Logger) ErrorF(ctx context.Context, template string, args ...interface{}) {
	ZapLog.Errorf(getTraceId(ctx)+template, args...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	ZapLog.Infof(format, v...)
}
func (l *Logger) Print(v ...interface{}) {
	ZapLog.Info(v...)
}
