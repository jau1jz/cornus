package cornus

import (
	"github.com/kataras/iris/v12"
	"testing"
)

func TestStart_Default_Server(t *testing.T) {
	server := GetCornusInstance()
	//http
	server.app.Default()
	server.StartServer(DatabaseService)

	Instance.app.Get("/ping", func(context iris.Context) {
		context.WriteString("123")
		panic("kkk")
	})
	server.WaitClose(iris.WithoutBodyConsumptionOnUnmarshal)

}
