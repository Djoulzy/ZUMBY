'use strict'

class Items
{
	constructor(game, id, face, startx, starty) {
		this.game = game
		this.id = id
		this.face = face
		this.X = startx
		this.Y = starty
		this.step = this.game.Properties.step
		// this.line = new Phaser.Line(0, 0, 100, 100);
	    // this.graphics=game.add.graphics(0,0);
		this.game.DynLoad.loadUser(face, this.initSprite.bind(this))
	}

	initSprite(face) {
		this.sprite = this.game.add.sprite(this.X*this.step, this.Y*this.step, this.face);
		this.game.midLayer.add(this.sprite)
		this.game.physics.arcade.enable(this.sprite);
	    // this.sprite.body.collideWorldBounds = true;
		this.sprite.body.setSize(this.step, this.step);
	}

	adjustSpritePosition() {
		var markerx = this.game.math.snapToFloor(Math.ceil(this.dest_X*this.step), this.step)
		var markery = this.game.math.snapToFloor(Math.ceil(this.dest_Y*this.step), this.step)
		// console.log("Adjusting : x="+this.sprite.body.x+" y="+this.sprite.body.y+" -> x="+ markerx +" y="+markery)
		this.sprite.body.x = markerx
		this.sprite.body.y = markery
		this.X = this.dest_X
		this.Y = this.dest_Y
		// this.graphics.clear();
	}

	destroy() {
		this.sprite.kill()
	}
}

module.exports = Items
