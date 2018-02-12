package main

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"

	"github.com/Djoulzy/Tools/clog"
)

const (
	writeWait      = 5 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var Upgrader *websocket.Upgrader

type HTTPManager struct {
	Httpaddr         string
	ServerName       string
	Hub              *Hub
	ReadBufferSize   int
	WriteBufferSize  int
	NBAcceptBySecond int
	HandshakeTimeout int
	CallToAction     func(*Client, []byte)
	Cryptor          *Cypher
	MapGenCallback   func(x, y int) []byte
	GetTilesList     func() []byte
	GetMapImg        func(x, y int) string
	ClientDisconnect func(string)
	WorldWidth       int
	WorldHeight      int
}

func (m *HTTPManager) statusPage(w http.ResponseWriter, r *http.Request) {
	handShake, _ := m.Cryptor.Encrypt_b64("MNTR|Monitoring|MNTR")
	var data = struct {
		Host   string
		Nb     int
		Users  map[string]*Client
		Stats  string
		HShake string
	}{
		m.Httpaddr,
		len(m.Hub.Users),
		m.Hub.Users,
		MachineLoad.String(),
		string(handShake),
	}

	zeFile := "../public/status.html"
	clog.Test("HTTPServer", "statusPage", "%s", zeFile)
	homeTempl, err := template.ParseFiles(zeFile)
	if err != nil {
		clog.Error("HTTPServer", "statusPage", "%s", err)
		return
	}
	homeTempl.Execute(w, &data)
}

func (m *HTTPManager) testPage(w http.ResponseWriter, r *http.Request) {
	handShake, _ := m.Cryptor.Encrypt_b64("LOAD_1|TestPage|USER")

	var data = struct {
		Host   string
		HShake string
	}{
		m.Httpaddr,
		string(handShake),
	}

	homeTempl, err := template.ParseFiles("../public/client.html")
	if err != nil {
		clog.Error("HTTPServer", "testPage", "%s", err)
		return
	}
	homeTempl.Execute(w, &data)
}

func (m *HTTPManager) Connect() *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: m.Httpaddr, Path: "/ws"}
	clog.Info("HTTPServer", "Connect", "Connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		clog.Error("HTTPServer", "Connect", "%s", err)
		return nil
	}

	return conn
}

func (m *HTTPManager) Reader(conn *websocket.Conn, cli *Client) {
	defer func() {
		m.ClientDisconnect(cli.Name)
		m.Hub.Unregister <- cli
		conn.Close()
	}()
	conn.SetReadLimit(maxMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		// clog.Debug("HTTPServer", "Reader", "PONG! from %s", cli.Name)
		return nil
	})
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				clog.Error("HTTPServer", "Reader", "%v - %s", err, message)
			}
			return
		}
		message = bytes.TrimSpace(bytes.Replace(message, Newline, Space, -1))
		go m.CallToAction(cli, message)
	}
}

func (m *HTTPManager) _write(ws *websocket.Conn, mt int, message []byte) error {
	ws.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.WriteMessage(mt, message)
}

func (m *HTTPManager) Writer(conn *websocket.Conn, cli *Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		conn.Close()
	}()

	for {
		select {
		case message, ok := <-cli.Send:
			if !ok {
				clog.Warn("HTTPServer", "Writer", "Error: %s", ok)
				cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Something went wrong !")
				if err := m._write(conn, websocket.CloseMessage, cm); err != nil {
					clog.Error("HTTPServer", "Writer", "Connection lost ! Cannot send CloseMessage to %s", cli.Name)
				}
				return
			}
			// clog.Debug("HTTPServer", "Writer", "Sending: %s", message)
			if err := m._write(conn, websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			clog.Debug("HTTPServer", "Writer", "Client %s Ping!", cli.Name)
			if err := m._write(conn, websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case <-cli.Quit:
			cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "An other device is using your account !")
			if err := m._write(conn, websocket.CloseMessage, cm); err != nil {
				clog.Error("HTTPServer", "Writer", "Cannot write CloseMessage to %s", cli.Name)
			}
			return
		}
	}
}

