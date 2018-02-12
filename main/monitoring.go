package main

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

func loadAverage(h *hubManager, p *monParams) {
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
				SID:      p.ServerID,
				TCPADDR:  p.Tcpaddr,
				HTTPADDR: p.Httpaddr,
				HOST:     fmt.Sprintf("HTTP: %s - TCP: %s", p.Httpaddr, p.Tcpaddr),
				CPU:      nbcpu,
				GORTNE:   runtime.NumGoroutine(),
				STTME:    startTime.Format("02/01/2006 15:04:05"),
				UPTME:    upTime.String(),
				LSTUPDT:  t.Format("02/01/2006 15:04:05"),
				LAVG:     loadIndice,
				MEM:      getMemUsage(),
				SWAP:     getSwapUsage(),
				NBMESS:   h.SentMessByTicks,
				NBI:      len(h.Incomming),
				MXI:      p.MaxIncommingConns,
				NBU:      len(h.Users),
				MXU:      p.MaxUsersConns,
				NBM:      len(h.Monitors),
				MXM:      p.MaxMonitorsConns,
				NBS:      len(h.Servers),
				MXS:      p.MaxServersConns,
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
				if len(h.Monitors)+len(h.Servers) > 0 {
					h.SentMessByTicks = 0
					mess := newDatamessage(nil, clientMonitor, nil, json)
					h.Broadcast <- mess
					mess = newDatamessage(nil, clientServer, nil, append([]byte("[MNIT]"), json...))
					h.Broadcast <- mess
					mess = newDatamessage(nil, clientUser, nil, append([]byte("[FLBK]"), brthJSON...))
					h.Broadcast <- mess
				}
			}
		}
	}
}

func monStart(hub *hubManager, p *monParams) {
	// addToBrothersList(list)
	loadAverage(hub, p)
}
