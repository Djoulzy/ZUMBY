package zserver

import (
	"encoding/json"
	"fmt"
	"math"
	"runtime"
	"time"

	"github.com/Djoulzy/Tools/clog"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

const statsTimer = 5 * time.Second

type hubClientList map[string]*hubClient

type brother struct {
	Tcpaddr  string
	Httpaddr string
}

type serverMetrics struct {
	SID      string
	TCPADDR  string
	HTTPADDR string
	HOST     string
	CPU      int
	GORTNE   int
	STTME    string
	UPTME    string
	LSTUPDT  string
	LAVG     int
	MEM      string
	SWAP     string
	NBMESS   int
	NBI      int
	MXI      int
	NBU      int
	MXU      int
	NBM      int
	MXM      int
	NBS      int
	MXS      int
	BRTHLST  map[string]brother
}

type brotherList struct {
	BRTHLST map[string]brother
}

type hubClientsRegister struct {
	ID   hubClientList
	Name hubClientList
	Type map[int]hubClientList
}

type monParams struct {
	ServerID          string
	Httpaddr          string
	Tcpaddr           string
	MaxUsersConns     int
	MaxMonitorsConns  int
	MaxServersConns   int
	MaxIncommingConns int
}

var startTime time.Time
var upTime time.Duration
var machineLoad *load.AvgStat
var nbcpu int
var cr hubClientsRegister
var addBrother = make(chan map[string]brother)
var brotherlist = make(map[string]brother)

func getMemUsage() string {
	v, _ := mem.VirtualMemory()
	return fmt.Sprintf("<th>Mem</th><td class='memCell'>%v Mo</td><td class='memCell'>%v Mo</td><td class='memCell'>%.1f%%</td>", (v.Total / 1048576), (v.Free / 1048576), v.UsedPercent)
}

func getSwapUsage() string {
	v, _ := mem.SwapMemory()
	return fmt.Sprintf("<th>Swap</th><td class='memCell'>%v Mo</td><td class='memCell'>%v Mo</td><td class='memCell'>%.1f%%</td>", (v.Total / 1048576), (v.Free / 1048576), v.UsedPercent)
}

func addToBrothersList(srv map[string]brother) {
	for name, infos := range srv {
		brotherlist[name] = infos
	}
}

func monStart() {
	ticker := time.NewTicker(statsTimer)
	machineLoad = &load.AvgStat{Load1: 0, Load5: 0, Load15: 0}
	nbcpu, _ := cpu.Counts(true)
	startTime = time.Now()

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case newSrv := <-addBrother:
			addToBrothersList(newSrv)
		case <-ticker.C:
			tmp, _ := load.Avg()
			machineLoad = tmp
			loadIndice := int(math.Ceil((((machineLoad.Load1*5 + machineLoad.Load5*3 + machineLoad.Load15*2) / 10) / float64(nbcpu)) * 100))
			// mess := NewdataMessage(nil, machineLoad.String())
			t := time.Now()
			upTime = time.Since(startTime)

			newStats := serverMetrics{
				SID:      ZConf.Name,
				TCPADDR:  ZConf.TCPaddr,
				HTTPADDR: ZConf.HTTPaddr,
				HOST:     fmt.Sprintf("HTTP: %s - TCP: %s", ZConf.HTTPaddr, ZConf.TCPaddr),
				CPU:      nbcpu,
				GORTNE:   runtime.NumGoroutine(),
				STTME:    startTime.Format("02/01/2006 15:04:05"),
				UPTME:    upTime.String(),
				LSTUPDT:  t.Format("02/01/2006 15:04:05"),
				LAVG:     loadIndice,
				MEM:      getMemUsage(),
				SWAP:     getSwapUsage(),
				NBMESS:   zehub.SentMessByTicks,
				NBI:      len(zehub.Incomming),
				MXI:      ZConf.MaxIncommingConns,
				NBU:      len(zehub.Users),
				MXU:      ZConf.MaxUsersConns,
				NBM:      len(zehub.Monitors),
				MXM:      ZConf.MaxMonitorsConns,
				NBS:      len(zehub.Servers),
				MXS:      ZConf.MaxServersConns,
				BRTHLST:  brotherlist,
			}

			newBrthList := brotherList{
				BRTHLST: brotherlist,
			}

			brthJSON, _ := json.Marshal(newBrthList)
			json, err := json.Marshal(newStats)
			if err != nil {
				clog.Error("Monitoring", "LoadAverage", "MON: Cannot send server metrics to listeners ...")
			} else {
				if len(zehub.Monitors)+len(zehub.Servers) > 0 {
					zehub.SentMessByTicks = 0
					mess := newDataMessage(nil, clientMonitor, nil, json)
					zehub.Broadcast <- mess
					mess = newDataMessage(nil, clientServer, nil, append([]byte("[MNIT]"), json...))
					zehub.Broadcast <- mess
					mess = newDataMessage(nil, clientUser, nil, append([]byte("[FLBK]"), brthJSON...))
					zehub.Broadcast <- mess
				}
			}
		}
	}
}
