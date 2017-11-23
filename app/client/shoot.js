'use strict'

class Shoot
{
	constructor(game) {
		this.game = game
		this.bullets = this.game.add.group();
	    this.bullets.enableBody = true;
	    this.bullets.physicsBodyType = Phaser.Physics.ARCADE;
		this.bulletTime = 0;
	}

	moveOver(sprt) {
		sprt.kill();
	}

	sendMoveToServer(move) {
	}

	showFire(from) {
		var fire = this.bullets.getFirstExists(false);
		if (!fire)
		{
			fire = this.bullets.create(from.body.x, from.body.y, 'shoot', 1000);
			fire.animations.add('fire', Phaser.Animation.generateFrameNames('fire', 1, 3));
		}
		fire.bringToTop()
        fire.reset(from.body.x, from.body.y);
		fire.play('fire', 30, false, true)
	}

	moveLeft(bullet, step, speed) {
		this.showFire(bullet)
		bullet.frameName = "bullet3"
		this.sendMoveToServer('left')
		bullet.body.moveTo(speed, step, 180);
	}

	moveRight(bullet, step, speed) {
		this.showFire(bullet)
		bullet.frameName = "bullet1"
		this.sendMoveToServer('right')
		bullet.body.moveTo(speed, step, 0);
	}

	moveUp(bullet, step, speed) {
		this.showFire(bullet)
		bullet.frameName = "bullet4"
		this.sendMoveToServer('up')
		bullet.body.moveTo(speed, step, 270);
	}

	moveDown(bullet, step, speed) {
		this.showFire(bullet)
		bullet.frameName = "bullet2"
		this.sendMoveToServer('down')
		bullet.body.moveTo(speed, step, 90);
	}

	fire(from, portee, speed) {

	    //  To avoid them being allowed to fire too fast we set a time limit
	    if (this.game.time.now > this.bulletTime)
	    {
			from.fire(portee)
	        //  Grab the first bullet we can from the pool
	        var bullet = this.bullets.getFirstExists(false);
	        if (!bullet)
	        {
				bullet = this.bullets.create(from.sprite.body.x, from.sprite.body.y, 'shoot');
				bullet.checkWorldBounds = true;
				bullet.outOfBoundsKill = true;
			}
            //  And fire it
            bullet.reset(from.sprite.body.x, from.sprite.body.y);
			bullet.body.onMoveComplete.add(this.moveOver, this);

			switch(from.bearing) {
				case "up": this.moveUp(bullet, portee*this.game.Properties.step, speed); break;
				case "down": this.moveDown(bullet, portee*this.game.Properties.step, speed); break;
				case "left": this.moveLeft(bullet, portee*this.game.Properties.step, speed); break;
				case "right": this.moveRight(bullet, portee*this.game.Properties.step, speed); break;
			}
            this.bulletTime = this.game.time.now + 500;
	    }

	}
}

module.exports = Shoot
