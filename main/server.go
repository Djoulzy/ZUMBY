package main

import (
	"encoding/json"
	"runtime"
	"syscall"

	"github.com/Djoulzy/Tools/clog"
	"github.com/Djoulzy/Tools/config"
)

var cryptor *cypher
var scaleList *serversList

// var Storage *storage.Driver

var zeWorld *world
var zehub *hubManager

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
	config.Load("server.ini", conf)

	clog.LogLevel = conf.LogLevel
	clog.StartLogging = conf.StartLogging

	// System Optims
	setMaxProcs(conf.SetMaxProcs)
	maxFiles := maxOpenFiles(conf.MaxOpenFiles)
	////////////////

	totalConn := conf.MaxUsersConns + conf.MaxMonitorsConns + conf.MaxServersConns + conf.MaxIncommingConns
	if totalConn > maxFiles {
		conf.MaxUsersConns = maxFiles - 120
		conf.MaxMonitorsConns = 3
		conf.MaxServersConns = 10
		conf.MaxIncommingConns = 100
		clog.Warn("server", "main", "Setting MaxUser to %d.", conf.MaxUsersConns)
	}

	cryptor = &cypher{
		HashSize: conf.HashSize,
		HexKey:   []byte(conf.HexKey),
		HexIV:    []byte(conf.HexIV),
	}

	confJSON, _ := json.Marshal(conf)

	///////////////////////////////////////////////////////////////////////////

	zehub = newhubManager()
	zeWorld = worldInit(zehub, confJSON)
	clog.ServiceCallback = zeWorld.sendServerMassage

	go monStart()

	scaleList = scaleInit(&conf.KnownBrothers.Servers)
	go scaleList.scaleStart()
	// go scaling.Start(ScalingServers)

	clog.Output("HTTP Server starting listening on %s", conf.HTTPaddr)
	go httpStart()

	clog.Output("TCP Server starting listening on %s", conf.TCPaddr)
	go tcpStart()

	go zeWorld.run()
	zehub.run()
}
