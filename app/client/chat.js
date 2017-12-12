'use strict'

class Chat
{
    constructor(game, cursors) {
        this.game = game
        this.cursors = cursors

		var ZeChat = new Draggabilly( '.draggable', {
			handle: '.handle'
        })

        this.ZChat = document.getElementById("ZChat")
        this.ZChat.onmouseover = this.onMouseOver.bind(this)
        this.ZChat.onmouseout = this.onMouseOut.bind(this)

        var chatForm = document.getElementById("ZForm")
        this.txtInput = document.getElementById("ZInput")
        this.txtWindow = document.getElementById("ZFlow")

        chatForm.onsubmit= this.onSendMessage.bind(this)
        this.txtWindow.scrollTop = this.txtWindow.scrollHeight
    }

    addMessage(obj) {
        this.txtWindow.innerHTML += "["+obj.from+"] "+obj.mess+"<br/>"
        this.txtWindow.scrollTop = this.txtWindow.scrollHeight
    }

    onMouseOver() {
        this.game.input.enabled = false
        // this.cursors.pickup.enabled = false
        // this.cursors.pickup.reset();
    }

    onMouseOut() {
        this.game.input.enabled = true
        // this.cursors.pickup.enabled = false
        // this.cursors.pickup.reset();
    }

    onSendMessage() {
        var textMessage = this.escapeHtml(this.txtInput.value)
        console.log(textMessage)
        var mess = {
            from: this.game.Properties.pseudo,
            type: 4,
            mess: textMessage
        }
        this.addMessage(mess);
        this.txtInput.value = ''
        return false
    }

    escapeHtml(unsafe) {
        return unsafe
            .replace(/&/g, "&amp;")
            .replace(/</g, "&lt;")
            .replace(/>/g, "&gt;")
            .replace(/"/g, "&quot;")
            .replace(/'/g, "&#039;");
    }
}

module.exports = Chat