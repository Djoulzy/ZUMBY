package world

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/Djoulzy/Tools/clog"
)

type FILELAYER struct {
	Data    []int  `bson:"data" json:"data"`
	Name    string `bson:"name" json:"name"`
	Width   int    `bson:"width" json:"width"`
	Height  int    `bson:"height" json:"height"`
	Opacity int    `bson:"opacity" json:"opacity"`
	Type    string `bson:"type" json:"type"`
	Visible bool   `bson:"visible" json:"visible"`
	X       int    `bson:"x" json:"x"`
	Y       int    `bson:"y" json:"y"`
	Offsetx int    `bson:"offsetx" json:"offsetx"`
	Offsety int    `bson:"offsety" json:"offsety"`
}

type FILETILESET struct {
	Columns     int    `bson:"columns" json:"columns"`
	Firstgid    int    `bson:"firstgid" json:"firstgid"`
	Image       string `bson:"image" json:"image"`
	Imageheight int    `bson:"imageheight" json:"imageheight"`
	Imagewidth  int    `bson:"imagewidth" json:"imagewidth"`
	Margin      int    `bson:"margin" json:"margin"`
	Name        string `bson:"name" json:"name"`
	Spacing     int    `bson:"spacing" json:"spacing"`
	Tilecount   int    `bson:"tilecount" json:"tilecount"`
	Tileheight  int    `bson:"tileheight" json:"tileheight"`
	Tilewidth   int    `bson:"tilewidth" json:"tilewidth"`
}

type FILEMAP struct {
	Width        int           `bson:"width" json:"width"`
	Height       int           `bson:"height" json:"height"`
	Nextobjectid int           `bson:"nextobjectid" json:"nextobjectid"`
	Orientation  string        `bson:"orientation" json:"orientation"`
	Renderorder  string        `bson:"renderorder" json:"renderorder"`
	Tiledversion string        `bson:"tiledversion" json:"tiledversion"`
	Tileheight   int           `bson:"tileheight" json:"tileheight"`
	Tilewidth    int           `bson:"tilewidth" json:"tilewidth"`
	Type         string        `bson:"type" json:"type"`
	Version      int           `bson:"version" json:"version"`
	Layers       []FILELAYER   `bson:"layers" json:"layers"`
	Tilesets     []FILETILESET `bson:"tilesets" json:"tilesets"`
}

type MapData struct {
	Width    int
	Height   int
	Entities [][]interface{}
	Over     [][]int
	Block    [][]int
	Ground   [][]int
	Items    [][]ITEM
	FileData FILEMAP
}

func (M *MapData) loadTiledJSONMap(file string) {
	dat, _ := ioutil.ReadFile(file)
	err := json.Unmarshal(dat, &M.FileData)
	if err != nil {
		clog.Error("", "", "%s", err)
	}

	// M.FileData := mapper.NewMap()
	M.Width = M.FileData.Layers[0].Width
	M.Height = M.FileData.Layers[0].Height

	M.Ground = make([][]int, M.Width)
	M.Over = make([][]int, M.Width)
	M.Block = make([][]int, M.Width)
	M.Items = make([][]ITEM, M.Width)
	M.Entities = make([][]interface{}, M.Width)
	for i := 0; i < M.Width; i++ {
		M.Entities[i] = make([]interface{}, M.Height)
		M.Ground[i] = make([]int, M.Height)
		M.Block[i] = make([]int, M.Height)
		M.Over[i] = make([]int, M.Height)
		M.Items[i] = make([]ITEM, M.Height)
	}

	row := 0
	for row < M.Height {
		col := 0
		for col < M.Width {
			M.Ground[col][row] = M.FileData.Layers[0].Data[(row*M.Width)+col]
			M.Block[col][row] = M.FileData.Layers[1].Data[(row*M.Width)+col]
			M.Over[col][row] = M.FileData.Layers[2].Data[(row*M.Width)+col]
			item := M.FileData.Layers[3].Data[(row*M.Width)+col]
			if item != 0 {
				M.Items[col][row] = ITEM{
					ID: M.FileData.Layers[3].Data[(row*M.Width)+col],
					X:  col,
					Y:  row,
				}
				M.Over[col][row] = item
			}
			M.Entities[col][row] = nil
			col++
		}
		row++
	}
}

