package main

import (
	"runtime"
	"syscall"

	"github.com/Djoulzy/Tools/clog"
	"github.com/Djoulzy/Tools/config"
	"github.com/Djoulzy/ZUMBY/zserver"
)

func setMaxProcs(nb int) {
	var procs int
	if nb == 0 {
		procs = runtime.NumCPU()
		runtime.GOMAXPROCS(procs)
	} else {
		procs = nb
		runtime.GOMAXPROCS(procs)
	}
	clog.Output("Using %d CPUs on %d.", runtime.GOMAXPROCS(procs), runtime.NumCPU())
}

func maxOpenFiles(max uint64) int {
	var rLimit syscall.Rlimit

	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		clog.Error("server", "maxOpenFiles", "Error Getting Rlimit %s", err)
	}

	if max != 0 {
		rLimit.Cur = max - 1
		rLimit.Max = max
	}
	if rLimit.Cur < rLimit.Max {
		rLimit.Cur = rLimit.Max
		err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
		if err != nil {
			clog.Error("server", "maxOpenFiles", "Error Setting Rlimit %s", err)
		}
	}
	clog.Output("Setting maxOpenFiles to %d.", rLimit.Cur)
	return int(rLimit.Cur)
}

func main() {
	config.Load("server.ini", zserver.ZConf)

	clog.LogLevel = zserver.ZConf.LogLevel
	clog.StartLogging = zserver.ZConf.StartLogging

	// System Optims
	setMaxProcs(zserver.ZConf.SetMaxProcs)
	maxFiles := maxOpenFiles(zserver.ZConf.MaxOpenFiles)
	////////////////

	totalConn := zserver.ZConf.MaxUsersConns + zserver.ZConf.MaxMonitorsConns + zserver.ZConf.MaxServersConns + zserver.ZConf.MaxIncommingConns
	if totalConn > maxFiles {
		zserver.ZConf.MaxUsersConns = maxFiles - 120
		zserver.ZConf.MaxMonitorsConns = 3
		zserver.ZConf.MaxServersConns = 10
		zserver.ZConf.MaxIncommingConns = 100
		clog.Warn("server", "main", "Setting MaxUser to %d.", zserver.ZConf.MaxUsersConns)
	}

	zserver.StartZServer()
}
