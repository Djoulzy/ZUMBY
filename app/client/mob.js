'use strict';

const Remote = require('client/remote');

class Mob extends Remote
{
	constructor(game, id, face, subview, startx, starty) {
		super(game, id, face, subview, startx, starty)
		// this.graphics.lineStyle(2, 0xf11010 , 1);
	}

	initAnims() {
		var visual = Number(this.subview)*12
	    this.sprite.animations.add('down', [visual+0, visual+1, visual+2], 10, true);
		this.sprite.animations.add('left', [visual+3, visual+4, visual+5], 10, true);
	    this.sprite.animations.add('right', [visual+6, visual+7, visual+8], 10, true);
	    this.sprite.animations.add('up', [visual+9, visual+10, visual+11], 10, true);
	}
}

module.exports = Mob
