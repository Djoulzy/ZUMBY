'use strict'

var Config = require('config')
var Connection = require('client/wsconnect')
var DynLoad = require('client/dynload')
var WMap = require('client/map')
var OSD = require('client/osd')
var Local = require('client/local')
var Remote = require('client/remote')
var Mob = require('client/mob')
var Shoot = require('client/shoot')
var Explode = require('client/explode')
var Bag = require('client/items')
var Chat = require('client/chat')

function Play(){}

Play.prototype = {

    create: function() {
		this.game.physics.startSystem(Phaser.Physics.ARCADE);
		this.running = false
		this.cursors = this.game.input.keyboard.addKeys({
			'pickup': Phaser.Keyboard.P,
			'space': Phaser.Keyboard.SPACEBAR,
			'up': Phaser.Keyboard.UP,
			'down': Phaser.Keyboard.DOWN,
			'left': Phaser.Keyboard.LEFT,
			'right': Phaser.Keyboard.RIGHT })
        this.game.DynLoad = new DynLoad(this.game)

		this.game.Chat = new Chat(this.game)

		this.entities = [];
        this.game.backLayer = this.game.add.group()
		this.game.midLayer = this.game.add.group()
		this.game.frontLayer = this.game.add.group()

		this.initSocket()
		this.bullets = new Shoot(this.game)
		this.explode = new Explode(this.game)
		this.inventory = new Bag(this.game, 10)

		this.game.WorldMap = new WMap(this.game)
		this.game.OSD = new OSD(this.game)

		this.game.camera.view = new Phaser.Rectangle(0,0,960,768)
		// this.game.camera.deadzone = new Phaser.Rectangle(100, 100, 600, 400);

		this.game.TilesList = this.game.cache.getJSON('tilesList')
    },

////////////////////////////////////////////////////
//                      NETWORK                   //
////////////////////////////////////////////////////
	initSocket: function() {
		this.game.socket = new Connection(Config.MMOServer.Host, this.onSocketConnected.bind(this))
       	this.game.socket.on("userlogged", this.onUserLogged.bind(this))
		this.game.socket.on("new_entity", this.newEntitie.bind(this))
      	this.game.socket.on("enemy_move", this.onEnemyMove.bind(this));
      	this.game.socket.on("kill_enemy", this.onRemoveEntity.bind(this));
      	this.game.socket.on("hide_item", this.onRemoveItem.bind(this));
		this.game.socket.on("show_item", this.onAddItem.bind(this));
      	this.game.socket.on("chat_message", this.game.Chat.addMessage.bind(this.game.Chat));
    },

	findGetParameter: function(parameterName) {
		var result = null
		var tmp = []
		location.search.substr(1).split("&").forEach(function (item) {
			tmp = item.split("=");
			if (tmp[0] === parameterName) result = decodeURIComponent(tmp[1]);
		})
		return result
	},

    onSocketConnected: function() {
		var passphrase = this.findGetParameter("key")
		this.game.socket.sendTextMessage(this.game.socket.HELLO, passphrase)
	},

    onUserLogged: function(data) {
		this.game.Properties.pseudo = data.id
		this.game.player = new Local(this.game, data.id, data.png, data.x, data.y)
		this.game.player.setAttr(data)
		this.inventory.loadInventory(data.i)
		this.running = true
    },

	onRemoveItem: function(data) {
		this.game.WorldMap.removeTileInArea(data.x, data.y)
		if (data.owner == this.game.Properties.pseudo) {
			this.inventory.addItem(data.id, data.tp)
		}
	},

	onAddItem: function(data) {
		this.game.WorldMap.addTileInArea(data.id, data.x, data.y)
	},

////////////////////////////////////////////////////
//                      PLAYERS                   //
////////////////////////////////////////////////////
	findplayerbyid: function(id) {
		for (var i = 0; i < this.entities.length; i++) {
			if (this.entities[i].User_id == id) {
				return this.entities[i];
			}
		}
		return false
	},

	newEntitie: function(data) {
		if (data.id == this.game.Properties.pseudo) return
		var movePlayer = this.findplayerbyid(data.id);
		console.log("New Entity: "+data.id)
		if (this.findplayerbyid(data.id)) return
		else {
			if (data.typ == "P") {
                console.log("New Remote Player")
				var new_enemy = new Remote(this.game, data.id, data.png, "", data.x, data.y);
            } else {
				var new_enemy = new Mob(this.game, data.id, "zombies", data.png, data.x, data.y);
            }
			this.entities.push(new_enemy);
		}
	},

////////////////////////////////////////////////////
//                       MOVES                    //
////////////////////////////////////////////////////
	onEnemyMove: function(data) {
		if (data.id == this.game.Properties.pseudo) {
			return
		}

		var movePlayer = this.findplayerbyid(data.id);
		if (!movePlayer) {
			this.newEntitie(data)
			return
		}
		movePlayer.moves.push(data)
	},

	onRemoveEntity: function(id) {
		var removePlayer = this.findplayerbyid(id);
		if (!removePlayer) {
			console.log('Player not found: ', id)
			return
		}

		this.explode.boom(removePlayer.sprite)
		removePlayer.destroy();
		this.entities.splice(this.entities.indexOf(removePlayer), 1);
	},

	updatePlayer: function() {
		// game.physics.arcade.collide(player.sprite, obstacles, playerBlocked);
        if (this.game.player.inGame) {
    		if (!this.game.player.isMoving()) {
    			if (this.cursors.left.isDown) this.game.player.moveLeft(this.game.Properties.step, this.game.Properties.speed)
    			else if (this.cursors.right.isDown) this.game.player.moveRight(this.game.Properties.step, this.game.Properties.speed)
    			else if (this.cursors.up.isDown) this.game.player.moveUp(this.game.Properties.step, this.game.Properties.speed)
    			else if (this.cursors.down.isDown) this.game.player.moveDown(this.game.Properties.step, this.game.Properties.speed)
				else if (this.cursors.pickup.isDown) this.game.player.getItem(this.inventory)
    			else if (this.cursors.space.isDown) {
    				var portee = 5
    				this.bullets.fire(this.game.player, portee, this.game.Properties.speed);
    				// this.loadNewMap()
    			}
			}
        }
	},

	updateRemotePlayers: function() {
		for (var i = 0; i < this.entities.length; i++) {
            if (this.entities[i].inGame) {
    			if (this.entities[i].moves.length > 0 && !this.entities[i].isMoving()) {
    				var move = this.entities[i].moves.shift()
    				this.entities[i].dest_X = move.x;
    				this.entities[i].dest_Y = move.y;
    				this.entities[i].PlayerIsMoving = true
    				var mobSpeed = Math.ceil((this.game.Properties.ServerSpeed * move.spd) / this.game.Properties.step) * this.game.Properties.step + 50;

    				if (move.mov == "left") this.entities[i].moveLeft(this.game.Properties.step, mobSpeed)
    				else if (move.mov == "right") this.entities[i].moveRight(this.game.Properties.step, mobSpeed)
    				else if (move.mov == "up") this.entities[i].moveUp(this.game.Properties.step, mobSpeed)
    				else if (move.mov == "down") this.entities[i].moveDown(this.game.Properties.step, mobSpeed)
    			}
            }
		}
	},

////////////////////////////////////////////////////
//                       LOOPS                    //
////////////////////////////////////////////////////
    update: function() {
		if (this.running) {
			this.updatePlayer()
			this.updateRemotePlayers()
			this.game.OSD.refresh()
		}
    },

	render: function() {
		// Night
	    // this.game.context.fillStyle = 'rgba(0,0,0,0.8)';	    
	    // this.game.context.fillRect(0, 0, 960, 768);

		this.game.DynLoad.start()
		// if (this.running) {
		// 	if (this.game.player.inGame) {
		// 		this.game.debug.spriteInfo(this.game.player.sprite, 32, 32)
		// 	}
		// }

		// Camera
	    // this.game.context.fillStyle = 'rgba(30,0,50,0.8)';
	    // this.game.context.fillRect(5, 620, 300, 140);
		// this.game.debug.cameraInfo(this.game.camera, 10, 640);

		// var zone = this.game.camera.deadzone;
	    // this.game.context.fillStyle = 'rgba(255,0,0,0.6)';
	    // this.game.context.fillRect(zone.x, zone.y, zone.width, zone.height);

		// this.game.debug.gameInfo(32, 500)

		// FPS
	    this.game.context.fillStyle = 'rgba(30,0,50,0.8)';
	    this.game.context.fillRect(970, 630, 300, 130);
		this.game.debug.gameTimeInfo(980, 650)
	}
}

module.exports = Play
