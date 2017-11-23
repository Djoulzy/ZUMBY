
module.exports = {
	Encrypt_b64: function(text) {

  		var CryptoJS = require("crypto-js");
		
		HASH_SIZE = 8
		HEX_KEY = "d87fbb277eefe245ee384b6098637513462f5151336f345778706b462f724473"
		HEX_IV = "046b51957f00c25929e8ccaad3bfe1a7"

		text_bin = CryptoJS.enc.Utf8.parse(text)
		key_bin = CryptoJS.enc.Hex.parse(HEX_KEY)
		iv_bin = CryptoJS.enc.Hex.parse(HEX_IV)

		hash = CryptoJS.MD5(text_bin).toString().substr(0, 16)
		console.log("Hash: "+hash)
		signedText = CryptoJS.enc.Hex.parse(hash + text_bin)
		console.log("SignedText: " + signedText)

		encrypted = CryptoJS.AES.encrypt(signedText, key_bin, { iv: iv_bin, mode: CryptoJS.mode.CBC, padding: CryptoJS.pad.Pkcs7 }).ciphertext
		b64_iv = CryptoJS.enc.Base64.stringify(iv_bin)
		b64_crypted = CryptoJS.enc.Base64.stringify(encrypted)
		console.log("IV = " + b64_iv)
		console.log("encrypted = " + b64_crypted)

		b64_iv_final = rtrim(b64_iv.replaceAll("/", "_").replaceAll("+", "-"), "=")
		b64_crypted_final = rtrim(b64_crypted.replaceAll("/", "_").replaceAll("+", "-"), "=")

		console.log("Final = " + b64_iv_final + "/" + b64_crypted_final)
		return b64_iv_final + "/" + b64_crypted_final
	}
}

function rtrim(str, chr) {
	var rgxtrim = (!chr) ? new RegExp('\\s+$') : new RegExp(chr+'+$');
	return str.replace(rgxtrim, '');
}

String.prototype.replaceAll = function(searchStr, replaceStr) {
	var str = this;
    if(str.indexOf(searchStr) === -1) {
        return str;
    }
    return (str.replace(searchStr, replaceStr)).replaceAll(searchStr, replaceStr);
}
