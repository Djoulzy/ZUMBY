package main

import (
	"fmt"
	"strings"

	"github.com/Djoulzy/Tools/clog"
)

func welcomeNewMonitor(c *Client, newName string, app_id string) {
	if len(zeHub.Monitors) >= conf.MaxMonitorsConns {
		zeHub.Unregister <- c
		// <-c.Consistent
	} else {
		zeHub.Newrole(&ConnModifier{Client: c, NewName: c.Name, NewType: ClientMonitor})
		c.App_id = app_id
		clog.Info("server", "welcomeNewMonitor", "Accepting %s", c.Name)
	}
}

func welcomeNewUser(c *Client, newName string, app_id string) {
	if zeHub.UserExists(c.Name, ClientUndefined) {
		if len(zeHub.Users) >= conf.MaxUsersConns && !zeHub.UserExists(newName, ClientUser) {
			clog.Warn("CallToAction", "welcomeNewUser", "Too many Users connections, rejecting %s (In:%d/Cl:%d).", c.Name, len(zeHub.Incomming), len(zeHub.Users))
			if !ScaleList.RedirectConnection(c) {
				clog.Error("CallToAction", "welcomeNewUser", "NO FREE SLOTS !!!")
			}
			zeHub.Unregister <- c
			// <-c.Consistent
		} else {
			clog.Info("CallToAction", "welcomeNewUser", "Identifying %s as %s", c.Name, newName)
			zeHub.Newrole(&ConnModifier{Client: c, NewName: newName, NewType: ClientUser})
			c.App_id = app_id
			ScaleList.DispatchNewConnection(zeHub, c.Name)

			// message := []byte(fmt.Sprintf("[NUSR]%s", c.Name))
			// mess := hub.NewMessage(c, hub.ClientUser, nil, message)
			// zeHub.Broadcast <- mess
			infos, err := zeWorld.LogUser(c)
			if err != nil {
				zeHub.Unregister <- c
			} else {
				message := []byte(fmt.Sprintf("[WLCM]%s", infos))
				mess := NewMessage(nil, ClientUser, c, message)
				zeHub.Unicast <- mess
			}
		}
	} else {
		clog.Warn("CallToAction", "welcomeNewUser", "Can't identify client... Disconnecting %s / %s.", c.Name, newName)
		zeHub.Unregister <- c
		// <-c.Consistent
	}
}

func welcomeNewServer(c *Client, newName string, addr string) {
	if len(zeHub.Servers) >= conf.MaxServersConns {
		clog.Warn("server", "welcomeNewServer", "Too many Server connections, rejecting %s (In:%d/Cl:%d).", c.Name, len(zeHub.Incomming), len(zeHub.Servers))
		zeHub.Unregister <- c
		// <-c.Consistent
		return
	}

	if zeHub.UserExists(c.Name, ClientUndefined) {
		clog.Info("server", "welcomeNewServer", "Identifying %s as %s", c.Name, newName)
		zeHub.Newrole(&ConnModifier{Client: c, NewName: newName, NewType: ClientServer})
		c.Addr = addr
		ScaleList.AddNewConnectedServer(c)
	} else {
		clog.Warn("server", "welcomeNewServer", "Can't identify server... Disconnecting %s.", c.Name)
		zeHub.Unregister <- c
		// <-c.Consistent
	}
}

func HandShake(c *Client, message []byte) {
	uncrypted_message, _ := Cryptor.Decrypt_b64(string(message))
	clog.Info("server", "HandShake", "New Incomming Client %s (%s)", c.Name, uncrypted_message)
	infos := strings.Split(string(uncrypted_message), "|")
	if len(infos) != 3 {
		clog.Warn("server", "HandShake", "Bad Handshake format ... Disconnecting")
		zeHub.Unregister <- c
		// <-c.Consistent
		return
	}

	App_id := strings.TrimSpace(infos[1])
	newName := strings.TrimSpace(infos[0])
	switch infos[2] {
	case "MNTR":
		welcomeNewMonitor(c, newName, App_id)
	case "SERV":
		welcomeNewServer(c, newName, App_id)
	case "USER":
		welcomeNewUser(c, newName, App_id)
	default:
		zeHub.Unregister <- c
		// <-c.Consistent
	}
}

func CallToAction(c *Client, message []byte) {
	if len(message) < 6 {
		clog.Warn("server", "CallToAction", "Bad Command '%s', disconnecting client %s.", message, c.Name)
		zeHub.Unregister <- c
		// <-c.Consistent
		return
	}

	cmd_group := string(message[0:6])
	action_group := message[6:]

	if c.CType != ClientUndefined {
		switch cmd_group {
		case "[BCST]":
			// clog.Trace("", "", "%s", message)
			mess := NewMessage(c, ClientUser, nil, message)
			zeHub.Broadcast <- mess
			if c.CType != ClientServer {
				mess = NewMessage(c, ClientServer, nil, message)
				zeHub.Broadcast <- mess
			}
		// case "[UCST]":
		// case "[STOR]":
		// 	Storage.NewRecord(string(action_group))
		case "[QUIT]":
			zeHub.Unregister <- c
			// <-c.Consistent
		case "[MNIT]":
			clog.Debug("server", "CallToAction", "Metrics received from %s (%s)", c.Name, c.Addr)
			ScaleList.UpdateMetrics(c.Addr, action_group)
		case "[KILL]":
			id := string(action_group)
			if zeHub.UserExists(id, ClientUser) {
				userToKill := zeHub.Users[id]
				clog.Info("server", "CallToAction", "Killing user %s", action_group)
				zeHub.Unregister <- userToKill
				// <-userToKill.Consistent
			}
		case "[GKEY]":
			crypted, _ := Cryptor.Encrypt_b64(string(action_group))
			mess := NewMessage(nil, c.CType, c, crypted)
			zeHub.Unicast <- mess
		default:
			// mess := hub.NewMessage(nil, c.CType, c, []byte(fmt.Sprintf("%s:?", cmd_group)))
			// zeHub.Unicast <- mess
			zeWorld.CallToAction(c, cmd_group, action_group)
		}
	} else {
		switch cmd_group {
		case "[HELO]":
			// [HELO]<unique_id>|<app_id ou addr_ip>|<client_type>
			HandShake(c, action_group)
		default:
			clog.Warn("server", "CallToAction", "Bad Command '%s', disconnecting client %s.", cmd_group, c.Name)
			zeHub.Unregister <- c
			// <-c.Consistent
		}
	}
}
