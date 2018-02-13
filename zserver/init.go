package zserver

import (
	"github.com/Djoulzy/Tools/clog"
)

var cryptor *cypher
var scaleList *serversList

// var Storage *storage.Driver

var zeWorld *world
var zehub *hubManager

// StartZServer Init and start the ZUMBY Server
func StartZServer() {

	cryptor = &cypher{
		HashSize: ZConf.HashSize,
		HexKey:   []byte(ZConf.HexKey),
		HexIV:    []byte(ZConf.HexIV),
	}

	///////////////////////////////////////////////////////////////////////////

	zehub = newhubManager()
	zeWorld = worldInit()
	clog.ServiceCallback = zeWorld.sendServerMassage

	go monStart()

	scaleList = scaleInit(&ZConf.KnownBrothers.Servers)
	go scaleList.scaleStart()
	// go scaling.Start(ScalingServers)

	clog.Output("HTTP Server starting listening on %s", ZConf.HTTPaddr)
	go httpStart()

	clog.Output("TCP Server starting listening on %s", ZConf.TCPaddr)
	go tcpStart()

	go zeWorld.run()
	zehub.run()
}
