'use strict'

const User = require('client/user');

class Remote extends User
{
	constructor(game, id, face, subview, startx, starty) {
		super(game, id, face, startx, starty)

		this.subview = subview
		this.isPlayer = false
		this.PlayerOrdersCount = 0
		// this.sprite.body.onMoveComplete.add(this.moveRemoteOver, this);
		// this.graphics.lineStyle(2, 0x14c818 , 1);
		this.moves = new Array()

		if (subview == "")
			this.game.DynLoad.loadUser(face, this.initSprite.bind(this))
	}

	// debugLine() {
	// 	this.graphics.moveTo(this.sprite.body.x + 16, this.sprite.body.y + 16);//moving position of graphic if you draw mulitple lines
	// 	this.graphics.lineTo(this.sprite.dest_x + 16, this.sprite.dest_y + 16);
	// 	this.graphics.endFill();
	// }

	moveOver() {
		this.adjustSpritePosition()
		this.PlayerIsMoving = false
		if (this.moves.length == 0) {
			this.sprite.animations.stop();
		}
		// this.sprite.frame = 1;
	}

	moveLeft(step, speed) {
		// this.debugLine()
		this.PlayerIsMoving = true
		this.sprite.body.moveTo(speed, step, 180);
		this.sprite.animations.play('left');
	}

	moveRight(step, speed) {
		// this.debugLine()
		this.PlayerIsMoving = true
		this.sprite.body.moveTo(speed, step, 0);
		this.sprite.animations.play('right');
	}

	moveUp(step, speed) {
		// this.debugLine()
		this.PlayerIsMoving = true
		this.sprite.body.moveTo(speed, step, 270);
		this.sprite.animations.play('up');
	}

	moveDown(step, speed) {
		// this.debugLine()
		this.PlayerIsMoving = true
		this.sprite.body.moveTo(speed, step, 90);
		this.sprite.animations.play('down');
	}
}

module.exports = Remote
