package db

import (
	"log/slog"
	"testing"
	"vezgammon/server/types"
)

func TestCreateUser(t *testing.T) {

	u := types.User{
		Username:  "gigi",
		Firstname: "giorgio",
		Lastname:  "galli",
		Mail:      "giorgio.galli@mail.it",
	}

	password := "2354wdfs"

	retuser, err := CreateUser(u, password)
	if err != nil {
		t.Fatalf("not creating user %s", err)
	}

	if retuser.ID == u.ID {
		t.Fatalf("id not set %d", retuser.ID)
	}
}

func TestGetUsers(t *testing.T) {

	u2 := types.User{
		Username:  "sa",
		Firstname: "salvatore",
		Lastname:  "esposito",
		Mail:      "salvatore.esposito@mail.it",
	}
	pass2 := "dfa24"

	retu2, _ := CreateUser(u2, pass2)
	u2.ID = retu2.ID

	users, err := GetUsers()
	if err != nil {
		t.Fatalf("error getting users %s", err)
	}

	slog.With("users", users).Debug("users array")
	if len(users) < 2 {
		t.Fatalf("not getting users correctly")
	}
}
