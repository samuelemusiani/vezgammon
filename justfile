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
    rm -rf server/dist
    cp -r client/dist server/dist

start-server: generate-swag
    go run ./server ./server/config/test-config.toml

start-client:
    cd client && npm i && npm run dev

test: test-server

test-server:
    #!/usr/bin/env sh
    sudo docker-compose -f docker-compose-test.yml up -d
    sleep 2
    if go test -coverprofile server/coverage.txt ./server/... -v -p 1; then
        sudo docker-compose -f docker-compose-test.yml down -v;
    else
        sudo docker-compose -f docker-compose-test.yml down -v;
        exit 1;
    fi

generate-swag: install-swag
    cd server && ~/go/bin/swag init

install-swag:
    go install github.com/swaggo/swag/cmd/swag@latest

test-client:
    cd client && npm run test
