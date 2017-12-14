'use strict'

var Config = require('config')

class Area
{
	// Status :
	// 0: init
	// 1: prending / loading
	// 2: active
	// 3: disabled

	constructor(game) {
		this.game = game
        this.data = null
        this.coord = new Phaser.Point(0, 0)
        this.status = 0
		this.name = '0_0'
		this.back = null
		this.items = null
		this.front = null
    }

	setWorldBound() {
		var newWidth = (this.coord.x+3)*this.game.Properties.areaWidth*this.game.Properties.step
		var newHeight = (this.coord.y+3)*this.game.Properties.areaHeight*this.game.Properties.step
		if (this.game.world.width > newWidth) newWidth = this.game.world.width
		if (this.game.world.height > newHeight) newHeight = this.game.world.height
		this.game.world.setBounds(0, 0, newWidth, newHeight)
		console.log("Area "+this.name+" Rendered - New World bounds : "+this.game.world.width+"x"+this.game.world.height)
	}

	render() {
		if (this.status == 1) {
			this.setWorldBound()
			this.data = this.game.add.tilemap(this.name);
			this.data.addTilesetImage('final');
			this.back = this.data.createLayer('back')
			this.items = this.data.createLayer('items')
			this.front = this.data.createLayer('front')

			this.game.backLayer.add(this.back)
			this.game.backLayer.add(this.items)
			this.game.frontLayer.add(this.front)
			this.game.frontLayer.sendToBack(this.front)

			// this.game.frontLayer.sort();
			this.status = 2
		}
	}

	getTileValueAt(x, y, layerType) {
		var result = 0
		var newX = x - (layerType.layer.offsetX/32) // - this.coord.x*this.game.Properties.areaWidth
		var newY = y - (layerType.layer.offsetY/32) // - this.coord.y*this.game.Properties.areaHeight
		var tmp = this.data.getTile(newX, newY, layerType)
		if (tmp != null) result = tmp.index
		// console.log("Tile for : "+x+"x"+y+" converted to: "+newX+"x"+newY+" = "+result)
		return result
	}

	removeTileAt(x, y) {
		var newX = x - (this.items.layer.offsetX/32) // - this.coord.x*this.game.Properties.areaWidth
		var newY = y - (this.items.layer.offsetY/32) // - this.coord.y*this.game.Properties.areaHeight
		this.data.removeTile(newX, newY, this.items).destroy()
	}

	addTileAt(id, x, y) {
		var newX = x - (this.items.layer.offsetX/32) // - this.coord.x*this.game.Properties.areaWidth
		var newY = y - (this.items.layer.offsetY/32) // - this.coord.y*this.game.Properties.areaHeight
		this.data.putTile(id, newX, newY, this.items)
	}

	findSameTileZone(x, y, tileID, tileList) {
		// console.log("Adding "+x+"x"+y)
		tileList.add(this.data.getTile(x, y, this.front))

		if ((this.data.getTile(x + 1, y, this.back).index == tileID) && (!tileList.has(this.data.getTile(x + 1, y, this.front))))
			this.findSameTileZone(x + 1, y, tileID, tileList)
		if ((this.data.getTile(x - 1, y, this.back).index == tileID) && (!tileList.has(this.data.getTile(x - 1, y, this.front))))
			this.findSameTileZone(x - 1, y, tileID, tileList)
		if ((this.data.getTile(x, y + 1, this.back).index == tileID) && (!tileList.has(this.data.getTile(x, y + 1, this.front))))
			this.findSameTileZone(x, y + 1, tileID, tileList)
		if ((this.data.getTile(x, y - 1, this.back).index == tileID) && (!tileList.has(this.data.getTile(x, y - 1, this.front))))
			this.findSameTileZone(x, y - 1, tileID, tileList)
	}

	destroy() {
		if (this.status == 2) {
			this.back.destroy()
			this.items.destroy()
			this.front.destroy()
			this.data.destroy()
		}
	}
}

class World
{
    constructor(game) {
		this.playerArea = new Phaser.Point(-1, -1)
        this.WorldMap = new Area(game)
		this.game = game
		this.buildings = new Map()
    }

    updateArea(x, y) {
		var newarea = new Phaser.Point(Math.floor(x/this.game.Properties.areaWidth), Math.floor(y/this.game.Properties.areaHeight))
		if (!Phaser.Point.equals(newarea, this.playerArea)) {
			console.log("Player reach new area: "+this.playerArea)
			this.playerArea = newarea
			this.checkLoadedMaps(this.playerArea.x, this.playerArea.y)
		}
	}

	checkLoadedMaps(x, y) {
		var name = x+'_'+y
		this.game.DynLoad.loadMap(name, this.swapMap.bind(this))
	}

	swapMap(key) {
		this.WorldMap.destroy()
		delete this.WorldMap
		this.WorldMap = new Area(this.game)
		this.WorldMap.name = key
		this.WorldMap.coord = this.playerArea
		this.WorldMap.status = 1
		this.WorldMap.render()
	}

    getTileInArea(x, y) {
		return this.WorldMap.getTileValueAt(x, y, this.WorldMap.front)
	}

	getItemInArea(x, y) {
		// console.log(x, y, this.WorldMap.getItemValueAt(x, y))
		return this.WorldMap.getTileValueAt(x, y, this.WorldMap.items)
	}

	isFreeSpace(x, y) {
		var b = this.game.TilesList[this.WorldMap.getTileValueAt(x, y, this.WorldMap.back)]
		var i = this.game.TilesList[this.WorldMap.getTileValueAt(x, y, this.WorldMap.items)]
		var f = this.game.TilesList[this.WorldMap.getTileValueAt(x, y, this.WorldMap.front)]

		// console.log(b, i, f)
		if (!b.block & !i.block & !f.block)
			return true
		else return false
	}

	removeTileInArea(x, y) {
		this.WorldMap.removeTileAt(x, y)
	}

	addTileInArea(id, x, y) {
		this.WorldMap.addTileAt(id, x, y)
	}

	enterBuilding(x, y, tileID) {
		var newX = x - (this.WorldMap.front.layer.offsetX/32) // - this.coord.x*this.game.Properties.areaWidth
		var newY = y - (this.WorldMap.front.layer.offsetY/32) // - this.coord.y*this.game.Properties.areaHeight
		var name = newX+"_"+newY
		if (!this.buildings.has(name)) {
			var tileList = new Set()
			this.WorldMap.findSameTileZone(newX, newY, tileID, tileList)
			this.buildings.set(name, tileList)
		} else {
			var tileList = this.buildings.get(name)
		}
		tileList.forEach(function(val1, val2, zeSet){
			val1.alpha = 0
		}, this)
		this.WorldMap.front.dirty = true
	}

	exitBuilding(x, y) {
		var newX = x - (this.WorldMap.front.layer.offsetX/32) // - this.coord.x*this.game.Properties.areaWidth
		var newY = y - (this.WorldMap.front.layer.offsetY/32) // - this.coord.y*this.game.Properties.areaHeight
		var name = newX+"_"+newY
		var tileList = this.buildings.get(name)
		tileList.forEach(function(val1, val2, zeSet){
			val1.alpha = 1
		}, this)
		this.WorldMap.front.dirty = true
	}
}

module.exports = World
