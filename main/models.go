package main

import (
	"fmt"
	"time"

	"github.com/Djoulzy/Tools/cmap"
)

var (
	newLine   = []byte{'\r', '\n'}
	spaceChar = []byte{' '}
)

type entity struct {
	ID        string `bson:"id" json:"id"`
	Type      string `bson:"typ" json:"typ"`
	Face      string `bson:"png" json:"png"`
	ComID     int    `bson:"num" json:"num"`
	Dir       string `bson:"mov" json:"mov"`
	X         int    `bson:"x" json:"x"` // Col nums
	Y         int    `bson:"y" json:"y"` // Row nums
	Pow       int    `bson:"pow" json:"pow"`
	Speed     int    `bson:"spd" json:"spd"`
	aoi       *AOI
	waitState int
}

type attributes struct {
	PV     int `bson:"pv" json:"pv"`
	Starv  int `bson:"st" json:"st"`
	Thirst int `bson:"th" json:"th"`
	Fight  int `bson:"fgt" json:"fgt"`
	Shoot  int `bson:"sht" json:"sht"`
	Craft  int `bson:"cft" json:"cft"`
	Breed  int `bson:"brd" json:"brd"`
	Grow   int `bson:"grw" json:"grw"`
}

type item struct {
	ID    int    `bson:"id" json:"id"`
	Owner string `bson:"-" json:"-"`
}

type user struct {
	hubClient *hubClient
	entity
	attributes
	Inventory []item `bson:"i" json:"i"`
}

type mob struct {
	entity
}

type tile struct {
	ID    int    `bson:"id" json:"id"`
	Type  string `bson:"type" json:"type"`
	Item  bool   `bson:"item" json:"item"`
	Block bool   `bson:"block" json:"block"`
	Name  string `bson:"name" json:"name"`
}

type world struct {
	hub *hubManager
	// MobList   map[string]*MOB
	MobList *cmap.CMap
	// UserList  map[string]*USER
	UserList  *cmap.CMap
	TilesList []tile
	Map       *mapData
	AOIs      *AOIList
	TimeStep  time.Duration
	TileSize  int
	AOIWidth  int
	AOIHeight int
	MobSpeed  int
	MaxMobNum int
}

type inventory struct {
	Owner      string `bson:"owner" json:"owner"`
	ID         int    `bson:"id" json:"id"`
	FromPocket int    `bson:"fp" json:"fp"`
	ToPocket   int    `bson:"tp" json:"tp"`
	X          int    `bson:"x" json:"x"`
	Y          int    `bson:"y" json:"y"`
}

type chatmsg struct {
	From string `bson:"from" json:"from"`
	Type int    `bson:"type" json:"type"`
	Mess string `bson:"mess" json:"mess"`
}

func (E entity) String() string {
	return fmt.Sprintf("\nID: %s [%s - %s]\nCoord: %dx%d - %s\nAOI: %s", E.ID, E.Type, E.Face, E.X, E.Y, E.Dir, E.aoi)
}
