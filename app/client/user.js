'use strict'

class User
{
	constructor(game, id, face, startx, starty) {
		this.game = game
		this.User_id = id;
		this.face = face
		this.X = startx
		this.Y = starty
		this.dest_X = startx
	    this.dest_Y = starty
		this.step = this.game.Properties.step
		this.inGame = false
		this.PlayerIsMoving = false
		// this.line = new Phaser.Line(0, 0, 100, 100);
	    // this.graphics=game.add.graphics(0,0);
		this.game.DynLoad.loadUser(face, this.initSprite.bind(this))
	}

	initAnims() {
		this.sprite.animations.add('left', [3, 4, 5], 10, true);
	    this.sprite.animations.add('right', [6, 7, 8], 10, true);
	    this.sprite.animations.add('up', [9, 10, 11], 10, true);
	    this.sprite.animations.add('down', [0, 1, 2], 10, true);
	}

	initSprite(face) {
		this.sprite = this.game.add.sprite(this.X*this.step, this.Y*this.step, this.face);
		this.game.midLayer.add(this.sprite)
		this.game.physics.arcade.enable(this.sprite);
	    // this.sprite.body.collideWorldBounds = true;
		this.sprite.body.setSize(this.step, this.step);
		this.sprite.body.onMoveComplete.add(this.moveOver, this);
		this.inGame = true
		this.initAnims()
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

	isMoving() {
		return this.PlayerIsMoving
	}

	destroy() {
		this.sprite.kill()
	}
}

module.exports = User
