'use strict'

class Chat
{
    constructor(game) {
		this.ZChat = new Draggabilly( '.draggable', {
			handle: '.handle'
        })
        this.txtWindow = document.getElementById("ZFlow")
        this.txtWindow.scrollTop = this.txtWindow.scrollHeight
    }

    addMessage(obj) {
        this.txtWindow.innerHTML += "["+obj.from+"] "+obj.mess+"<br/>"
        this.txtWindow.scrollTop = this.txtWindow.scrollHeight
    }
}

module.exports = Chat