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
        this.ZoneWindow = document.getElementById("ZZone")
        this.GlobalWindow = document.getElementById("ZGlobal")
        this.ServerWindow = document.getElementById("ZServer")

        this.selectedZone = "ZGlobal"

        this.ZoneWindow.style.display = "none"
        this.GlobalWindow.style.display = "inline"
        this.ServerWindow.style.display = "none"

        this.tabLabels = new Map([["T_ZZone", " ZONE "], ["T_ZGlobal", " GLOBAL "], ["T_ZServer", " SERVER "]])

        // Events
        this.sticky = false
        this.lock = document.getElementById("sticky")
        this.lock.onclick = this.onStickyClick.bind(this)
        this.reduce = false
        this.resize = document.getElementById("resize")
        this.resize.onclick = this.onResizeClick.bind(this)

        var tabs = document.getElementsByClassName("tab");
        for(var i = 0; i < tabs.length; i++)
            tabs.item(i).onclick = this.onSelectTab.bind(this)

        var modes = document.getElementsByClassName("mode");
        for(var i = 0; i < modes.length; i++)
            modes.item(i).onclick = this.onSelectMode.bind(this)
        this.mode = 3

        chatForm.onsubmit = this.onSendMessage.bind(this)
    }

    onStickyClick(obj) {
        this.sticky = !this.sticky
        if (this.sticky) this.lock.innerHTML = "lock"
        else this.lock.innerHTML = "lock_open"
    }

    onResizeClick(obj) {
        var container = document.getElementsByClassName("container")
        if (container.item(0).style.display == "none") {
            container.item(0).style.display = "inline"
            this.resize.innerHTML = "vertical_align_top"
        } else {
            container.item(0).style.display = "none"
            this.resize.innerHTML = "vertical_align_bottom"
        }
    }

    addMessage(obj) {
        var tab
        switch(obj.type) {
            case 1:
            case 2: tab = this.ZoneWindow; break;
            case 3: tab = this.GlobalWindow; break;
            default: tab = this.ServerWindow; break;
        }
        tab.innerHTML += "["+obj.from+"] "+obj.mess+"<br/>"

        if (tab.getAttribute("id") == this.selectedZone)
            this.txtWindow.scrollTop = this.txtWindow.scrollHeight
        else {
            var tabName = "T_"+tab.getAttribute("id")
            document.getElementById(tabName).innerHTML = this.tabLabels.get(tabName)+"<span style='color:red;'>*</span>"
        }
    }

    onSelectMode(obj) {
        var selectedMode = obj.target.getAttribute("id")
        switch(selectedMode) {
            case "talk": this.mode = 1; break;
            case "shout": this.mode = 2; break;
            default: this.mode = 3; break;
        }

        var modes = document.getElementsByClassName("mode");
        for(var i = 0; i < modes.length; i++) {
            if (modes.item(i).getAttribute("id") == selectedMode)
                modes.item(i).classList.add('selected')
            else
                modes.item(i).classList.remove('selected')
        }
    }

    onSelectTab(obj) {
        this.selectedZone = obj.target.getAttribute("tab")

        var tabs = document.getElementsByClassName("tab");
        for(var i = 0; i < tabs.length; i++) {
            if (tabs.item(i).getAttribute("tab") == this.selectedZone) {
                tabs.item(i).classList.add('selected')
                var tabName = "T_"+this.selectedZone
                tabs.item(i).innerHTML = this.tabLabels.get(tabName)
            }
            else
                tabs.item(i).classList.remove('selected')
        }

        var zones = document.getElementsByClassName("txtArea");
        for(var i = 0; i < zones.length; i++)
        {
            if (zones.item(i).getAttribute("id") == this.selectedZone) {
                zones.item(i).style.display = "inline"
            }
            else zones.item(i).style.display = "none"
        }
        this.txtWindow.scrollTop = this.txtWindow.scrollHeight
    }

    onMouseOver() {
        this.ZChat.style.opacity = "0.9"
        this.game.input.enabled = false
    }

    onMouseOut() {
        if (!this.sticky)
            this.ZChat.style.opacity = "0.5"
        this.game.input.enabled = true
    }

    onSendMessage() {
        var textMessage = this.escapeHtml(this.txtInput.value)
        var mess = {
            from: this.game.Properties.pseudo,
            type: this.mode,
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