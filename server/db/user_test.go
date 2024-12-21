package db

import (
	"log/slog"
	"testing"
	"time"
	"vezgammon/server/types"

	"golang.org/x/crypto/bcrypt"
	"gotest.tools/v3/assert"
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

func TestLoginUser(t *testing.T) {
	user := types.User{
		Username:  "testlogin",
		Firstname: "testlogin",
		Lastname:  "testlogin",
		Mail:      "testlogin",
	}
	password := "testlogin"
	hashpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NilError(t, err)

	retuser, _ := CreateUser(user, string(hashpass))
	retuser2, err := LoginUser(retuser.Username, password)
	assert.NilError(t, err)
	assert.Equal(t, retuser.ID, retuser2.ID)
}

func TestGenerateSessionToken(t *testing.T) {
	token := GenerateSessionToken()
	assert.Assert(t, len(token) > 0)
}

func TestSessionToken(t *testing.T) {
	user := types.User{
		Username:  "testsst",
		Firstname: "testsst",
		Lastname:  "testsst",
		Mail:      "testsst",
	}
	password := "testsst"
	hashpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NilError(t, err)

	retuser, _ := CreateUser(user, string(hashpass))

	token := GenerateSessionToken()
	err = SaveSessionToken(retuser.ID, token)
	assert.NilError(t, err)

	userid, err := ValidateSessionToken(token)
	assert.NilError(t, err)
	assert.Equal(t, retuser.ID, userid)

	err = Logout(token)
	assert.NilError(t, err)
}

func TestStatsAndBadge(t *testing.T) {
	user1 := types.User{
		Username:  "tseb1",
		Firstname: "tseb1",
		Lastname:  "tseb1",
		Mail:      "tseb1",
	}
	user2 := types.User{
		Username:  "tseb2",
		Firstname: "tseb2",
		Lastname:  "tseb2",
		Mail:      "tseb2",
	}
	password := "tseb"

	retuser1, _ := CreateUser(user1, password)
	retuser2, _ := CreateUser(user2, password)

	gameonline := types.Game{
		Player1: retuser1.ID,
		Player2: retuser2.ID,
		Start:   time.Now(),
		End:     time.Now(),
	}
	rgameonline, err := CreateGame(gameonline)
	assert.NilError(t, err)
	rgameonline.Status = types.GameStatusWinP1
	err = UpdateGame(rgameonline)
	assert.NilError(t, err)

	gamebot := types.Game{
		Player1: retuser1.ID,
		Player2: GetEasyBotID(),
		Start:   time.Now(),
		End:     time.Now(),
	}
	rgamebot, err := CreateGame(gamebot)
	assert.NilError(t, err)
	rgamebot.Status = types.GameStatusWinP1
	err = UpdateGame(rgamebot)

	_, err = GetStats(retuser1.ID)
	assert.NilError(t, err)
	_, err = GetStats(retuser2.ID)
	assert.NilError(t, err)
	_, err = GetBadge(retuser1.ID)
	assert.NilError(t, err)
	_, err = GetBadge(retuser2.ID)
	assert.NilError(t, err)
}

func TestChangeAvatar(t *testing.T) {
	user := types.User{
		Username:  "testavatarc",
		Firstname: "testavatarc",
		Lastname:  "testavatarc",
		Mail:      "testavatarc",
	}
	password := "testavatarc"

	retuser, err := CreateUser(user, password)
	assert.NilError(t, err)

	err = ChangeAvatar(retuser.ID, "https://api.dicebear.com/9.x/adventurer/svg?seed=Amaya")
	assert.NilError(t, err)
}

func TestChangePass(t *testing.T) {
	user := types.User{
		Username:  "testchangepass",
		Firstname: "testchangepass",
		Lastname:  "testchangepass",
		Mail:      "testchangepass",
	}
	password := "testchangepass"
	hashpass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NilError(t, err)

	retuser, err := CreateUser(user, string(hashpass))
	assert.NilError(t, err)

	newpass := "testchangepass2"
	err = ChangePass(retuser.Username, newpass, password)
	assert.NilError(t, err)
}
