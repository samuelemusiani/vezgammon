package matchmaking

import (
	//"errors"
	"database/sql"
	"log/slog"
	"math"
	"sync"
	"time"
	"vezgammon/server/types"
)

var mutex sync.Mutex

type DB interface {
	GetUser(int64) (*types.User, error)
	CreateGame(types.Game) (*types.Game, error)
}

type WS interface {
	SendGameFound(int64) error
}

var db DB
var ws WS

func Init(databse DB, websocket WS) {
	db = databse
	ws = websocket
	go worker(false)
}

func worker(finite bool) {
	var l int = 1
	for l > 0 {
		if finite {
			l--
		}

		if length() < 2 {
			time.Sleep(2 * time.Second)
			continue
		}

		mutex.Lock()
		for i := range length() - 1 {
			p1 := queue[(start+i)%qlen]
			for j := range length() - i {
				p2 := queue[(start+i+j+1)%qlen]
				if checkIfValidOpponent(p1.Elo, p2.Elo) {
					err := remove(p1)
					if err != nil {
						slog.With("err", err, "p1", p1).Error("Removing user1 from queue")
						continue
					}

					err = remove(p2)
					if err != nil {
						slog.With("err", err, "p2", p2).Error("Removing user2 from queue")
						continue
					}

					if _, err := CreateGame(p1.User_id, p2.User_id, sql.NullInt64{Valid: false}); err != nil {
						slog.With("err", err, "p1", p1.User_id, "p2", p2.User_id).Error("Creating game")
						continue
					}

					err = ws.SendGameFound(p1.User_id)
					if err != nil {
						slog.With("err", err).Error("Sending message to player")
					}

					err = ws.SendGameFound(p2.User_id)
					if err != nil {
						slog.With("err", err).Error("Sending message to player")
					}
					continue
				}
			}
		}
		mutex.Unlock()
		if !finite {
			time.Sleep(5 * time.Second)
		}
	}
}

func CreateGame(userID1, userID2 int64, tournament sql.NullInt64) (*types.Game, error) {
	user1, err := db.GetUser(userID1)
	if err != nil {
		return nil, err
	}

	if user1 == nil {
		panic("User1 is nil")
	}

	user2, err := db.GetUser(userID2)
	if err != nil {
		return nil, err
	}

	var dices = types.NewDices()
	var CurrentPlayer string
	if dices[0] >= dices[1] {
		CurrentPlayer = types.GameCurrentPlayerP1
	} else {
		CurrentPlayer = types.GameCurrentPlayerP2
	}

	var game types.Game
	game.Player1 = user1.ID
	game.Player2 = user2.ID
	game.Elo1 = user1.Elo
	game.Elo2 = user2.Elo
	game.Start = time.Now()
	game.End = time.Now()
	game.CurrentPlayer = CurrentPlayer
	game.Dices = types.NewDices()

	game.Tournament = tournament

	retgame, err := db.CreateGame(game)
	if err != nil {
		return nil, err
	}

	return retgame, nil
}

func StopSearch(uid int64) error {
	mutex.Lock()
	defer mutex.Unlock()
	err := remove(qel{User_id: uid})
	return err
}

func SearchGame(uid int64) error {
	u, err := db.GetUser(uid)
	if err != nil {
		return err
	}
	slog.With("user stats: ", u).Debug("User stats")

	mutex.Lock()
	err = push(qel{User_id: u.ID, Elo: u.Elo})
	if err != nil {
		slog.With("err", err).Error("Pushing into queue")
	}
	mutex.Unlock()

	//cerco l'opponent nel db in base al player e in base a quanto è in queue
	return err
}

func checkIfValidOpponent(elo1, elo2 int64) bool {
	diff := elo1 - elo2
	if diff < 0 {
		diff = -diff
	}
	return math.Abs(float64(elo1-elo2)) < 200
}
