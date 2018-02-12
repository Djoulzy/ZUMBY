package main

import (
	"fmt"
	"strings"

	"github.com/Djoulzy/Tools/clog"
)

func welcomeNewMonitor(c *hubClient, newName string, appID string) {
	if len(zehubManager.Monitors) >= conf.MaxMonitorsConns {
		zehubManager.Unregister <- c
		// <-c.Consistent
	} else {
		zehubManager.newRole(&connModifier{hubClient: c, NewName: c.Name, NewType: clientMonitor})
		c.AppID = appID
		clog.Info("server", "welcomeNewMonitor", "Accepting %s", c.Name)
	}
}

func welcomeNewUser(c *hubClient, newName string, appID string) {
	if zehubManager.userExists(c.Name, clientUndefined) {
		if len(zehubManager.Users) >= conf.MaxUsersConns && !zehubManager.userExists(newName, clientUser) {
			clog.Warn("CallToAction", "welcomeNewUser", "Too many Users connections, rejecting %s (In:%d/Cl:%d).", c.Name, len(zehubManager.Incomming), len(zehubManager.Users))
			if !scaleList.redirectConnection(c) {
				clog.Error("CallToAction", "welcomeNewUser", "NO FREE SLOTS !!!")
			}
			zehubManager.Unregister <- c
			// <-c.Consistent
		} else {
			clog.Info("CallToAction", "welcomeNewUser", "Identifying %s as %s", c.Name, newName)
			zehubManager.newRole(&connModifier{hubClient: c, NewName: newName, NewType: clientUser})
			c.AppID = appID
			scaleList.dispatchNewConnection(zehubManager, c.Name)

			// message := []byte(fmt.Sprintf("[NUSR]%s", c.Name))
			// mess := hub.NewdataMessage(c, hub.clientUser, nil, message)
			// zehubManager.Broadcast <- mess
			infos, err := zeWorld.logUser(c)
			if err != nil {
				zehubManager.Unregister <- c
			} else {
				message := []byte(fmt.Sprintf("[WLCM]%s", infos))
				mess := newDatamessage(nil, clientUser, c, message)
				zehubManager.Unicast <- mess
			}
		}
	} else {
		clog.Warn("CallToAction", "welcomeNewUser", "Can't identify client... Disconnecting %s / %s.", c.Name, newName)
		zehubManager.Unregister <- c
		// <-c.Consistent
	}
}

func welcomeNewServer(c *hubClient, newName string, addr string) {
	if len(zehubManager.Servers) >= conf.MaxServersConns {
		clog.Warn("server", "welcomeNewServer", "Too many Server connections, rejecting %s (In:%d/Cl:%d).", c.Name, len(zehubManager.Incomming), len(zehubManager.Servers))
		zehubManager.Unregister <- c
		// <-c.Consistent
		return
	}

	if zehubManager.userExists(c.Name, clientUndefined) {
		clog.Info("server", "welcomeNewServer", "Identifying %s as %s", c.Name, newName)
		zehubManager.newRole(&connModifier{hubClient: c, NewName: newName, NewType: clientServer})
		c.Addr = addr
		scaleList.addNewConnectedServer(c)
	} else {
		clog.Warn("server", "welcomeNewServer", "Can't identify server... Disconnecting %s.", c.Name)
		zehubManager.Unregister <- c
		// <-c.Consistent
	}
}

// HandShake valide la connexion entre le client et le serveur
func handShake(c *hubClient, message []byte) {
	uncryptedMessage, _ := cryptor.decryptB64(string(message))
	clog.Info("server", "HandShake", "New Incomming hubClient %s (%s)", c.Name, uncryptedMessage)
	infos := strings.Split(string(uncryptedMessage), "|")
	if len(infos) != 3 {
		clog.Warn("server", "HandShake", "Bad Handshake format ... Disconnecting")
		zehubManager.Unregister <- c
		// <-c.Consistent
		return
	}

	AppID := strings.TrimSpace(infos[1])
	newName := strings.TrimSpace(infos[0])
	switch infos[2] {
	case "MNTR":
		welcomeNewMonitor(c, newName, AppID)
	case "SERV":
		welcomeNewServer(c, newName, AppID)
	case "USER":
		welcomeNewUser(c, newName, AppID)
	default:
		zehubManager.Unregister <- c
		// <-c.Consistent
	}
}

// CallToAction Analyse la requette et appelle l'action demandée
func CallToAction(c *hubClient, message []byte) {
	if len(message) < 6 {
		clog.Warn("server", "CallToAction", "Bad Command '%s', disconnecting client %s.", message, c.Name)
		zehubManager.Unregister <- c
		// <-c.Consistent
		return
	}

	cmdGroup := string(message[0:6])
	actionGroup := message[6:]

	if c.CType != clientUndefined {
		switch cmdGroup {
		case "[BCST]":
			// clog.Trace("", "", "%s", message)
			mess := newDatamessage(c, clientUser, nil, message)
			zehubManager.Broadcast <- mess
			if c.CType != clientServer {
				mess = newDatamessage(c, clientServer, nil, message)
				zehubManager.Broadcast <- mess
			}
		// case "[UCST]":
		// case "[STOR]":
		// 	Storage.NewRecord(string(actionGroup))
		case "[QUIT]":
			zehubManager.Unregister <- c
			// <-c.Consistent
		case "[MNIT]":
			clog.Debug("server", "CallToAction", "Metrics received from %s (%s)", c.Name, c.Addr)
			scaleList.updateMetrics(c.Addr, actionGroup)
		case "[KILL]":
			id := string(actionGroup)
			if zehubManager.userExists(id, clientUser) {
				userToKill := zehubManager.Users[id]
				clog.Info("server", "CallToAction", "Killing user %s", actionGroup)
				zehubManager.Unregister <- userToKill
				// <-userToKill.Consistent
			}
		case "[GKEY]":
			crypted, _ := cryptor.encryptB64(string(actionGroup))
			mess := newDatamessage(nil, c.CType, c, crypted)
			zehubManager.Unicast <- mess
		default:
			// mess := hub.NewdataMessage(nil, c.CType, c, []byte(fmt.Sprintf("%s:?", cmd_group)))
			// zehubManager.Unicast <- mess
			zeWorld.callToAction(c, cmdGroup, actionGroup)
		}
	} else {
		switch cmdGroup {
		case "[HELO]":
			// [HELO]<unique_id>|<app_id ou addr_ip>|<client_type>
			handShake(c, actionGroup)
		default:
			clog.Warn("server", "CallToAction", "Bad Command '%s', disconnecting client %s.", cmdGroup, c.Name)
			zehubManager.Unregister <- c
			// <-c.Consistent
		}
	}
}
