'use strict'

class OSD
{
	constructor(game) {
		this.game = game
		this.border = 20
		this.cartouche = this.game.add.sprite(960, 0, 'cartouche')
		this.cartouche.fixedToCamera = true

		var Gros = { font: "bold 32px Arial", fill: "#fff", boundsAlignH: "center", boundsAlignV: "middle" };
		var Normal = { font: "bold 16px Arial", fill: "#fff", boundsAlignH: "center", boundsAlignV: "middle" };

		this.label_name = this.game.add.text(this.border, 20, "", Gros)
		this.label_name.setShadow(3, 3, 'rgba(150, 67, 49, 4)', 2);
		this.cartouche.addChild(this.label_name);

		this.Attr_PV = this.game.add.text(this.border, 200, "PV: " + 0, Normal)
		this.cartouche.addChild(this.Attr_PV)
		this.Attr_ST = this.game.add.text(this.border, 220, "Starve: " + 0, Normal)
		this.cartouche.addChild(this.Attr_ST)
		this.Attr_TH = this.game.add.text(this.border, 240, "Thirst: " + 0, Normal)
		this.cartouche.addChild(this.Attr_TH)

		this.Attr_FT = this.game.add.text(this.border, 300, "Fight: " + 0, Normal)
		this.cartouche.addChild(this.Attr_FT)
		this.Attr_SH = this.game.add.text(this.border, 320, "Shoot: " + 0, Normal)
		this.cartouche.addChild(this.Attr_SH)
		this.Attr_CR = this.game.add.text(this.border, 340, "Craft: " + 0, Normal)
		this.cartouche.addChild(this.Attr_CR)
		this.Attr_BR = this.game.add.text(this.border, 360, "Breed: " + 0, Normal)
		this.cartouche.addChild(this.Attr_BR)
		this.Attr_GR = this.game.add.text(this.border, 380, "Grow: " + 0, Normal)
		this.cartouche.addChild(this.Attr_GR)

		this.label_area = this.game.add.text(this.border, 500, "Area: 0 x 0", Normal)
		this.cartouche.addChild(this.label_area);

		// this.game.frontLayer.add(this.label_score)
		this.game.frontLayer.add(this.cartouche)
	}

	refresh() {
		this.label_name.setText(this.game.player.User_id)
		this.label_area.setText("Area: " + this.game.WorldMap.playerArea.x + " x " + this.game.WorldMap.playerArea.y)
	}
}

module.exports = OSD
