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
		this.tarrain = null
		this.obstacles = null
		this.hauteurs = null
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
			this.terrain = this.data.createLayer('terrain')
			this.obstacles = this.data.createLayer('obstacles')
			this.hauteurs = this.data.createLayer('hauteurs')

			this.game.backLayer.add(this.terrain)
			this.game.backLayer.add(this.obstacles)
			this.game.frontLayer.add(this.hauteurs)

			// this.game.frontLayer.sort();
			this.status = 2
		}
	}

	getTileValueAt(x, y) {
		var result = 0
		var newX = x - (this.obstacles.layer.offsetX/32) // - this.coord.x*this.game.Properties.areaWidth
		var newY = y - (this.obstacles.layer.offsetY/32) // - this.coord.y*this.game.Properties.areaHeight
		var tmp = this.data.getTile(newX, newY, this.obstacles)
		if (tmp != null) result = tmp.index
		// console.log("Tile for : "+x+"x"+y+" converted to: "+newX+"x"+newY+" = "+result)
		return result
	}

	getItemValueAt(x, y) {
		var result = 0
		var newX = x - (this.hauteurs.layer.offsetX/32) // - this.coord.x*this.game.Properties.areaWidth
		var newY = y - (this.hauteurs.layer.offsetY/32) // - this.coord.y*this.game.Properties.areaHeight
		var tmp = this.data.getTile(newX, newY, this.hauteurs)
		if (tmp != null) result = tmp.index
		// console.log("Tile for : "+x+"x"+y+" converted to: "+newX+"x"+newY+" = "+result)
		return result
	}

	removeTileAt(x, y) {
		var newX = x - (this.hauteurs.layer.offsetX/32) // - this.coord.x*this.game.Properties.areaWidth
		var newY = y - (this.hauteurs.layer.offsetY/32) // - this.coord.y*this.game.Properties.areaHeight
		this.data.removeTile(newX, newY, this.hauteurs)
	}

	destroy() {
		if (this.status == 2) {
			this.terrain.destroy()
			this.obstacles.destroy()
			this.hauteurs.destroy()
			this.data.destroy()
		}
	}
}

class Map
{
    constructor(game) {
		this.playerArea = new Phaser.Point(-1, -1)
        this.WorldMap = new Area(game)
        this.game = game
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
		return this.WorldMap.getTileValueAt(x, y)
	}

	getItemInArea(x, y) {
		return this.WorldMap.getItemValueAt(x, y)
	}

	removeTileInArea(x, y) {
		this.WorldMap.removeTileAt(x, y)
	}
}

module.exports = Map
