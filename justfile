build: build-client copy-client build-server
  cp ./server/config.toml ./build

build-server:
    mkdir build
    go build -o ./build/vezgammon ./server

build-client:
    cd client && npm install && npm run build

copy-client:
    cp -r client/dist server/dist

start-server:
    go run ./server ./server/config.toml

start-client:
    cd client && npm run dev