func (M *MapData) ExportMapArea(x, y, AOIWidth, AOIHeight int) []byte {
	var startx, starty int
	tmp := M.FileData

	if x-1 < 0 {
		startx = 0
		tmp.Width = AOIWidth * 2
	} else {
		startx = (x - 1) * AOIWidth
		tmp.Width = AOIWidth * 3
	}

	if y-1 < 0 {
		starty = 0
		tmp.Height = AOIHeight * 2
	} else {
		starty = (y - 1) * AOIHeight
		tmp.Height = AOIHeight * 3
	}

	tmp.Layers = make([]FILELAYER, 3)
	tmp.Layers[0].Data = make([]int, tmp.Width*tmp.Height)
	tmp.Layers[1].Data = make([]int, tmp.Width*tmp.Height)
	tmp.Layers[2].Data = make([]int, tmp.Width*tmp.Height)

	cpt := 0
	for j := starty; j < starty+tmp.Height; j++ {
		for i := startx; i < startx+tmp.Width; i++ {
			tmp.Layers[0].Data[cpt] = M.Ground[i][j]
			tmp.Layers[1].Data[cpt] = M.Block[i][j]
			if M.Items[i][j].ID != 0 {
				tmp.Layers[2].Data[cpt] = M.Items[i][j].ID
			} else {
				tmp.Layers[2].Data[cpt] = M.Over[i][j]
			}
			cpt++
		}
	}

	tmp.Layers[0].Name = "terrain"
	tmp.Layers[0].Type = "tilelayer"
	tmp.Layers[0].Opacity = 1
	tmp.Layers[0].Visible = true
	tmp.Layers[0].Width = tmp.Width
	tmp.Layers[0].Height = tmp.Height
	// tmp.Layers[0].X = startx * tmp.Tilewidth
	// tmp.Layers[0].Y = starty * tmp.Tileheight
	// tmp.Layers[0].Offsetx = 0
	// tmp.Layers[0].Offsety = 0
	tmp.Layers[0].X = 0
	tmp.Layers[0].Y = 0
	tmp.Layers[0].Offsetx = startx * tmp.Tilewidth
	tmp.Layers[0].Offsety = starty * tmp.Tileheight

	tmp.Layers[1].Name = "obstacles"
	tmp.Layers[1].Type = "tilelayer"
	tmp.Layers[1].Opacity = 1
	tmp.Layers[1].Visible = true
	tmp.Layers[1].Width = tmp.Width
	tmp.Layers[1].Height = tmp.Height
	// tmp.Layers[1].X = startx * tmp.Tilewidth
	// tmp.Layers[1].Y = starty * tmp.Tileheight
	// tmp.Layers[1].Offsetx = 0
	// tmp.Layers[1].Offsety = 0
	tmp.Layers[1].X = 0
	tmp.Layers[1].Y = 0
	tmp.Layers[1].Offsetx = startx * tmp.Tilewidth
	tmp.Layers[1].Offsety = starty * tmp.Tileheight

	tmp.Layers[2].Name = "hauteurs"
	tmp.Layers[2].Type = "tilelayer"
	tmp.Layers[2].Opacity = 1
	tmp.Layers[2].Visible = true
	tmp.Layers[2].Width = tmp.Width
	tmp.Layers[2].Height = tmp.Height
	// tmp.Layers[2].X = startx * tmp.Tilewidth
	// tmp.Layers[2].Y = starty * tmp.Tileheight
	// tmp.Layers[2].Offsetx = 0
	// tmp.Layers[2].Offsety = 0
	tmp.Layers[2].X = 0
	tmp.Layers[2].Y = 0
	tmp.Layers[2].Offsetx = startx * tmp.Tilewidth
	tmp.Layers[2].Offsety = starty * tmp.Tileheight

	json, _ := json.MarshalIndent(tmp, "", "    ")
	return json
}

