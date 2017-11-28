'use strict';

const Config = require('config');

function Boot(){}

Boot.prototype = {
    preload: function(){
        // Debbug
        // this.game.plugins.add(Phaser.Plugin.Inspector)

        this.game.stage.disableVisibilityChange = true;
        this.game.stage.backgroundColor = 0x3b0760;
        this.game.load.onFileComplete.add(this.onFileComplete, this);
        this.game.load.onLoadComplete.addOnce(this.onLoadComplete, this);

        this.showLoadingText()
		this.loadAssets()
    },

    onFileComplete: function(progress, cacheKey, success, totalLoaded, totalFiles) {
        console.log("File Complete: " + progress + "% - " + totalLoaded + " out of " + totalFiles + "(" + cacheKey + ")")
    },

    onLoadComplete: function() {
        this.game.state.start('play')
    },

    loadAssets: function() {
		// Graphics
	  	this.game.load.spritesheet('h1', 'http://'+Config.MMOServer.Host+'/data/h1.png', 32, 32);
	  	this.game.load.atlas('zombies', 'assets/ZombieSheet.png', 'assets/ZombieSheet.json', Phaser.Loader.TEXTURE_ATLAS_JSON_HASH);
	  	this.game.load.atlas('shoot', 'assets/shoot.png', 'assets/shoot.json', Phaser.Loader.TEXTURE_ATLAS_JSON_HASH);
		this.game.load.spritesheet('final', 'http://'+Config.MMOServer.Host+'/data/final.png', 32, 32);
		this.game.load.image('cartouche', 'assets/cartouche.png');
        this.game.load.start();
    },

    showLoadingText: function(){
        var loadingText = "- Loading -";
        var text = this.game.add.text(this.game.world.centerX, this.game.world.centerY, loadingText);
        //  Centers the text
        text.anchor.set(0.5);
        text.align = 'center';

        //  Our font + size
        text.font = 'Arial';
        text.fontWeight = 'bold';
        text.fontSize = 70;
        text.fill = '#ffffff';
    }
}

module.exports = Boot
