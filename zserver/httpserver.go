package zserver

import (
	"bytes"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Djoulzy/Tools/clog"
	"github.com/gorilla/websocket"
)

const (
	writeWait          = 5 * time.Second
	pongWait           = 60 * time.Second
	pingPeriod         = (pongWait * 9) / 10
	maxdataMessageSize = 512
)

var upgrader *websocket.Upgrader

func statusPage(w http.ResponseWriter, r *http.Request) {
	handShake, _ := cryptor.encryptB64("MNTR|Monitoring|MNTR")
	var data = struct {
		Host   string
		Nb     int
		Users  map[string]*hubClient
		Stats  string
		HShake string
	}{
		ZConf.HTTPaddr,
		len(zehub.Users),
		zehub.Users,
		machineLoad.String(),
		string(handShake),
	}

	zeFile := "public/status.html"
	clog.Test("HTTPServer", "statusPage", "%s", zeFile)
	homeTempl, err := template.ParseFiles(zeFile)
	if err != nil {
		clog.Error("HTTPServer", "statusPage", "%s", err)
		return
	}
	homeTempl.Execute(w, &data)
}

func testPage(w http.ResponseWriter, r *http.Request) {
	handShake, _ := cryptor.encryptB64("LOAD_1|TestPage|USER")

	var data = struct {
		Host   string
		HShake string
	}{
		ZConf.HTTPaddr,
		string(handShake),
	}

	homeTempl, err := template.ParseFiles("public/client.html")
	if err != nil {
		clog.Error("HTTPServer", "testPage", "%s", err)
		return
	}
	homeTempl.Execute(w, &data)
}

func httpConnect() *websocket.Conn {
	u := url.URL{Scheme: "ws", Host: ZConf.HTTPaddr, Path: "/ws"}
	clog.Info("HTTPServer", "Connect", "Connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		clog.Error("HTTPServer", "Connect", "%s", err)
		return nil
	}

	return conn
}

func httpReader(conn *websocket.Conn, cli *hubClient) {
	defer func() {
		zeWorld.dropUser(cli.Name)
		zehub.Unregister <- cli
		conn.Close()
	}()
	conn.SetReadLimit(maxdataMessageSize)
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		// clog.Debug("HTTPServer", "Reader", "PONG! from %s", cli.Name)
		return nil
	})
	for {
		mTyp, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				clog.Error("HTTPServer", "Reader", "[%d] %v - %s", mTyp, err, message)
			}
			return
		}
		message = bytes.TrimSpace(bytes.Replace(message, newLine, spaceChar, -1))
		go callToAction(cli, message)
	}
}

func _httpWriter(ws *websocket.Conn, mt int, message []byte) error {
	ws.SetWriteDeadline(time.Now().Add(writeWait))
	return ws.WriteMessage(mt, message)
}

func httpWriter(conn *websocket.Conn, cli *hubClient) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		// conn.Close()
	}()

	for {
		select {
		case message, ok := <-cli.Send:
			if !ok {
				clog.Warn("HTTPServer", "Writer", "Error: %s", ok)
				cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Something went wrong !")
				if err := _httpWriter(conn, websocket.CloseMessage, cm); err != nil {
					clog.Error("HTTPServer", "Writer", "Connection lost ! Cannot send ClosedataMessage to %s", cli.Name)
				}
				return
			}
			// clog.Debug("HTTPServer", "Writer", "Sending: %s", message)
			if err := _httpWriter(conn, websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			clog.Debug("HTTPServer", "Writer", "hubClient %s Ping!", cli.Name)
			if err := _httpWriter(conn, websocket.PingMessage, []byte{}); err != nil {
				return
			}
		case <-cli.Quit:
			cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "An other device is using your account !")
			_httpWriter(conn, websocket.CloseMessage, cm)
			// if err := _httpWriter(conn, websocket.CloseMessage, cm); err != nil {
			// 	clog.Error("HTTPServer", "Writer", "Cannot write ClosedataMessage to %s", cli.Name)
			// }
			return
		}
	}
}

