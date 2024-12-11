package matchmaking

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

var users map[int64]string = make(map[int64]string)
var links map[string]int64 = make(map[string]int64)

func GenerateLink(userID int64) (string, error) {
	uuid := uuid.New()

	suuid := uuid.String()

	users[userID] = suuid
	links[suuid] = userID

	return suuid, nil
}

var ErrLinkNotFound = errors.New("Link not found")

func JoinLink(link string, userID int64) error {
	oppID, ok := links[link]
	if !ok {
		return ErrLinkNotFound
	}

	delete(links, link)
	delete(users, oppID)

	_, err := CreateGame(userID, oppID, sql.NullInt64{Valid: false})
	if err != nil {
		return err
	}

	err = ws.SendGameFound(oppID)
	return err
}
