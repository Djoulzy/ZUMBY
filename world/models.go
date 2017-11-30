package world

import (
	"fmt"
	"time"

	"github.com/Djoulzy/ZUMBY/hub"
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
	ID    int    `bson:"id" json:"id"`
	Owner string `bson:"owner" json:"owner"`
	X     int    `bson:"x" json:"x"` // Col nums
	Y     int    `bson:"y" json:"y"` // Row nums
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
	hub       *hub.Hub
	MobList   map[string]*MOB
	UserList  map[string]*USER
	Map       *MapData
	AOIs      *AOIList
	TimeStep  time.Duration
	TileSize  int
	AOIWidth  int
	AOIHeight int
	MobSpeed  int
	MaxMobNum int
}

func (E Entity) String() string {
	return fmt.Sprintf("\nID: %s [%s - %s]\nCoord: %dx%d - %s\nAOI: %s", E.ID, E.Type, E.Face, E.X, E.Y, E.Dir, E.AOI)
}
