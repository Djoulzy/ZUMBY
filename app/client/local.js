'use strict';

const User = require('client/user');

class Local extends User
{
	constructor(game, id, face, startx, starty) {
		super(game, id, face, startx, starty)

		this.isPlayer = true
		this.PlayerOrdersCount = 0
		// this.sprite.body.onMoveComplete.add(this.moveLocalOver, this);
		// this.graphics.lineStyle(2, 0xffd900, 1);
		this.bearing = "down"
		this.area = new Phaser.Point(0, 0)
		// this.game.WorldMap.updateArea(startx, starty)

		this.PV = 0
		this.ST = 0
		this.TH = 0
		this.FGT = 0
		this.SHT = 0
		this.CFT = 0
		this.BRD = 0
		this.GRW = 0

		// this.game.DynLoad.loadUser(face, this.initSprite.bind(this))
		this.game.WorldMap.updateArea(startx, starty)
	}

	initSprite(face) {
		super.initSprite(face)
		this.game.camera.follow(this.sprite)
	}

	setAttr(data) {
		this.PV = data.pv
		this.ST = data.st
		this.TH = data.th
		this.FGT = data.fgt
		this.SHT = data.sht
		this.CFT = data.cft
		this.BRD = data.brd
		this.GRW = data.grw
	}

	fire(portee) {
		this.game.socket.playerShoot({
			typ: "P",
			id: this.User_id,
			x: this.X,
			y: this.Y,
			mov: this.bearing,
			pow: portee
		})
	}

	sendMoveToServer(move) {
		if (this.isPlayer) {
			this.bearing = move
			this.PlayerOrdersCount += 1;
			// console.log("Sending: "+player.sprite.dest_x+"  "+player.sprite.dest_y)
			// this.graphics.moveTo(this.sprite.body.x + 16, this.sprite.body.y + 16);//moving position of graphic if you draw mulitple lines
		    // this.graphics.lineTo(this.sprite.dest_x + 16, this.sprite.dest_y + 16);
		    // this.graphics.endFill();
			this.game.socket.playerMove({
				typ: "P",
				id: this.User_id,
				png: this.face,
				num: this.PlayerOrdersCount,
				mov: move,
				spd: 1,
				x: this.dest_X,
				y: this.dest_Y
			})
		}
		this.game.WorldMap.updateArea(this.dest_X, this.dest_Y)
		this.PlayerIsMoving = true
	}

	moveOver() {
		this.adjustSpritePosition()
		this.PlayerIsMoving = false
		this.sprite.animations.stop();
	}

	moveLeft(step, speed) {
		if (this.game.WorldMap.getTileInArea(this.X - 1, this.Y) == 0) {
			this.dest_X = this.X - 1
			this.dest_Y = this.Y
			this.sendMoveToServer('left')
			this.sprite.body.moveTo(speed, step, 180);
			this.sprite.animations.play('left');
		} else {
			this.PlayerIsMoving = false
			this.sprite.frame = 4;
		}
	}

	moveRight(step, speed) {
		if (this.game.WorldMap.getTileInArea(this.X + 1, this.Y) == 0) {
			this.dest_X = this.X + 1
			this.dest_Y = this.Y
			this.sendMoveToServer('right')
			this.sprite.body.moveTo(speed, step, 0);
			this.sprite.animations.play('right');
		} else {
			this.PlayerIsMoving = false
			this.sprite.frame = 7;
		}
	}

	moveUp(step, speed) {
		if (this.game.WorldMap.getTileInArea(this.X, this.Y - 1) == 0) {
			this.dest_X = this.X
			this.dest_Y = this.Y - 1
			this.sendMoveToServer('up')
			this.sprite.body.moveTo(speed, step, 270);
			this.sprite.animations.play('up');
		} else {
			this.PlayerIsMoving = false
			this.sprite.frame = 10;
		}
	}

	moveDown(step, speed) {
		if (this.game.WorldMap.getTileInArea(this.X, this.Y + 1) == 0) {
			this.dest_X = this.X
			this.dest_Y = this.Y + 1
			this.sendMoveToServer('down')
			this.sprite.body.moveTo(speed, step, 90);
			this.sprite.animations.play('down');
		} else {
			this.PlayerIsMoving = false
			this.sprite.frame = 1;
		}
	}
}

module.exports = Local
