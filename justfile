
default: build


build: build-client copy-client build-server

start: build
    cd server && ./server

build-server:
    cd server && go build

start-server: build-server
    cd server && ./server

build-client:
    cd client && npm run build

start-client:
    cd client && npm run dev

copy-client:
    cp -r client/dist server/dist
