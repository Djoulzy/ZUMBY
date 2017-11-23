'use strict'

var Config = require('config')

class DynLoad
{
    constructor(game) {
        this.game = game
        this.queue = new Map()
        this.game.load.onLoadComplete.add(this.emptyQueue, this);
    }

    loadMap(key, callback) {
        console.log("Added new map to load queue: "+key)
        this.game.load.tilemap(key, 'http://'+Config.MMOServer.Host+'/map/'+key, null, Phaser.Tilemap.TILED_JSON)
        this.queue.set(key, callback)
    }

    loadUser(key, callback) {
        console.log("Added new png to load queue: "+key)
        this.game.load.spritesheet(key, 'http://'+Config.MMOServer.Host+'/data/'+key+'.png', 32, 32);
        this.queue.set(key, callback)
    }

    emptyQueue() {
        this.queue.forEach(function(element, index, theSet) {
			console.log("Calling callback for "+index)
            element(index)
            this.queue.delete(index)
        }, this)
    }

	start() {
		if (this.queue.size) this.game.load.start()
	}
}

module.exports = DynLoad
