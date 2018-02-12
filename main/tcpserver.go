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

func tcpReader(conn *net.TCPConn, cli *hubClient) {
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
		message = bytes.TrimSpace(bytes.Replace(message, newLine, spaceChar, -1))
		// clog.Trace("TCPserver", "reader", "Reading %s", message)
		// long, err := conn.Read(message)
		// if err != nil {
		// 	break
		// }
		// message = message[:long-1]
		// spew.Dump(message)
		go callToAction(cli, message)
	}
}

func tcpWriter(conn *net.TCPConn, cli *hubClient) {
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
			message = append(message, newLine...)
			conn.Write(message)
		}
	}
}

// func GetAddr(c *hub.hubClient) string {
// 	addr := c.Conn.(*net.TCPConn).RemoteAddr().String()
// 	ip := strings.Split(string(addr), "|")
// 	return ip[0]
// }

func tcpConnect(addr string) (*net.TCPConn, error) {
	conn, err := net.DialTimeout("tcp", addr, time.Second*time.Duration(conf.ConnectTimeOut))
	// addr, _ := net.ResolveTCPAddr("tcp", m.Tcpaddr)
	// conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		clog.Error("TCPserver", "Connect", "Can't connect to server %s", addr)
		return nil, err
	}
	return conn.(*net.TCPConn), err
}

func newhubClient(addr string, name string) *hubClient {
	client := &hubClient{Quit: make(chan bool),
		CType: clientUndefined, Send: make(chan []byte, 256), Addr: addr,
		Name: name, AppID: "", Country: "", UserAgent: "TCP Socket"}
	zehub.Register <- client
	// <-client.Consistent
	return client
}

func newOutgoingTCPConn(conn *net.TCPConn, toName string, wg *sync.WaitGroup) {
	clog.Debug("TCPserver", "NewOutgoingConn", "Contacting %s", conn.RemoteAddr().String())
	client := newhubClient(conn.RemoteAddr().String(), toName)
	handShake, _ := cryptor.encryptB64(fmt.Sprintf("%s|%s|SERV", conf.Name, conf.TCPaddr))
	mess := newDatamessage(nil, client.CType, client, append([]byte("[HELO]"), handShake...))
	zehub.Unicast <- mess

	go tcpWriter(conn, client)
	(*wg).Done()
	tcpReader(conn, client)
	zehub.Unregister <- client
	// <-client.Consistent
}

func newIncommingTCPConn(conn *net.TCPConn, wg *sync.WaitGroup) {
	client := newhubClient(conn.RemoteAddr().String(), conn.RemoteAddr().String())
	handShake, _ := cryptor.encryptB64(fmt.Sprintf("%s|%s|SERV", conf.Name, conf.TCPaddr))
	mess := newDatamessage(nil, client.CType, client, append([]byte("[HELO]"), handShake...))
	zehub.Unicast <- mess

	go tcpWriter(conn, client)
	(*wg).Done()
	tcpReader(conn, client)
	zehub.Unregister <- client
	// <-client.Consistent
}

// TCPStart lance le server TCP
func tcpStart() {
	var wg sync.WaitGroup

	formatedaddr, _ := net.ResolveTCPAddr("tcp", conf.TCPaddr)
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
		go newIncommingTCPConn(conn, &wg)
		wg.Wait()
	}
}
