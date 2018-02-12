package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"sync"
	"time"

	// "github.com/davecgh/go-spew/spew"
	"github.com/Djoulzy/Tools/clog"
)

type TCPManager struct {
	Tcpaddr                  string
	ServerName               string
	Hub                      *Hub
	MaxServersConns          int
	ConnectTimeOut           int
	WriteTimeOut             int
	ScalingCheckServerPeriod int
	CallToAction             func(*Client, []byte)
	Cryptor                  *Cypher
}

func (m *TCPManager) reader(conn *net.TCPConn, cli *Client) {
	defer func() {
		conn.Close()
	}()

	// message := make([]byte, 1024)
	for {
		// conn.SetReadDeadline(time.Now().Add(time.Second * 10))
		message, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			clog.Trace("TCPserver", "reader", "closing conn %s", err)
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, Newline, Space, -1))
		// clog.Trace("TCPserver", "reader", "Reading %s", message)
		// long, err := conn.Read(message)
		// if err != nil {
		// 	break
		// }
		// message = message[:long-1]
		// spew.Dump(message)
		go m.CallToAction(cli, message)
	}
}

func (m *TCPManager) writer(conn *net.TCPConn, cli *Client) {
	defer func() {
		conn.Close()
	}()

	for {
		select {
		case <-cli.Quit:
			clog.Trace("TCPserver", "writer", "closing conn")
			return
		case message, ok := <-cli.Send:
			// clog.Debug("TCPserver", "writer", "Sending %s", message)
			if !ok {
				// The hub closed the channel.
				return
			}

			err := conn.SetWriteDeadline(time.Now().Add(time.Second))
			if err != nil {

				return
			}
			message = append(message, Newline...)
			conn.Write(message)
		}
	}
}

// func GetAddr(c *hub.Client) string {
// 	addr := c.Conn.(*net.TCPConn).RemoteAddr().String()
// 	ip := strings.Split(string(addr), "|")
// 	return ip[0]
// }

func (m *TCPManager) Connect(addr string) (*net.TCPConn, error) {
	conn, err := net.DialTimeout("tcp", addr, time.Second*time.Duration(m.ConnectTimeOut))
	// addr, _ := net.ResolveTCPAddr("tcp", m.Tcpaddr)
	// conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		clog.Error("TCPserver", "Connect", "Can't connect to server %s", addr)
		return nil, err
	}
	return conn.(*net.TCPConn), err
}

func (m *TCPManager) newClient(addr string, name string) *Client {
	client := &Client{Quit: make(chan bool),
		CType: ClientUndefined, Send: make(chan []byte, 256), CallToAction: m.CallToAction, Addr: addr,
		Name: name, Content_id: 0, Front_id: "", App_id: "", Country: "", User_agent: "TCP Socket"}
	m.Hub.Register <- client
	// <-client.Consistent
	return client
}

func (m *TCPManager) NewOutgoingConn(conn *net.TCPConn, toName string, wg *sync.WaitGroup) {
	clog.Debug("TCPserver", "NewOutgoingConn", "Contacting %s", conn.RemoteAddr().String())
	client := m.newClient(conn.RemoteAddr().String(), toName)
	handShake, _ := m.Cryptor.Encrypt_b64(fmt.Sprintf("%s|%s|SERV", m.ServerName, m.Tcpaddr))
	mess := NewMessage(nil, client.CType, client, append([]byte("[HELO]"), handShake...))
	m.Hub.Unicast <- mess

	go m.writer(conn, client)
	(*wg).Done()
	m.reader(conn, client)
	m.Hub.Unregister <- client
	// <-client.Consistent
}

func (m *TCPManager) NewIncommingConn(conn *net.TCPConn, wg *sync.WaitGroup) {
	client := m.newClient(conn.RemoteAddr().String(), conn.RemoteAddr().String())
	handShake, _ := m.Cryptor.Encrypt_b64(fmt.Sprintf("%s|%s|SERV", m.ServerName, m.Tcpaddr))
	mess := NewMessage(nil, client.CType, client, append([]byte("[HELO]"), handShake...))
	m.Hub.Unicast <- mess

	go m.writer(conn, client)
	(*wg).Done()
	m.reader(conn, client)
	m.Hub.Unregister <- client
	// <-client.Consistent
}

func (m *TCPManager) Start(conf *TCPManager) {
	var wg sync.WaitGroup

	m = conf

	formatedaddr, _ := net.ResolveTCPAddr("tcp", m.Tcpaddr)
	ln, err := net.ListenTCP("tcp", formatedaddr)
	if err != nil {
		clog.Error("TCPserver", "Start", "%s", err)
	}

	for {
		conn, err := ln.AcceptTCP()
		if err != nil {
			// handle error
		}
		wg.Add(1)
		go m.NewIncommingConn(conn, &wg)
		wg.Wait()
	}
}
