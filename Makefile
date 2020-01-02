.build-admin-win:
	npm install && npm run build:win

.build-admin-linux:
	npm install && npm run build:linux

.publish-admin-win:
	npm install && npm run publish:win

.publish-admin-linux:
	npm install && npm run publish:linux

default: .build-admin-linux
