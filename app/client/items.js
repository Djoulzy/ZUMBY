'use strict'

class Pocket {
	constructor(game) {
	}
}

class Bag
{
	constructor(game, nbPockets) {
		this.game = game
		this.bagSize = nbPockets
		this.inventory = this.game.add.group()
		this.inventory.enableBody = true
		this.inventory.fixedToCamera = true
		// this.inventory.inputEnabled = true
		this.cartouche = new Phaser.Point(960, 0)

		this.sac = new Array(nbPockets)
		this.sac.fill(0)

		this.pocketPos = new Array()
		this.pocketPos[0] = new Phaser.Point(this.cartouche.x + 118, this.cartouche.y + 474)
		this.pocketPos[1] = new Phaser.Point(this.cartouche.x + 170, this.cartouche.y + 474)
		this.pocketPos[2] = new Phaser.Point(this.cartouche.x + 66, this.cartouche.y + 526)
		this.pocketPos[3] = new Phaser.Point(this.cartouche.x + 118, this.cartouche.y + 526)
		this.pocketPos[4] = new Phaser.Point(this.cartouche.x + 170, this.cartouche.y + 526)
		this.pocketPos[5] = new Phaser.Point(this.cartouche.x + 222, this.cartouche.y + 526)
		this.pocketPos[6] = new Phaser.Point(this.cartouche.x + 66, this.cartouche.y + 578)
		this.pocketPos[7] = new Phaser.Point(this.cartouche.x + 118, this.cartouche.y + 578)
		this.pocketPos[8] = new Phaser.Point(this.cartouche.x + 170, this.cartouche.y + 578)
		this.pocketPos[9] = new Phaser.Point(this.cartouche.x + 222, this.cartouche.y + 578)
	}

	findEmptyZone() {
		for (var i=0; i<this.bagSize; i++) {
			if (this.sac[i] == 0) return i
		}
		return false
	}

	addItem(id) {
		var zone = this.findEmptyZone()
		console.log(zone)
		if (zone !== false) {
			var newItem = this.inventory.getFirstExists(false);
			if (!newItem)
			{
				newItem = this.inventory.create(this.pocketPos[zone].x, this.pocketPos[zone].y, 'final', id-1)
				newItem.inputEnabled = true
				newItem.input.enableDrag(true)
			}
			newItem.reset(this.pocketPos[zone].x, this.pocketPos[zone].y);
			this.sac[zone] = id
		}
	}
}

module.exports = Bag
