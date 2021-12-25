package commons

import (
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
)

// ResponseCode define the error code
type ResponseCode int

const (
	HttpNotFound   ResponseCode = -2
	UnKnowError    ResponseCode = -1
	OK             ResponseCode = 0
	ParameterError ResponseCode = 1
	ValidateError  ResponseCode = 2
	TokenError     ResponseCode = 3
	CheckAuthError ResponseCode = 4
)

// CodeMsg global code and msg
var CodeMsg = map[ResponseCode]string{
	OK:             "成功",
	UnKnowError:    "未知错误",
	HttpNotFound:   "404",
	ParameterError: "参数错误",
	ValidateError:  "参数验证错误",
	TokenError:     "Token错误",
	CheckAuthError: "检查权限错误",
}

// GetCodeAndMsg construct the code and msg
func GetCodeAndMsg(code ResponseCode) string {
	value, ok := CodeMsg[code]
	if ok {
		return value
	}
	return "{}"
}

// RegisterCodeAndMsg msg will be used as default msg, and you can change msg with function 'BuildFailedWithMsg' or 'BuildSuccessWithMsg' or 'response.WithMsg' for once.
func RegisterCodeAndMsg(arr map[ResponseCode]string) {
	if len(arr) == 0 {
		return
	}
	for k, v := range arr {
		CodeMsg[k] = v
	}
}

const (
	Debug = iota - 1
	Info
	Warn
	Error
)

var LogLevel = map[string]int{
	"debug": -1,
	"info":  0,
	"warn":  1,
	"error": 2,
}

var ZapLogLevel = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
}

var GormLogLevel = map[string]logger.LogLevel{
	"debug": logger.Info,
	"info":  logger.Info,
	"warn":  logger.Warn,
	"error": logger.Error,
}
