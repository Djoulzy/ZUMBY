package zserver

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Djoulzy/Tools/clog"
)

var serverCheckPeriod = 10 * time.Second

type nearbyServer struct {
	hubclient   *hubClient
	distantName string
	connected   bool
	cpuload     int
	freeslots   int
	httpaddr    string
	tcpaddr     string
}

type serversList struct {
	nodes           map[string]*nearbyServer
	localName       string
	localAddr       string
	MaxServersConns int
	hubManager      *hubManager
}

func (slist *serversList) updateMetrics(addr string, message []byte) {
	serv := slist.nodes[addr]
	h := slist.hubManager
	if len(h.Monitors)+len(h.Servers) > 0 {
		clog.Debug("Scaling", "updateMetrics", "Update Metrics for %s", serv.tcpaddr)

		var metrics serverMetrics

		err := json.Unmarshal(message, &metrics)
		if err != nil {
			clog.Error("Scaling", "updateMetrics", "Cannot reading distant server metrics")
			return
		}
		serv.cpuload = metrics.LAVG
		serv.freeslots = (metrics.MXU - metrics.NBU)
		serv.httpaddr = metrics.HTTPADDR

		for name, infos := range metrics.BRTHLST {
			slist.addNewPotentialServer(name, infos.Tcpaddr)
		}

		newSrv := make(map[string]brother)
		newSrv[metrics.SID] = brother{
			Httpaddr: metrics.HTTPADDR,
			Tcpaddr:  metrics.TCPADDR,
		}
		addBrother <- newSrv

		if len(h.Monitors) > 0 {
			mess := newDataMessage(nil, clientMonitor, nil, message)
			h.Broadcast <- mess
		}
	}
}

func (slist *serversList) checkingNewServers() {
	// var wg sync.WaitGroup

	// spew.Dump(slist)
	// for addr, node := range slist.nodes {
	// 	if node.hubclient == nil || node.hubclient.hubManager == nil {
	// 		conn, err := slist.tcpmanager.Connect(addr)
	// 		if err == nil {
	// 			clog.Trace("Scaling", "checkingNewServers", "Trying new server -> %s (%s)", node.distantName, addr)
	// 			wg.Add(1)
	// 			go slist.tcpmanager.NewOutgoingConn(conn, node.distantName, &wg)
	// 			wg.Wait()
	// 			node.connected = true
	// 		}
	// 	}
	// }
}

func (slist *serversList) addNewConnectedServer(c *hubClient) {
	clog.Info("Scaling", "AddNewConnectedServer", "Commit of server %s to scaling procedure.", c.Name)
	slist.nodes[c.Addr] = &nearbyServer{
		// manager: &tcpserver.Manager{
		// 	ServerName: c.Name,
		// 	hubManager:        c.hubManager,
		// 	Tcpaddr:    c.Addr,
		// },
		distantName: c.Name,
		tcpaddr:     c.Addr,
		connected:   true,
		hubclient:   c,
	}
}

func (slist *serversList) addNewPotentialServer(name string, addr string) {
	if slist.nodes[addr] == nil {
		if addr != slist.localAddr {
			clog.Info("Scaling", "AddNewPotentialServer", "New server : %s (%s)", name, addr)
			slist.nodes[addr] = &nearbyServer{
				// manager: &tcpserver.Manager{
				// 	ServerName: name,
				// 	Tcpaddr:    addr,
				// 	hubManager:        slist.hubManager,
				// },

				distantName: name,
				tcpaddr:     addr,
				connected:   false,
			}
		}
	}
}

func scaleInit(list *map[string]string) *serversList {
	slist := &serversList{
		nodes:           make(map[string]*nearbyServer),
		localName:       ZConf.Name,
		localAddr:       ZConf.TCPaddr,
		MaxServersConns: ZConf.MaxServersConns,
	}

	if list != nil {
		for name, serv := range *list {
			slist.addNewPotentialServer(name, serv)
		}
	}
	return slist
}

func (slist *serversList) redirectConnection(client *hubClient) bool {
	for _, node := range slist.nodes {
		if node.connected {
			clog.Trace("Scaling", "RedirectConnection", "Server %s CPU: %d Slots: %d", node.hubclient.Name, node.cpuload, node.freeslots)
			if node.cpuload < 80 && node.freeslots > 0 {
				redirect := fmt.Sprintf("[RDCT]%s", node.httpaddr)
				client.Send <- []byte(redirect)
				clog.Info("Scaling", "RedirectConnection", "hubClient redirect -> %s (%s)", node.hubclient.Name, node.httpaddr)
				return true
			}
			clog.Warn("Scaling", "RedirectConnection", "Server %s full ...", node.hubclient.Name)
		}
	}
	return false
}

func (slist *serversList) dispatchNewConnection(h *hubManager, name string) {
	message := []byte(fmt.Sprintf("[KILL]%s", name))
	mess := newDataMessage(nil, clientServer, nil, message)
	h.Broadcast <- mess
}

func (slist *serversList) scaleStart() {
	ticker := time.NewTicker(serverCheckPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			slist.checkingNewServers()
		}
	}
}
