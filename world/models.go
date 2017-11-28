package world

import (
	"fmt"
	"time"

	"github.com/Djoulzy/ZUMBY/hub"
)

const (
	timeStep  = 100 * time.Millisecond // Actualisation 10 par seconde
	tileSize  = 32
	AOIWidth  = 30
	AOIHeight = 30
	mobSpeed  = 8
	maxMobNum = 5
)

type Entity struct {
	ID        string `bson:"id" json:"id"`
	Type      string `bson:"typ" json:"typ"`
	Face      string `bson:"png" json:"png"`
	ComID     int    `bson:"num" json:"num"`
	Dir       string `bson:"mov" json:"mov"`
	X         int    `bson:"x" json:"x"` // Col nums
	Y         int    `bson:"y" json:"y"` // Row nums
	Pow       int    `bson:"pow" json:"pow"`
	Speed     int    `bson:"spd" json:"spd"`
	AOI       *AOI   `bson:"-" json:"-"`
	waitState int    `bson:"-" json:"-"`
}

type Attributes struct {
	PV     int `bson:"pv" json:"pv"`
	Starv  int `bson:"st" json:"st"`
	Thirst int `bson:"th" json:"th"`
	Fight  int `bson:"fgt" json:"fgt"`
	Shoot  int `bson:"sht" json:"sht"`
	Craft  int `bson:"cft" json:"cft"`
	Breed  int `bson:"brd" json:"brd"`
	Grow   int `bson:"grw" json:"grw"`
}

type ITEM struct {
	ID int `bson:"id" json:"id"`
	X  int `bson:"x" json:"x"` // Col nums
	Y  int `bson:"y" json:"y"` // Row nums
}

type USER struct {
	hubClient *hub.Client
	Entity
	Attributes
	Inventory []ITEM `bson:"i" json:"i"`
}

type MOB struct {
	Entity
}

type TILE struct {
	Type int
	ID   string
}

type WORLD struct {
	hub      *hub.Hub
	MobList  map[string]*MOB
	UserList map[string]*USER
	Map      *MapData
	AOIs     *AOIList
}

func (E Entity) String() string {
	return fmt.Sprintf("\nID: %s [%s - %s]\nCoord: %dx%d - %s\nAOI: %s", E.ID, E.Type, E.Face, E.X, E.Y, E.Dir, E.AOI)
}
