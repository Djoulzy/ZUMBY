package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/Djoulzy/Tools/clog"
	"github.com/stretchr/testify/assert"
)

var tmphubManager *hubManager

func newhubClient(name string, userType int) *hubClient {
	tmphubClient := &hubClient{
		Quit:  make(chan bool, 8),
		CType: userType, Send: make(chan []byte, 256),
		CallToAction: nil, Addr: "127.0.0.1:8080",
		Name: name, Content_id: 0, Front_id: "", App_id: "", Country: "", User_agent: "Test Socket",
	}
	return tmphubClient
}

func TestRegister(t *testing.T) {
	tmphubClient1 := newhubClient("TestRegister", clientUser)
	tmphubClient2 := newhubClient("TestRegister", clientUser)

	tmphubManager.Register <- tmphubClient1
	assert.Equal(t, true, tmphubManager.UserExists(tmphubClient1.Name, tmphubClient1.CType), "hubClient should be found")
	assert.Equal(t, tmphubClient1, tmphubManager.GethubClientByName(tmphubClient1.Name, tmphubClient1.CType), "Registered hubClient should equal original hubClient")

	tmphubManager.Register <- tmphubClient2

	assert.Equal(t, true, tmphubManager.UserExists(tmphubClient1.Name, tmphubClient1.CType), "hubClient should be found")
	assert.NotEqual(t, tmphubClient1, tmphubManager.GethubClientByName(tmphubClient1.Name, tmphubClient1.CType), "hubClient should be replaced")
	assert.Equal(t, tmphubClient2, tmphubManager.GethubClientByName(tmphubClient1.Name, tmphubClient1.CType), "Registered hubClient should equal to second client")
}

func TestUnregister(t *testing.T) {
	tmphubClient := newhubClient("test", clientUser)

	tmphubManager.Register <- tmphubClient
	tmphubManager.Unregister <- tmphubClient

	assert.Nil(t, tmphubManager.GethubClientByName(tmphubClient.Name, tmphubClient.CType))

	tmphubClient3 := newhubClient("TestRegister", clientServer)
	tmphubManager.Register <- tmphubClient3
	assert.Equal(t, tmphubClient3, tmphubManager.GethubClientByName(tmphubClient3.Name, tmphubClient3.CType), "Registered hubClient should equal original hubClient")
	tmphubManager.Register <- tmphubClient3
}

func TestConcurrency(t *testing.T) {
	var tmphubClient *hubClient

	for i := 0; i < 100; i++ {
		tmphubClient = newhubClient(fmt.Sprintf("%d", i), clientUser)

		tmphubManager.Register <- tmphubClient
		assert.Equal(t, tmphubClient, tmphubManager.GethubClientByName(tmphubClient.Name, tmphubClient.CType), "Registered hubClient should equal original hubClient")
		tmphubManager.Unregister <- tmphubClient
		assert.Nil(t, tmphubManager.GethubClientByName(tmphubClient.Name, tmphubClient.CType))
	}
}

func TestdataMessages(t *testing.T) {
	var tmphubClient *hubClient

	for i := 0; i < 10; i++ {
		tmphubClient = newhubClient(fmt.Sprintf("%d", i), clientUser)

		tmphubManager.Register <- tmphubClient
	}

	mess := NewdataMessage(tmphubClient, clientUser, nil, []byte("BROADCAST"))
	tmphubManager.Broadcast <- mess

	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("%d", i)
		client := tmphubManager.GethubClientByName(name, clientUser)

		message, ok := <-client.Send
		if ok {
			assert.Equal(t, "BROADCAST", string(message), "dataMessage cannot be read from channel")
		} else {
			t.Fail()
		}
	}

	mess = NewdataMessage(nil, clientUser, tmphubClient, []byte("UNICAST"))
	tmphubManager.Unicast <- mess

	message, ok := <-tmphubClient.Send
	if ok {
		assert.Equal(t, "UNICAST", string(message), "dataMessage cannot be read from channel")
	} else {
		t.Fail()
	}
}

func TestNewRole(t *testing.T) {
	tmphubClient := newhubClient("0", clientUser)
	tmphubManager.Register <- tmphubClient

	newRole := &ConnModifier{
		hubClient: tmphubClient,
		NewName:   "1",
		NewType:   clientUser,
	}

	tmphubManager.Newrole(newRole)
	assert.Equal(t, tmphubClient, tmphubManager.GethubClientByName("1", clientUser), "Bad new Role")
	assert.Nil(t, tmphubManager.GethubClientByName("0", clientUser))
}

func TestMain(m *testing.M) {
	clog.LogLevel = 5
	clog.StartLogging = false

	tmphubManager = NewhubManager()
	go tmphubManager.Run()
	os.Exit(m.Run())
}
