package zserver

import (
	"os"
	"testing"

	"github.com/Djoulzy/Tools/clog"
)

var tmphubManager *hubManager
var slist *serversList

func TestMain(m *testing.M) {
	clog.LogLevel = 5
	clog.StartLogging = true

	tmphubManager = newhubManager()
	go tmphubManager.run()

	srvList := make(map[string]string)
	srvList["srv2"] = "127.0.0.3"
	srvList["srv2"] = "127.0.0.5"

	os.Exit(m.Run())
}
