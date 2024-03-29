package utils

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/jau1jz/cornus/commons"
	"github.com/jau1jz/cornus/commons/log"
	"github.com/kataras/iris/v12"
	"strings"
)

func GenerateUUID() string {
	newUUID, _ := uuid.NewUUID()
	return strings.Replace(newUUID.String(), "-", "", -1)
}

func StringToMd5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func StringToSha256(str string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(str)))
}
func RetryFunction(c func() bool, times int) bool {
	for i := times + 1; i > 0; i-- {
		if c() == true {
			return true
		}
	}
	return false
}
func ValidateAndBindParameters(entity interface{}, ctx iris.Context, info string) (commons.ResponseCode, string) {
	if err := ctx.UnmarshalBody(entity, iris.UnmarshalerFunc(json.Unmarshal)); err != nil {
		log.Slog.ErrorF(ctx.Values().Get("ctx").(context.Context), "%s error %s", info, err.Error())
		return commons.ParameterError, err.Error()
	}
	if err := Validate(entity); err != nil {
		log.Slog.ErrorF(ctx.Values().Get("ctx").(context.Context), "%s error %s", info, err.Error())
		return commons.ValidateError, err.Error()
	}

	return commons.OK, ""
}

func ValidateAndBindCtxParameters(entity interface{}, ctx iris.Context, info string) (commons.ResponseCode, string) {
	err := json.Unmarshal(ctx.Values().Get(commons.CtxValueParameter).([]byte), entity)
	if err != nil {
		log.Slog.ErrorF(ctx.Values().Get("ctx").(context.Context), "%s error %s", info, err.Error())
		return commons.ParameterError, err.Error()
	}
	if err := Validate(entity); err != nil {
		log.Slog.ErrorF(ctx.Values().Get("ctx").(context.Context), "%s error %s", info, err.Error())
		return commons.ValidateError, err.Error()
	}

	return commons.OK, ""
}
