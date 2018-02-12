package main

import (
	"os"
	"testing"

	"github.com/Djoulzy/Tools/clog"
	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	var cryptor = &cypher{
		HashSize: 8,
		HexKey:   []byte("d87fbb277eefe245ee384b6098637513462f5151336f345778706b462f724473"),
		HexIV:    []byte("046b51957f00c25929e8ccaad3bfe1a7"),
	}

	crypted, _ := cryptor.encryptB64("iphone1|xcode|USER")
	assert.Equal(t, "BGtRlX8Awlkp6Myq07_hpw/QvGzLgBaPZiJgeKdpfg7HZzBhEaspxOJaCBv-05d96k", string(crypted), "Bad encryption")
}

func TestMain(m *testing.M) {
	clog.LogLevel = 5
	clog.StartLogging = true

	os.Exit(m.Run())
}
