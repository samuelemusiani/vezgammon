package db

import (
	"log/slog"
	"testing"
	"vezgammon/server/types"
)

func TestInitUser(t *testing.T) {
	err := initUser()
	if err != nil {
		t.Fatalf("cannot initialize user, %s", err)
	}

	q := `
	SELECT EXISTS (
    SELECT FROM 
        pg_tables
    WHERE  
        tablename  = 'users'
    )
	`

	row := conn.QueryRow(q)

	var isin bool
	err = row.Scan(&isin)
	if err != nil {
		t.Fatalf("%s", err)
	}

	if !isin {
		t.Fatalf("user table doesn't exists")
	}
}

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
	q := `DELETE FROM users`
	_, err := conn.Exec(q)
	if err != nil {
		t.Fatalf("%s", err)
	}

	u1 := types.User{
		Username:  "gigi",
		Firstname: "giorgio",
		Lastname:  "galli",
		Mail:      "giorgio.galli@mail.it",
	}
	pass1 := "asda13"

	u2 := types.User{
		Username:  "sa",
		Firstname: "salvatore",
		Lastname:  "esposito",
		Mail:      "salvatore.esposito@mail.it",
	}
	pass2 := "dfa24"

	retu1, _ := CreateUser(u1, pass1)
	retu2, _ := CreateUser(u2, pass2)
	u1.ID = retu1.ID
	u2.ID = retu2.ID

	users, err := GetUsers()
	if err != nil {
		t.Fatalf("error getting users %s", err)
	}

	slog.With("users", users).Debug("users array")
	if users[0] != u1 || users[1] != u2 {
		t.Fatalf("not getting users correctly")
	}
}
