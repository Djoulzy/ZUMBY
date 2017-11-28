'use strict'

var Connection = function (addr, callback) {
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
	    			var obj = JSON.parse(cmd[i].substr(6));
	    			for (var k in obj.BRTHLST){
	    			    if (obj.BRTHLST.hasOwnProperty(k))
	    					 brothers.add(obj.BRTHLST[k].Httpaddr)
	    			}
	    			break;
	            case "[BCST]":
	                var obj = JSON.parse(cmd[i].substr(6));
                    if (connEvt["enemy_move"]) connEvt["enemy_move"](obj);
	                break;
				case "[KILL]":
	                // console.log(cmd[i].substr(6));
					connEvt["kill_enemy"](cmd[i].substr(6));
	                break;
				case "[NENT]":
					var obj = JSON.parse(cmd[i].substr(6));
					console.log(obj);
					connEvt["new_entity"](obj);
					break;
				case "[WLCM]":
					// console.log(cmd[i])
					var obj = JSON.parse(cmd[i].substr(6))
					connEvt["userlogged"](obj)
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

	this.logon = function(pass) {
        ws.send("[HELO]" + pass);
		// connEvt["userlogged"].call(this);
	}

    // this.bcast = function(message) {
	// 	// console.log(message);
    //     ws.send("[BCST]" + JSON.stringify(message))
    // }

    this.playerMove = function(message) {
		// console.log(message);
        ws.send("[PMOV]" + JSON.stringify(message))
    }

	this.playerShoot = function(message) {
		// console.log(message);
        ws.send("[FIRE]" + JSON.stringify(message))
    }

	this.newPlayer = function(message) {
		// console.log(message);
        ws.send("[NUSR]" + JSON.stringify(message))
    }
}

module.exports = Connection;