// serveWs handles websocket requests from the peer.
func (m *HTTPManager) wsConnect(w http.ResponseWriter, r *http.Request) {
	var ua string
	name := r.Header["Sec-Websocket-Key"][0]
	if len(r.Header["User-Agent"]) > 0 {
		ua = r.Header["User-Agent"][0]
	} else {
		ua = "n/a"
	}

	if m.Hub.UserExists(name, ClientUser) {
		clog.Warn("HTTPServer", "wsConnect", "Client %s already exists ... Refusing connection", name)
		return
	}

	httpconn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		clog.Error("HTTPServer", "wsConnect", "%s", err)
		httpconn.Close()
		return
	}

	client := &Client{Quit: make(chan bool),
		CType: ClientUndefined, Send: make(chan []byte, 256), Enqueue: make(chan []byte, 256),
		CallToAction: m.CallToAction, Addr: httpconn.RemoteAddr().String(),
		Name: name, Content_id: 0, Front_id: "", App_id: "", Country: "", User_agent: ua}

	m.Hub.Register <- client
	go m.Writer(httpconn, client)
	go m.Reader(httpconn, client)
}

func throttleClients(h http.Handler, n int) http.Handler {
	ticker := time.NewTicker(time.Second / time.Duration(n))
	// sema := make(chan struct{}, n)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// sema <- struct{}{}
		// defer func() { <-sema }()
		<-ticker.C
		h.ServeHTTP(w, r)
	})
}

func (m *HTTPManager) dataServe(w http.ResponseWriter, r *http.Request) {
	name := ".." + r.URL.Path
	file, err := os.Open(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.ServeContent(w, r, name, time.Now(), file)
	// m := mapper.NewMap()
	// mapJSON, _ := json.Marshal(m)
	// w.Write(mapJSON)
}

func (m *HTTPManager) getGameData(w http.ResponseWriter, r *http.Request) {
	var str []byte
	query := strings.Split(string(r.URL.Path[1:]), "/")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if query[1] == "TilesList.json" {
		str = m.GetTilesList()
	} else {
		str = []byte("")
	}

	w.Write(str)
}

func (m *HTTPManager) getMapArea(w http.ResponseWriter, r *http.Request) {
	query := strings.Split(string(r.URL.Path[1:]), "/")
	coord := strings.Split(query[1], "_")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	x, _ := strconv.Atoi(coord[0])
	y, _ := strconv.Atoi(coord[1])
	str := m.MapGenCallback(x, y)
	w.Write(str)
}

func (m *HTTPManager) getMonPage(w http.ResponseWriter, r *http.Request) {
	query := strings.Split(string(r.URL.Path[1:]), "/")
	coord := strings.Split(query[1], "_")

	if len(coord) == 2 {
		x, _ := strconv.Atoi(coord[0])
		y, _ := strconv.Atoi(coord[1])
		page := m.GetMapImg(x, y)
		w.Write([]byte(page))
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/html")
		m.GetMapImg(-1, -1)
		http.ServeFile(w, r, "../public/mon.html")
	}
}

func (m *HTTPManager) Start(conf *HTTPManager) {
	m = conf

	// myHttp := &http.Server{
	// 	Addr:           ":8080",
	// 	ReadTimeout:    10 * time.Second,
	// 	WriteTimeout:   10 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }

	Upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			clog.Error("httpserver", "Start", "Error %s", reason)
		},
		ReadBufferSize:   m.ReadBufferSize,
		WriteBufferSize:  m.WriteBufferSize,
		HandshakeTimeout: time.Duration(m.HandshakeTimeout) * time.Second,
	} // use default options

	fs := http.FileServer(http.Dir("../public/"))
	http.Handle("/client/", http.StripPrefix("/client/", fs))

	http.HandleFunc("/data/", m.dataServe)
	http.HandleFunc("/test", m.testPage)
	http.HandleFunc("/status", m.statusPage)
	http.HandleFunc("/map/", m.getMapArea)
	http.HandleFunc("/GameData/", m.getGameData)
	http.HandleFunc("/mon/", m.getMonPage)

	handler := http.HandlerFunc(m.wsConnect)
	http.Handle("/ws", throttleClients(handler, m.NBAcceptBySecond))

	clog.Fatal("httpserver", "Start", http.ListenAndServe(m.Httpaddr, nil))
}
