package utils

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/jau1jz/cornus/commons"
	"github.com/jau1jz/cornus/commons/log"
	"github.com/kataras/iris/v12"
	uuid "github.com/satori/go.uuid"
)

func GenerateUUID() string {
	return uuid.NewV4().String()
}

func StringToMd5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func RetryFunction(c func() bool, times int) bool {
	for i := times + 1; i > 0; i-- {
		if c() == true {
			return true
		}
	}
	return false
}

func ValidateAndBindParameters(entity interface{}, ctx *iris.Context, info string) (commons.ResponseCode, string) {
	if err := (*ctx).UnmarshalBody(entity, iris.UnmarshalerFunc(json.Unmarshal)); err != nil {
		log.Slog.ErrorF((*ctx).Value("ctx").(context.Context), "%s error %s", info, err.Error())
		return commons.ParameterError, err.Error()
	}
	if err := Validate(entity); err != nil {
		log.Slog.ErrorF((*ctx).Value("ctx").(context.Context), "%s error %s", info, err.Error())
		return commons.ValidateError, err.Error()
	}
	return commons.OK, ""
}
