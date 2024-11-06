# VezGammon
### Team 1: Diego Barbieri, Lorenzo Peronese, Samuele Musiani, Fabio Murer, Emanuele Argonni, Omar Ayache 
## Description

Vezgammon is a web-based bologna-style backgammon game that allows users to play against each other in real-time. The game is built in Vue and Go. The game is designed to be simple and easy to use, with a clean and intuitive interface. Players can create a new game, join an existing game, or watch a game in progress. The game also features a chat system that allows players to communicate with each other during the game.


## requirements

- `docker` with `compose` plugin (newer) or  `docker-compose` (older)
- `golang compiler` version `1.23.2`
- `npm`
- `just`

## Testing 

### server 

```bash
just test-server
```

## set up backend developing environment

- start the database locally with `sudo docker-compose -f docker-compose-test.yml up -d`
- run `just start-server` as many times as you want
- when done stop the database with `sudo docker-compose -f docker-compose-test.yml down`


## Deployment

- create database password
    - create db directory `mkdir db`
    - create `password.txt` containing a password es. `echo -n $(dd if=/dev/random bs=1 count=32 | base32 | sed \'s/=//g\') > db/password.txt` 

- start services with `sudo docker-compose -f docker-compose-release.yml up -d --build`

stops services with `sudo docker-compose -f docker-compose-release.yml down` add `-v` flag to reset db

## swagger

swagger documentation is generated with [swag](https://github.com/swaggo/swag)

- install it with `go install github.com/swaggo/swag/cmd/swag@latest`
- create documentation by commenting handler function with [this syntax](https://github.com/swaggo/swag?tab=readme-ov-file#declarative-comments-format)
- generate documentation with `cd server && ~/go/bin/swag init`

when running `go start-server` this is done automatically

swagger ui is accessible at [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) when server is started with `swagger = true` configuration

## bgweb-api

[bgweb-api](https://github.com/foochu/bgweb-api) is a evalution engine for Backgammon moves. we use to
- get a list of legal moves
- creating game bots of various difficulties
- analyze games

the interface to this api is coded in bgweb go module

## docker and docker compose problems

- __requirements__
    - docker and docker compose installed (for debian `docker.io`, `docker-compose` packages)

justfile scripts are build for debian witch has `docker-compose` older package instead of `docker compose` docker module
if you have problems with `docker-compose` replace it with `docker compose`

if a docker compose starts with error try to stop all docker containers with `sudo docker stop $(sudo docker ps -a -q)` and delete all volumes with `sudo docker volume rm $(sudo docker volume ls -q)`