package main

import (
	"encoding/json"
	"runtime"
	"syscall"

	"github.com/Djoulzy/Tools/clog"
	"github.com/Djoulzy/Tools/config"
)

var Cryptor *Cypher
var ScaleList *ServersList

// var Storage *storage.Driver

var zeWorld *WORLD
var zeHub *Hub

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

	Cryptor = &Cypher{
		HASH_SIZE: conf.HASH_SIZE,
		HEX_KEY:   []byte(conf.HEX_KEY),
		HEX_IV:    []byte(conf.HEX_IV),
	}

	conf_json, _ := json.Marshal(conf)

	///////////////////////////////////////////////////////////////////////////

	zeHub = NewHub()
	zeWorld = WorldInit(zeHub, conf_json)
	clog.ServiceCallback = zeWorld.SendServerMassage

	mon_params := &MonParams{
		ServerID:          conf.Name,
		Httpaddr:          conf.HTTPaddr,
		Tcpaddr:           conf.TCPaddr,
		MaxUsersConns:     conf.MaxUsersConns,
		MaxMonitorsConns:  conf.MaxMonitorsConns,
		MaxServersConns:   conf.MaxServersConns,
		MaxIncommingConns: conf.MaxIncommingConns,
	}
	go MonStart(zeHub, mon_params)

	tcp_params := &TCPManager{
		ServerName:               conf.Name,
		Tcpaddr:                  conf.TCPaddr,
		Hub:                      zeHub,
		ConnectTimeOut:           conf.ConnectTimeOut,
		WriteTimeOut:             conf.WriteTimeOut,
		ScalingCheckServerPeriod: conf.ScalingCheckServerPeriod,
		MaxServersConns:          conf.MaxServersConns,
		CallToAction:             CallToAction,
		Cryptor:                  Cryptor,
	}

	ScaleList = ScaleInit(tcp_params, &conf.KnownBrothers.Servers)
	go ScaleList.ScaleStart()
	// go scaling.Start(ScalingServers)

	http_params := &HTTPManager{
		ServerName:       conf.Name,
		Httpaddr:         conf.HTTPaddr,
		Hub:              zeHub,
		ReadBufferSize:   conf.ReadBufferSize,
		WriteBufferSize:  conf.WriteBufferSize,
		HandshakeTimeout: conf.HandshakeTimeout,
		NBAcceptBySecond: conf.NBAcceptBySecond,
		CallToAction:     CallToAction,
		Cryptor:          Cryptor,
		MapGenCallback:   zeWorld.GetMapArea,
		ClientDisconnect: zeWorld.DropUser,
		GetTilesList:     zeWorld.GetTilesList,
		GetMapImg:        zeWorld.GetMapImg,
		WorldWidth:       conf.AOIWidth,
		WorldHeight:      conf.AOIHeight,
	}
	clog.Output("HTTP Server starting listening on %s", conf.HTTPaddr)
	go http_params.Start(http_params)

	clog.Output("TCP Server starting listening on %s", conf.TCPaddr)
	go tcp_params.Start(tcp_params)

	go zeWorld.Run()
	zeHub.Run()
}
