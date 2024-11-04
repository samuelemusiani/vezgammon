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
    go run ./server ./server/config/test-config.toml

start-client:
    cd client && npm run dev

test: test-server test-client

test-server:
    sudo docker-compose -f docker-compose-test.yml up -d
    
    if go test -v ./server/... ; then \
        sudo docker-compose -f docker-compose-test.yml down; \
    else \
        sudo docker-compose -f docker-compose-test.yml down; \
        exit 1; \
    fi

test-client: 
    cd client && npm run test