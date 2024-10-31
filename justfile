build: clean build-client copy-client build-server

clean:
  rm -rf ./build/vezgammon
  rm -rf ./server/dist

build-server:
    mkdir -p build
    go build -o ./build/vezgammon ./server
    cp -u ./server/config/config.toml ./build

build-client:
    cd client && npm install && npm run build

copy-client:
    cp -r client/dist server/dist

start-server:
    go run ./server ./server/config.toml

start-client:
    cd client && npm run dev
