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

function Play(){}

Play.prototype = {

    create: function() {
		this.game.physics.startSystem(Phaser.Physics.ARCADE);
		this.running = false
		this.cursors = this.game.input.keyboard.addKeys({ 'space': Phaser.Keyboard.SPACEBAR, 'up': Phaser.Keyboard.UP, 'down': Phaser.Keyboard.DOWN, 'left': Phaser.Keyboard.LEFT, 'right': Phaser.Keyboard.RIGHT })
        this.game.DynLoad = new DynLoad(this.game)

		this.entities = [];

        this.game.backLayer = this.game.add.group()
		this.game.midLayer = this.game.add.group()
		this.game.frontLayer = this.game.add.group()

		this.initSocket()
		this.bullets = new Shoot(this.game)
		this.explode = new Explode(this.game)

		this.game.WorldMap = new WMap(this.game)
		this.game.OSD = new OSD(this.game)

		this.game.camera.view = new Phaser.Rectangle(0,0,960,768)
		// this.game.camera.deadzone = new Phaser.Rectangle(100, 100, 600, 400);
    },

////////////////////////////////////////////////////
//                      NETWORK                   //
////////////////////////////////////////////////////
	initSocket: function() {
		this.game.socket = new Connection(Config.MMOServer.Host, this.onSocketConnected.bind(this))
       	this.game.socket.on("userlogged", this.onUserLogged.bind(this))
      	this.game.socket.on("enemy_move", this.onEnemyMove.bind(this));
      	this.game.socket.on("kill_enemy", this.onRemoveEntity.bind(this));
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
		this.game.socket.logon(passphrase);
	},

    onUserLogged: function(data) {
		this.game.Properties.pseudo = data.id
		this.game.player = new Local(this.game, data.id, data.png, data.x, data.y)
		this.game.player.setAttr(data)
		this.running = true
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
		var movePlayer = this.findplayerbyid(data.id);
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
			this.game.DynLoad.start()
		// this.game.debug.spriteInfo(this.game.player.sprite, 32, 32);
		// this.game.debug.cameraInfo(this.game.camera, 32, 500);

		// var zone = this.game.camera.deadzone;
	    // this.game.context.fillStyle = 'rgba(255,0,0,0.6)';
	    // this.game.context.fillRect(zone.x, zone.y, zone.width, zone.height);
	}
}

module.exports = Play
