package zserver

import (
	"encoding/json"
	"fmt"
	"time"

	// "github.com/davecgh/go-spew/spew"
	"github.com/Djoulzy/Tools/clog"
)

const (
	clientUndefined = 0
	clientUser      = 1
	clientServer    = 2
	clientMonitor   = 3
	everybody       = 4
	timeStep        = 100 * time.Millisecond // Actualisation 10 par seconde
)

var clientTypeName = [4]string{"Incomming", "Users", "Servers", "Monitors"}

// const (
// 	ReadOnly  = 1
// 	WriteOnly = 2
// 	ReadWrite = 3
// )

// hubClient is a middleman between the websocket connection and the hub.
type hubClient struct {
	ID      string
	Send    chan []byte
	Enqueue chan []byte
	Quit    chan bool

	Addr      string
	CType     int
	Name      string
	AppID     string
	Country   string
	UserAgent string
}

type dataMessage struct {
	Level    int
	From     *hubClient
	UserType int
	Dest     *hubClient
	Content  []byte
}

var broadcastQueue [][]byte

type connModifier struct {
	hubClient *hubClient
	NewName   string
	NewType   int
}

type hubManager struct {
	// Registered clients.
	Incomming map[string]*hubClient
	Users     map[string]*hubClient
	Servers   map[string]*hubClient
	Monitors  map[string]*hubClient

	SentMessByTicks int

	FullUsersList [4](map[string]*hubClient)

	// Inbound messages from the clients.
	Register   chan *hubClient
	Unregister chan *hubClient
	Broadcast  chan *dataMessage
	// Status     chan *dataMessage
	Unicast chan *dataMessage
	Action  chan *dataMessage
	Done    chan bool
}

func newhubManager() *hubManager {
	hub := &hubManager{
		Register:   make(chan *hubClient),
		Unregister: make(chan *hubClient),

		Broadcast: make(chan *dataMessage),
		// Status:    make(chan *dataMessage),
		Unicast: make(chan *dataMessage),
		Action:  make(chan *dataMessage),
		Done:    make(chan bool),

		Users:     make(map[string]*hubClient),
		Incomming: make(map[string]*hubClient),
		Servers:   make(map[string]*hubClient),
		Monitors:  make(map[string]*hubClient),
	}
	hub.FullUsersList = [4](map[string]*hubClient){hub.Incomming, hub.Users, hub.Servers, hub.Monitors}
	return hub
}

func newDataMessage(from *hubClient, userType int, c *hubClient, content []byte) *dataMessage {
	m := &dataMessage{
		Level:    1,
		From:     from,
		UserType: userType,
		Dest:     c,
		Content:  content,
	}
	return m
}

func (h *hubManager) gethubClientByName(name string, userType int) *hubClient {
	return h.FullUsersList[userType][name]
}

func (h *hubManager) userExists(name string, userType int) bool {
	if h.FullUsersList[userType][name] != nil {
		return true
	}
	return false
}

func (h *hubManager) isRegistered(client *hubClient) bool {
	if h.FullUsersList[client.CType][client.Name] != nil {
		if h.FullUsersList[client.CType][client.Name].ID == client.ID {
			return true
		}
	}
	return false
}

func (h *hubManager) register(client *hubClient) {
	client.ID = fmt.Sprintf("%p", client)

	if h.userExists(client.Name, client.CType) {
		clog.Warn("hubManager", "Register", "hubClient %s already exists ... replacing", client.Name)
		h.unregister(h.FullUsersList[client.CType][client.Name])
	}

	h.FullUsersList[client.CType][client.Name] = client
	clog.Info("hubManager", "Register", "hubClient %s registered [%s] as %s.", client.Name, client.ID, clientTypeName[client.CType])
}

func (h *hubManager) unregister(client *hubClient) {
	if h.isRegistered(client) {
		delete(h.FullUsersList[client.CType], client.Name)

		// select {
		// case client.Quit <- true:
		// }

		close(client.Send)
		close(client.Quit)

		if client.CType == clientServer {
			data := struct {
				SID  string
				DOWN bool
			}{
				client.Name,
				true,
			}
			json, _ := json.Marshal(data)
			mess := newDataMessage(client, clientMonitor, nil, json)
			clog.Trace("hubManager", "Unregister", "Broadcasting close of server %s : %s", client.Name, json)
			h.broadcast(mess)
		}
		clog.Info("hubManager", "Unregister", "hubClient %s unregistered [%s] from %s.", client.Name, client.ID, clientTypeName[client.CType])
	}
}

func (h *hubManager) newRole(modif *connModifier) {
	if h.userExists(modif.NewName, modif.NewType) {
		clog.Warn("hubManager", "Newrole", "hubClient already exists ... Deleting")
		h.Unregister <- h.gethubClientByName(modif.NewName, modif.NewType)
	}
	delete(h.FullUsersList[modif.hubClient.CType], modif.hubClient.Name)
	modif.hubClient.Name = modif.NewName
	modif.hubClient.CType = modif.NewType
	h.FullUsersList[modif.NewType][modif.NewName] = modif.hubClient
}

func (h *hubManager) broadcast(message *dataMessage) {
	if message.UserType == clientUser {
		broadcastQueue = append(broadcastQueue, message.Content)
		// clog.Info("hubManager", "broadcast", "New message queued: (%d) - %s", len(BroadcastQueue), message.Content)
		h.flushBroadcastQueue()
	} else {
		list := h.FullUsersList[message.UserType]
		for _, client := range list {
			client.Send <- message.Content
			h.SentMessByTicks++
		}
	}
}

func (h *hubManager) flushBroadcastQueue() {
	list := h.FullUsersList[clientUser]
	for i, mess := range broadcastQueue {
		for _, client := range list {
			client.Send <- mess
			h.SentMessByTicks++
		}
		broadcastQueue[i] = nil
	}
	broadcastQueue = broadcastQueue[:0]
	// clog.Debug("hubManager", "flushBroadcastQueue", "Queue flushed")
}

func (h *hubManager) unicast(message *dataMessage) {
	message.Dest.Send <- message.Content
	// clog.Debug("hubManager", "unicast", "Unicast dataMessage to %s : %s", message.Dest.Name, message.Content)
	h.SentMessByTicks++
}

func (h *hubManager) action(message *dataMessage) {
	// clog.Debug("hubManager", "action", "dataMessage %s : %s", message.Dest.Name, message.Content)
	go callToAction(message.Dest, message.Content)
}

func (h *hubManager) run() {
	// ticker := time.NewTicker(timeStep)
	// defer func() {
	// 	ticker.Stop()
	// }()

	for {
		select {
		// case <-ticker.C:
		// 	if len(BroadcastQueue) > 0 {
		// 		h.flushBroadcastQueue()
		// 	}
		case client := <-h.Register:
			h.register(client)
			// client.Consistent <- true
		case client := <-h.Unregister:
			h.unregister(client)
			// client.Consistent <- true
		// case message := <-h.Status:
		// 	h.updateStatus(message)
		case message := <-h.Broadcast:
			h.broadcast(message)
		case message := <-h.Unicast:
			h.unicast(message)
		case message := <-h.Action:
			h.action(message)
		case <-h.Done:
			return
		}
	}
}
