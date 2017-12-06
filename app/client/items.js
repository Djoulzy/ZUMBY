'use strict'

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


		// this.addItem(151)
		// this.addItem(151)
		// this.addItem(151)
	}

	findEmptyZone() {
		for (var i=0; i<this.bagSize; i++) {
			if (this.sac[i] === 0) return i
		}
		return false
	}

	addItem(id, pocket) {
		if (this.sac[pocket] == 0) {
			var newItem = this.inventory.getFirstExists(false);
			if (!newItem)
			{
				newItem = this.inventory.create(this.pocketPos[pocket].x, this.pocketPos[pocket].y, 'final', id-1)
				newItem.inputEnabled = true
				newItem.input.enableDrag(true)
				newItem.events.onDragStart.add(this.startDrag, this)
				newItem.events.onDragStop.add(this.stopDrag, this)
				newItem.originalPosition = newItem.position.clone()
				newItem.tileID = id
				newItem.pocket = pocket
			}
			newItem.reset(this.pocketPos[pocket].x, this.pocketPos[pocket].y)
			this.sac[pocket] = id
		}
	}

	loadInventory(itemList) {
		for (var i=0; i<10; i++) {
			if (itemList[i].id != 0) this.addItem(itemList[i].id, i)
		}
	}

	combineItems(item1, item2) {

	}

	sendUpdate(id, from, to) {
		this.game.socket.sendJsonMessage(this.game.socket.UPDATEIVENTORY, {
				owner: this.game.player.User_id,
				id: id,
				fp: from,
				tp: to
			})
	}

	throwItem(item, pointer) {
		var fullx = Math.floor((this.game.camera.x + item.x)/32)
		var fully = Math.floor((this.game.camera.y + item.y)/32)
		if (this.game.WorldMap.isFreeSpace(fullx, fully)) {
			this.game.socket.sendJsonMessage(this.game.socket.DROPITEM, {
				owner: this.game.player.User_id,
				id: item.tileID,
				fp: item.pocket,
				x: fullx,
				y: fully
			})
			this.sac[item.pocket] = 0
			item.kill()
			return true
		}
		return false
	}

	startDrag(draggedSprite, pointer) {

	}

	checkDropZone(pointer) {
		for (var i=0; i<this.bagSize; i++) {
			var xx = pointer.x - this.pocketPos[i].x
			var yy = pointer.y - this.pocketPos[i].y
			if (( xx <= 32) && (xx >= 0) && ( yy < 32) && (yy >= 0))
				return i
		}
		return false
	}

	stopDrag(draggedSprite, pointer) {
		var zone = this.checkDropZone(pointer)
		if (zone !== false) {
			if (this.sac[zone] === 0) {
				draggedSprite.position.x = this.pocketPos[zone].x
				draggedSprite.position.y = this.pocketPos[zone].y
				draggedSprite.originalPosition = draggedSprite.position.clone()
				this.sendUpdate(draggedSprite.tileID, draggedSprite.pocket, zone)
				this.sac[zone] = draggedSprite.tileID
				this.sac[draggedSprite.pocket] = 0
				draggedSprite.pocket = zone
				return
			} else {
				if (this.combineItems(draggedSprite.tileID, this.sac[zone])) return
			}
		} 
		if (pointer.x < this.cartouche.x) {
			if (this.throwItem(draggedSprite, pointer)) return
		}
		draggedSprite.position.copyFrom(draggedSprite.originalPosition)
	}
}

module.exports = Bag