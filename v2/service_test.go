package cornus

import (
	"testing"
)

func TestStart_Default_Server(t *testing.T) {
	server := GetCornusInstance()
	server.StartServer(HttpService)
	server.WaitClose()
}
