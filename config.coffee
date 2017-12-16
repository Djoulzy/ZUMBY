module.exports = {
	config: {
		paths: {
			"watched": ["app"],
			"public": "public",
		},
		plugins: {
			babel: {
				ignore: /^(node_modules|vendor)/
			}
		},
		files: {
			javascripts: {
				joinTo: {
				'js/app.js': ['app/config.js',
					/^app\/client/],
				'js/crypt.js': /^app\/crypt/,
				'js/vendor.js': /(^node_modules|vendor)\//
				},
			},
			stylesheets: {
				joinTo: {
				'styles/login.css': 'app/css/login.less',
				'styles/chat.css': 'app/css/chat.less'
				}
			}
		},
		npm: {
			static: [
				'node_modules/crypto-js',
				'node_modules/phaser-ce/build/phaser.min.js',
				'node_modules/draggabilly/dist/draggabilly.pkgd.min.js'
			]
		},
		plugins: {
			babel: {
				ignore: /^(node_modules|vendor)/
			}
		}
	}
};