// serveWs handles websocket requests from the peer.
func wsConnect(w http.ResponseWriter, r *http.Request) {
	var ua string
	name := r.Header["Sec-Websocket-Key"][0]
	if len(r.Header["User-Agent"]) > 0 {
		ua = r.Header["User-Agent"][0]
	} else {
		ua = "n/a"
	}

	if zehub.userExists(name, clientUser) {
		clog.Warn("HTTPServer", "wsConnect", "hubClient %s already exists ... Refusing connection", name)
		return
	}

	httpconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		clog.Error("HTTPServer", "wsConnect", "%s", err)
		httpconn.Close()
		return
	}

	client := &hubClient{Quit: make(chan bool),
		CType: clientUndefined, Send: make(chan []byte, 256), Enqueue: make(chan []byte, 256),
		Addr: httpconn.RemoteAddr().String(),
		Name: name, AppID: "", Country: "", UserAgent: ua}

	zehub.Register <- client
	go httpWriter(httpconn, client)
	go httpReader(httpconn, client)
}

func throttlehubClients(h http.Handler, n int) http.Handler {
	ticker := time.NewTicker(time.Second / time.Duration(n))
	// sema := make(chan struct{}, n)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// sema <- struct{}{}
		// defer func() { <-sema }()
		<-ticker.C
		h.ServeHTTP(w, r)
	})
}

func dataServe(w http.ResponseWriter, r *http.Request) {
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

func getGameData(w http.ResponseWriter, r *http.Request) {
	var str []byte
	query := strings.Split(string(r.URL.Path[1:]), "/")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	if query[1] == "TilesList.json" {
		str = zeWorld.getTilesList()
	} else {
		str = []byte("")
	}

	w.Write(str)
}

func getMapArea(w http.ResponseWriter, r *http.Request) {
	query := strings.Split(string(r.URL.Path[1:]), "/")
	coord := strings.Split(query[1], "_")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")

	x, _ := strconv.Atoi(coord[0])
	y, _ := strconv.Atoi(coord[1])
	str := zeWorld.getMapArea(x, y)
	w.Write(str)
}

func getMonPage(w http.ResponseWriter, r *http.Request) {
	query := strings.Split(string(r.URL.Path[1:]), "/")
	coord := strings.Split(query[1], "_")

	if len(coord) == 2 {
		x, _ := strconv.Atoi(coord[0])
		y, _ := strconv.Atoi(coord[1])
		page := zeWorld.getMapImg(x, y)
		w.Write([]byte(page))
	} else {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/html")
		zeWorld.getMapImg(-1, -1)
		http.ServeFile(w, r, "public/mon.html")
	}
}

// HTTPStart lance le serveur HTTP
func httpStart() {

	// myHttp := &http.Server{
	// 	Addr:           ":8080",
	// 	ReadTimeout:    10 * time.Second,
	// 	WriteTimeout:   10 * time.Second,
	// 	MaxHeaderBytes: 1 << 20,
	// }

	upgrader = &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
			clog.Error("httpserver", "Start", "Error %s", reason)
		},
		ReadBufferSize:   ZConf.ReadBufferSize,
		WriteBufferSize:  ZConf.WriteBufferSize,
		HandshakeTimeout: time.Duration(ZConf.HandshakeTimeout) * time.Second,
	} // use default options

	fs := http.FileServer(http.Dir("public/"))
	http.Handle("/client/", http.StripPrefix("/client/", fs))

	http.HandleFunc("/data/", dataServe)
	http.HandleFunc("/test", testPage)
	http.HandleFunc("/status", statusPage)
	http.HandleFunc("/map/", getMapArea)
	http.HandleFunc("/GameData/", getGameData)
	http.HandleFunc("/mon/", getMonPage)

	handler := http.HandlerFunc(wsConnect)
	http.Handle("/ws", throttlehubClients(handler, ZConf.NBAcceptBySecond))

	clog.Fatal("httpserver", "Start", http.ListenAndServe(ZConf.HTTPaddr, nil))
}
