package zserver

// ServerID Nom du server
type ServerID struct {
	Name string
}

// Globals variables systemes
type Globals struct {
	LogLevel     int
	StartLogging bool
	SetMaxProcs  int
	MaxOpenFiles uint64
}

// ConnectionLimit Gestion des connections
type ConnectionLimit struct {
	MaxUsersConns     int
	MaxMonitorsConns  int
	MaxServersConns   int
	MaxIncommingConns int
}

// ServersAddresses Adresses IP/NS du server
type ServersAddresses struct {
	HTTPaddr string
	TCPaddr  string
}

// KnownBrothers serveurs voisins
type KnownBrothers struct {
	Servers map[string]string
}

// HTTPServerConfig params HTTP
type HTTPServerConfig struct {
	ReadBufferSize   int
	WriteBufferSize  int
	NBAcceptBySecond int
	HandshakeTimeout int
}

// TCPServerConfig params TCP
type TCPServerConfig struct {
	ConnectTimeOut           int
	WriteTimeOut             int
	ScalingCheckServerPeriod int
}

// Encryption secret keys
type Encryption struct {
	HashSize int
	HexKey   string
	HexIV    string
}

// World param du monde
type World struct {
	TimeStep  int
	TileSize  int
	AOIWidth  int
	AOIHeight int
	MobSpeed  int
	MaxMobNum int
}

// AppConfig Structure globale
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

// ZConf global conf exportation
var ZConf = &AppConfig{
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
		HashSize: 8,
		HexKey:   "0000000000000000000000000000000000000000000000000000000000000000",
		HexIV:    "00000000000000000000000000000000",
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
