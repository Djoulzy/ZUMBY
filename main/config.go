package main

type ServerID struct {
	Name string
}

type Globals struct {
	LogLevel     int
	StartLogging bool
	SetMaxProcs  int
	MaxOpenFiles uint64
}

type ConnectionLimit struct {
	MaxUsersConns     int
	MaxMonitorsConns  int
	MaxServersConns   int
	MaxIncommingConns int
}

type ServersAddresses struct {
	HTTPaddr string
	TCPaddr  string
}

type KnownBrothers struct {
	Servers map[string]string
}

type HTTPServerConfig struct {
	ReadBufferSize   int
	WriteBufferSize  int
	NBAcceptBySecond int
	HandshakeTimeout int
}

type TCPServerConfig struct {
	ConnectTimeOut           int
	WriteTimeOut             int
	ScalingCheckServerPeriod int
}

type Encryption struct {
	HASH_SIZE int
	HEX_KEY   string
	HEX_IV    string
}

type World struct {
	TimeStep  int
	TileSize  int
	AOIWidth  int
	AOIHeight int
	MobSpeed  int
	MaxMobNum int
}

type AppConfig struct {
	ServerID
	Globals
	ConnectionLimit
	ServersAddresses
	KnownBrothers
	HTTPServerConfig
	TCPServerConfig
	Encryption
	World
}

var conf *AppConfig = &AppConfig{
	ServerID{},
	Globals{
		LogLevel:     4,
		StartLogging: true,
		SetMaxProcs:  0,
		MaxOpenFiles: 0,
	},
	ConnectionLimit{
		MaxUsersConns:     100,
		MaxMonitorsConns:  3,
		MaxServersConns:   5,
		MaxIncommingConns: 50,
	},
	ServersAddresses{
		HTTPaddr: "localhost:8080",
		TCPaddr:  "localhost:8081",
	},
	KnownBrothers{},
	HTTPServerConfig{
		ReadBufferSize:   4096,
		WriteBufferSize:  4096,
		NBAcceptBySecond: 20,
		HandshakeTimeout: 5,
	},
	TCPServerConfig{
		ConnectTimeOut:           2,
		WriteTimeOut:             1,
		ScalingCheckServerPeriod: 10,
	},
	Encryption{
		HASH_SIZE: 8,
		HEX_KEY:   "0000000000000000000000000000000000000000000000000000000000000000",
		HEX_IV:    "00000000000000000000000000000000",
	},
	World{
		TimeStep:  100,
		TileSize:  32,
		AOIWidth:  30,
		AOIHeight: 30,
		MobSpeed:  8,
		MaxMobNum: 5,
	},
}
