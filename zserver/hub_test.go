package zserver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newTesthubClient(name string, userType int) *hubClient {
	tmphubClient := &hubClient{
		Quit:  make(chan bool, 8),
		CType: userType, Send: make(chan []byte, 256),
		Addr: "127.0.0.1:8080",
		Name: name, AppID: "", Country: "", UserAgent: "Test Socket",
	}
	return tmphubClient
}

func TestRegister(t *testing.T) {
	tmphubClient1 := newTesthubClient("TestRegister", clientUser)
	tmphubClient2 := newTesthubClient("TestRegister", clientUser)

	// Enregistrement d'un user
	tmphubManager.register(tmphubClient1)
	assert.Equal(t, true, tmphubManager.userExists(tmphubClient1.Name, tmphubClient1.CType), "hubClient should be found")
	assert.Equal(t, tmphubClient1, tmphubManager.gethubClientByName(tmphubClient1.Name, tmphubClient1.CType), "Registered hubClient should equal original hubClient")

	// Enregistrement du même user
	tmphubManager.register(tmphubClient2)
	assert.Equal(t, true, tmphubManager.userExists(tmphubClient1.Name, tmphubClient1.CType), "hubClient should be found")
	// Le user existe déjà, l'ancien doit être remplacé par le nouveau
	assert.NotEqual(t, tmphubClient1, tmphubManager.gethubClientByName(tmphubClient1.Name, tmphubClient1.CType), "hubClient should be replaced")
	assert.Equal(t, tmphubClient2, tmphubManager.gethubClientByName(tmphubClient1.Name, tmphubClient1.CType), "Registered hubClient should equal to second client")
}

func TestUnregister(t *testing.T) {
	tmphubClient := newTesthubClient("test", clientUser)

	tmphubManager.register(tmphubClient)
	tmphubManager.unregister(tmphubClient)

	assert.Nil(t, tmphubManager.gethubClientByName(tmphubClient.Name, tmphubClient.CType))

	tmphubClient3 := newTesthubClient("TestRegister", clientServer)
	tmphubManager.register(tmphubClient3)
	assert.Equal(t, tmphubClient3, tmphubManager.gethubClientByName(tmphubClient3.Name, tmphubClient3.CType), "Registered hubClient should equal original hubClient")
	tmphubManager.register(tmphubClient3)
}

func TestConcurrency(t *testing.T) {
	var tmphubClient *hubClient

	for i := 0; i < 100; i++ {
		tmphubClient = newTesthubClient(fmt.Sprintf("%d", i), clientUser)

		tmphubManager.register(tmphubClient)
		assert.Equal(t, tmphubClient, tmphubManager.gethubClientByName(tmphubClient.Name, tmphubClient.CType), "Registered hubClient should equal original hubClient")
		tmphubManager.unregister(tmphubClient)
		assert.Nil(t, tmphubManager.gethubClientByName(tmphubClient.Name, tmphubClient.CType))
	}
}

func TestDataMessages(t *testing.T) {
	var tmphubClient *hubClient

	for i := 0; i < 10; i++ {
		tmphubClient = newTesthubClient(fmt.Sprintf("%d", i), clientUser)

		tmphubManager.register(tmphubClient)
	}

	mess := newDataMessage(tmphubClient, clientUser, nil, []byte("BROADCAST"))
	tmphubManager.Broadcast <- mess

	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("%d", i)
		client := tmphubManager.gethubClientByName(name, clientUser)

		message, ok := <-client.Send
		if ok {
			assert.Equal(t, "BROADCAST", string(message), "dataMessage cannot be read from channel")
		} else {
			t.Fail()
		}
	}

	mess = newDataMessage(nil, clientUser, tmphubClient, []byte("UNICAST"))
	tmphubManager.Unicast <- mess

	message, ok := <-tmphubClient.Send
	if ok {
		assert.Equal(t, "UNICAST", string(message), "dataMessage cannot be read from channel")
	} else {
		t.Fail()
	}
}

func TestNewRole(t *testing.T) {
	tmphubClient := newTesthubClient("0", clientUser)
	tmphubManager.register(tmphubClient)

	newRole := &connModifier{
		hubClient: tmphubClient,
		NewName:   "1",
		NewType:   clientUser,
	}

	tmphubManager.newRole(newRole)
	assert.Equal(t, tmphubClient, tmphubManager.gethubClientByName("1", clientUser), "Bad new Role")
	assert.Nil(t, tmphubManager.gethubClientByName("0", clientUser))
}
