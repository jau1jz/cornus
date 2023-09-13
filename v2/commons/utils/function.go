package utils

import (
	"context"
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"github.com/jau1jz/cornus/v2/commons"
	"github.com/jau1jz/cornus/v2/commons/log"
	"net/http"
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
func ValidateAndBindCtxParameters(entity interface{}, ctx *gin.Context, info string) (commons.ResponseCode, string) {
	if ctx.Request.Method == http.MethodPost {
		err := ctx.MustBindWith(entity, binding.JSON)
		if err != nil {
			log.Slog.ErrorF(ctx.Value("ctx").(context.Context), "%s error %s", info, err.Error())
			return commons.ParameterError, err.Error()
		}
		if err := Validate(entity); err != nil {
			log.Slog.ErrorF(ctx.Value("ctx").(context.Context), "%s error %s", info, err.Error())
			return commons.ValidateError, err.Error()
		}
	}

	return commons.OK, ""
}