func (M *MapData) Draw() {
	fmt.Printf("%c[H", 27)
	visuel := ""
	display := "*"
	for y := 0; y < M.Height; y++ {
		for x := 0; x < M.Width; x++ {
			val := M.Block[x][y]
			if val == 0 {
				visuel = "   "
			} else if val == -1 {
				visuel = clog.GetColoredString(" + ", "black", "green")
			} else if val == 1000 {
				visuel = clog.GetColoredString(" D ", "black", "yellow")
			} else if val == 2000 {
				visuel = clog.GetColoredString(" F ", "white", "blue")
			} else {
				visuel = clog.GetColoredString(" X ", "white", "white")
			}
			if M.Entities[x][y] != nil {
				switch M.Entities[x][y].(type) {
				case *MOB:
					visuel = clog.GetColoredString(" Z ", "white", "red")
				case *USER:
					visuel = clog.GetColoredString(" P ", "black", "yellow")
				}
			}
			display = fmt.Sprintf("%s%s", display, visuel)
		}
		display = fmt.Sprintf("%s*\n*", display)
	}
	fmt.Printf("%s", display)
}

func (M *MapData) buildMonPage(AOIWidth, AOIHeight int) {
	var area string = "<!DOCTYPE html><html ><head><meta charset='UTF-8'><title>Mon</title></head><body>"
	for y := 0; y < (M.Height / AOIHeight); y++ {
		for x := 0; x < (M.Width / AOIWidth); x++ {
			area = fmt.Sprintf("%s\n<AREA shape='rect' coords='%d,%d,%d,%d' href='assets/%d_%d.png'>", area, x*AOIWidth*32, y*AOIHeight*32, (x+1)*AOIWidth*32, (y+1)*AOIHeight*32, x, y)
		}
	}
	area = fmt.Sprintf("%s\n<img USEMAP='#map' src='assets/mon.png' /></body></html>", area)
	f, _ := os.OpenFile("../public/mon.html", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	f.WriteString(area)
}

func drawRect(img *image.RGBA, x, y, width, height int, c color.Color) {
	for ry := 0; ry < height; ry++ {
		for rx := 0; rx < width; rx++ {
			img.Set(x+rx, y+ry, c)
		}
	}
}

func (M *MapData) genAOI(x, y, AOIWidth, AOIHeight int) {
	pixel := 11
	img := image.NewRGBA(image.Rect(0, 0, AOIWidth*pixel, AOIHeight*pixel))

	startx := x * AOIWidth
	starty := y * AOIHeight

	for aoiy := 0; aoiy < AOIHeight; aoiy++ {
		for aoix := 0; aoix < AOIWidth; aoix++ {
			val := M.Block[startx+aoix][starty+aoiy]
			if val == 0 {
				drawRect(img, aoix*pixel, aoiy*pixel, pixel, pixel, color.RGBA{0, 0, 0, 255})
			} else {
				drawRect(img, aoix*pixel, aoiy*pixel, pixel, pixel, color.RGBA{255, 255, 255, 255})
			}
			if M.Entities[startx+aoix][starty+aoiy] != nil {
				switch M.Entities[startx+aoix][starty+aoiy].(type) {
				case *MOB:
					drawRect(img, aoix*pixel, aoiy*pixel, pixel, pixel, color.RGBA{255, 0, 0, 255})
				case *USER:
					drawRect(img, aoix*pixel, aoiy*pixel, pixel, pixel, color.RGBA{0, 255, 0, 255})
				}
			}
		}
	}

	// Save to out.png
	f, _ := os.OpenFile("../public/assets/0_0.png", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)
}

func (M *MapData) genImage() {
	img := image.NewRGBA(image.Rect(0, 0, M.Width, M.Height))

	for y := 0; y < M.Height; y++ {
		for x := 0; x < M.Width; x++ {
			val := M.Block[x][y]
			if val == 0 {
				img.Set(x, y, color.RGBA{0, 0, 0, 255})
			} else {
				img.Set(x, y, color.RGBA{255, 255, 255, 255})
			}
			if M.Entities[x][y] != nil {
				switch M.Entities[x][y].(type) {
				case *MOB:
					img.Set(x, y, color.RGBA{255, 0, 0, 255})
				case *USER:
					img.Set(x, y, color.RGBA{0, 255, 0, 255})
				}
			}
		}
	}

	// Save to out.png
	f, _ := os.OpenFile("../public/assets/mon.png", os.O_WRONLY|os.O_CREATE, 0600)
	defer f.Close()
	png.Encode(f, img)
}
