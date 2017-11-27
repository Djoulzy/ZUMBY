package world

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Djoulzy/Tools/clog"
)

type AOI struct {
	X            int
	Y            int
	Adjacents    []*AOI
	EntitiesList map[string]interface{}
	Actions      [][]byte
	update       []byte
}

type AOIList [][]*AOI

func BuildAOIList(W *WORLD) *AOIList {
	width := W.Map.Width / AOIWidth
	height := W.Map.Height / AOIHeight
	clog.Info("AOI", "BuildAOIList", "Map Size: %dx%d -> AOIList: %dx%d", W.Map.Width, W.Map.Height, width, height)

	var L AOIList
	L = make(AOIList, width)
	for i := 0; i < width; i++ {
		L[i] = make([]*AOI, height)
		for j := 0; j < height; j++ {
			L[i][j] = &AOI{
				X:            i,
				Y:            j,
				EntitiesList: make(map[string]interface{}),
			}
		}
	}

	for i, cols := range L {
		for j, rows := range cols {
			for x := -1; x < 2; x++ {
				for y := -1; y < 2; y++ {
					dx := i + x
					dy := j + y
					if (x == 0 && y == 0) || dx < 0 || dy < 0 || dx >= width || dy >= height {
						continue
					}
					// clog.Trace("", "", "Adding %dx%d (%s) to %dx%d", dx, dy, L[dx][dy], i, j)
					rows.Adjacents = append(rows.Adjacents, L[dx][dy])
				}
			}
			// clog.Trace("", "", "%#v\n", rows.Adjacents)
		}
	}

	return &L
}

func (aoi AOI) String() string {
	return fmt.Sprintf("%dx%d", aoi.X, aoi.Y)
}

func (L AOIList) String() string {
	var str string
	for _, cols := range L {
		for _, rows := range cols {
			str = fmt.Sprintf("%s %s", str, rows)
		}
	}
	return str
}

func (aoi *AOI) addEvent(mess []byte) {
	for _, adj := range aoi.Adjacents {
		adj.Actions = append(adj.Actions, mess)
	}
	aoi.Actions = append(aoi.Actions, mess)
}

func (aoi *AOI) addEntity(entity interface{}) {
	if typedEnt, ok := entity.(*MOB); ok {
		typedEnt.AOI = aoi
		aoi.EntitiesList[typedEnt.ID] = entity
	} else {
		if typedEnt, ok := entity.(*USER); ok {
			typedEnt.AOI = aoi
			aoi.EntitiesList[typedEnt.ID] = entity
		}
	}
}

// func (L *AOIList) addItemsToAOI(items [][]ITEM) {
// 	clog.Info("AOI", "addItemsToAOI", "Adding items to AOIs")
// 	for _, cols := range items {
// 		for _, item := range cols {
// 			aoi := L.getAOIfromCoord(item.X, item.Y)
// 			aoi.Items = append(aoi.Items, item)
// 		}
// 	}
// }

func (L *AOIList) moveEntity(x, y int, entity interface{}) {
	aoi := L.getAOIfromCoord(x, y)
	if typedEnt, ok := entity.(*MOB); ok {
		if typedEnt.AOI != aoi {
			typedEnt.AOI = nil
			aoi.EntitiesList[typedEnt.ID] = nil
			aoi.addEntity(entity)
		}
		typedEnt.waitState = typedEnt.Speed
	} else {
		if typedEnt, ok := entity.(*USER); ok {
			if typedEnt.AOI != aoi {
				typedEnt.AOI = nil
				aoi.EntitiesList[typedEnt.ID] = nil
				aoi.addEntity(entity)
			}
		}
	}
	json, _ := json.Marshal(entity)
	message := []byte(fmt.Sprintf("[BCST]%s", json))
	aoi.addEvent(message)
}

func (L *AOIList) addEntity(x, y int, entity interface{}) {
	aoi := L.getAOIfromCoord(x, y)
	aoi.addEntity(entity)
	json, _ := json.Marshal(entity)
	message := []byte(fmt.Sprintf("[NENT]%s", json))
	aoi.addEvent(message)
}

func (L *AOIList) dropEntity(x, y int, entity interface{}) {
	var message []byte
	aoi := L.getAOIfromCoord(x, y)
	if typedEnt, ok := entity.(*MOB); ok {
		typedEnt.AOI = nil
		aoi.EntitiesList[typedEnt.ID] = nil
		message = []byte(fmt.Sprintf("[KILL]%s", typedEnt.ID))
	} else {
		if typedEnt, ok := entity.(*USER); ok {
			typedEnt.AOI = nil
			aoi.EntitiesList[typedEnt.ID] = nil
			message = []byte(fmt.Sprintf("[KILL]%s", typedEnt.ID))
		}
	}
	aoi.addEvent(message)
}

func (L *AOIList) addEvent(x, y int, mess []byte) {
	aoi := L.getAOIfromCoord(x, y)
	aoi.addEvent(mess)
}

func (L AOIList) computeUpdates() {
	for _, cols := range L {
		for _, aoi := range cols {
			aoi.update = bytes.Join(aoi.Actions, []byte("|"))
			aoi.Actions = aoi.Actions[:0]
		}
	}
}

func (L AOIList) getAOIfromCoord(x, y int) *AOI {
	AOIx := x / AOIWidth
	AOIy := y / AOIHeight
	return L[AOIx][AOIy]
}

func (L *AOIList) getUpdateForPlayer(x, y int) ([]byte, error) {
	aoi := L.getAOIfromCoord(x, y)
	if len(aoi.update) > 0 {
		return aoi.update, nil
	} else {
		return aoi.update, errors.New("No update")
	}
}
