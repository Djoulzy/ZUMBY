'use strict'

class Items
{
	constructor(game) {
		this.game = game
		this.inventory = this.game.add.group()
		this.inventory.enableBody = true
		this.inventory.fixedToCamera = true
		this.cartouche = new Phaser.Point(960, 0)

		this.sac = new Array(10)
		this.sac[0] = new Phaser.Point(this.cartouche.x + 118, this.cartouche.y + 474)
		this.sac[1] = new Phaser.Point(this.cartouche.x + 170, this.cartouche.y + 474)
		this.sac[2] = new Phaser.Point(this.cartouche.x + 66, this.cartouche.y + 526)
		this.sac[3] = new Phaser.Point(this.cartouche.x + 118, this.cartouche.y + 526)
	}

	add(id) {
		var zone = this.inventory.length
		var newItem = this.inventory.getFirstExists(false);
		if (!newItem)
		{
			newItem = this.inventory.create(this.sac[zone].x, this.sac[zone].y, 'final', id-1);
		}
		newItem.reset(this.sac[zone].x, this.sac[zone].y);
	}
}

module.exports = Items
