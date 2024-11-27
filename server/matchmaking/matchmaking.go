package matchmaking

import (
	//"errors"
	"log/slog"
	"sync"
	"time"
	"vezgammon/server/db"
	"vezgammon/server/types"
	"vezgammon/server/ws"
)

var mutex sync.Mutex

func Init() {
	go func() {
		for {
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

						if err := createGame(p1.User_id, p2.User_id); err != nil {
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
			time.Sleep(5 * time.Second)
		}
	}()
}

func createGame(user_id1, user_id2 int64) error {
	user1, err := db.GetUser(user_id1)
	if err != nil {
		return err
	}

	user2, err := db.GetUser(user_id2)
	if err != nil {
		return err
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

	_, err = db.CreateGame(game)
	if err != nil {
		return err
	}

	return nil
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
	return (elo1 - elo2) < 200
}
