'use strict'

var Connection = function (addr, callback) {

	this.HELLO = "[HELO]"
	this.BROADCAST = "[BCST]"
	this.PLAYERMOVE = "[PMOV]"
	this.PLAYERSHOOT = "[FIRE]"
	this.NEWUSER = "[NUSR]"
	this.PICKITEM = "[PICK]"
	this.DROPITEM = "[DROP]"
	this.UPDATEIVENTORY = "[UPDI]"

    var ws = new WebSocket ('ws://'+addr+'/ws');
    var brothers = new Set();
    var connEvt = new Set();

	this.on = function(evt, newCallback) {
		connEvt[evt] = newCallback
	}

    ws.onopen = callback

    ws.onmessage = function(evt) {
		var cmd = evt.data.split("|");
		var len = cmd.length
		// console.log(cmd)
		for (var i = 0; i < len; i++) {
	    	switch(cmd[i].substr(0, 6))
	    	{
	    		case "[RDCT]":
	    			reconnect(cmd[i].substr(6))
	    			break;
	    		case "[FLBK]":
	    			var obj = JSON.parse(cmd[i].substr(6))
	    			for (var k in obj.BRTHLST){
	    			    if (obj.BRTHLST.hasOwnProperty(k))
	    					 brothers.add(obj.BRTHLST[k].Httpaddr)
	    			}
	    			break;
	            case "[BCST]":
	                var obj = JSON.parse(cmd[i].substr(6))
                    if (connEvt["enemy_move"]) connEvt["enemy_move"](obj)
	                break;
				case "[KILL]":
	                // console.log(cmd[i].substr(6));
					connEvt["kill_enemy"](cmd[i].substr(6))
	                break;
				case "[NENT]":
					var obj = JSON.parse(cmd[i].substr(6))
					connEvt["new_entity"](obj);
					break;
				case "[HIDE]":
					var obj = JSON.parse(cmd[i].substr(6))
					connEvt["hide_item"](obj);
					break;
				case "[SHOW]":
					var obj = JSON.parse(cmd[i].substr(6))
					connEvt["show_item"](obj);
					break;
				case "[WLCM]":
					var obj = JSON.parse(cmd[i].substr(6))
					connEvt["userlogged"](obj)
					break;
				case "[CHAT]":
					var obj = JSON.parse(cmd[i].substr(6))
					connEvt["chat_message"](obj)
					break;
	    		default:;
	    	}
		}
    }

	ws.onclose = function(evt) {
		switch(evt.code)
		{
			case 1005:
				console.log("CLOSE By Client");
				ws = null;
				break;
			case 1000:
				console.log("CLOSE By SERVER: " + evt.reason);
				ws = null;
				break;
			case 1006:
			default:
				// console.log("Lost Connection: " + evt.reason);
				// for (let item of brothers) {
				// 	reconnect(item)
				// 	if (ws.readyState == 0) {
				// 		brothers.delete(item)
				// 	}
				// 	else break;
				// }
				break;
		}
	}

	this.sendTextMessage = function(action, message) {
		ws.send(action + message);
	}

	this.sendJsonMessage = function(action, message) {
		ws.send(action + JSON.stringify(message));
	}
}

module.exports = Connection;
