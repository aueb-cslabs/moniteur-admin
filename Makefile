.build-admin-win:
	cd adminPanel && npm install && npm run build:win

.build-admin-linux:
	cd adminPanel && npm install && npm run build:linux

.publish-admin-win:
	cd adminPanel && npm install && npm run publish:win

.publish-admin-linux:
	cd adminPanel && npm install && npm run publish:linux

default: .build-admin-linux
