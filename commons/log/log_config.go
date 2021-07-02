package slog

import (
	"context"
	"github.com/kataras/iris/v12"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
	"io"
	"os"
	"path"
	"time"
)

var ZapLogger *zap.Logger
var Log *zap.SugaredLogger
var Slog LogConfig

type LogConfig struct {
	Level    string `yaml:"level"`
	Path     string `yaml:"path"`
	FileName string `yaml:"filename"`
}

func init() {
	InitLogger(LogConfig{}, nil)
}

func InitLogger(logConfig LogConfig, app *iris.Application) {
	encoder := getEncoder()
	var writer io.Writer
	if logConfig.FileName != "" {
		writer = io.MultiWriter(os.Stdout, getLogWriter(logConfig.Path, logConfig.FileName))
	} else {
		writer = os.Stdout
	}
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(writer), zapcore.InfoLevel),
	)
	// develop mode
	caller := zap.AddCaller()
	// open the code line
	development := zap.Development()
	ZapLogger = zap.New(core, caller, development, zap.AddCallerSkip(1))
	Log = ZapLogger.Sugar()

	//set iris log level
	if app != nil {
		app.Logger().SetLevel(logConfig.Level)
		app.Logger().SetOutput(writer)
	}
}

/**
 * time format
 */
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("[2006-01-02 15:04:05]"))
}

/**
 * get zap log encoder
 */
func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.LineEnding = zapcore.DefaultLineEnding
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(logPath, level string) io.Writer {
	logFullPath := path.Join(logPath, level)
	hook, err := rotatelogs.New(
		logFullPath+"-%Y%m%d%H"+".txt",
		// log file split
		rotatelogs.WithRotationTime(24*time.Hour),
	)
	if err != nil {
		panic(err)
	}
	return hook
}

func (LogConfig) InfoF(template string, args ...interface{}) {
	Log.Infof(template, args)
}

func (LogConfig) DebugF(template string, args ...interface{}) {
	Log.Debugf(template, args)
}

func (LogConfig) ErrorF(template string, args ...interface{}) {
	Log.Errorf(template, args)
}
func (LogConfig) Warn(ctx context.Context, template string, args ...interface{}) {
	Log.Warnf(template, args)
}

func (l *LogConfig) Print(v ...interface{}) {
	Log.Info(v)
}
func (l *LogConfig) Printf(format string, v ...interface{}) {
	Log.Infof(format, v)
}
func (l *LogConfig) LogMode(level logger.LogLevel) logger.Interface {
	return l
}
func (l *LogConfig) Info(ctx context.Context, template string, args ...interface{}) {
	Log.Infof(template, args)
}
func (l *LogConfig) Error(ctx context.Context, template string, args ...interface{}) {
	Log.Errorf(template, args)
}
func (l *LogConfig) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
}
