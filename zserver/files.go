package zserver

import (
	"fmt"
	"io/ioutil"

	"github.com/Djoulzy/Tools/clog"
)

const (
	fileDBPath = "../../FileDB"
)

// SaveUser flush les data joueur dans un fichier texte
func SaveUser(ID string, json []byte) {
	fileName := fmt.Sprintf("%s/users/%s.txt", fileDBPath, ID)
	err := ioutil.WriteFile(fileName, json, 0644)
	if err != nil {
		clog.Error("storage", "SaveUser", "Can't save User %s : %s", ID, err)
	}
}

// LoadUser Charge les data joueur
func LoadUser(ID string) ([]byte, error) {
	fileName := fmt.Sprintf("%s/users/%s.txt", fileDBPath, ID)
	json, err := ioutil.ReadFile(fileName)
	if err != nil {
		clog.Error("storage", "SaveUser", "Can't read User %s : %s", ID, err)
	}
	return json, err
}
