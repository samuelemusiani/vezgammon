package bgweb

import (
	"encoding/json"
	"testing"
)

var domain = "localhost:3030/api/v1/"

func TestGetmoves(t *testing.T) {
	jsonArg := `
{
  "board": {
    "o": {
      "6": 5,
      "8": 3,
      "13": 5,
      "24": 2
    },
    "x": {
      "6": 5,
      "8": 3,
      "13": 5,
      "24": 2
    }
  },
  "cubeful": false,
  "dice": [
    3,
    1
  ],
  "max-moves": 3,
  "player": "x",
  "score-moves": true
}
	`

	moveargs := DefaultMoveArgs
	err := json.Unmarshal([]byte(jsonArg), &moveargs)
	if err != nil {
		t.Fatalf("cannot extract from json")
	}
	t.Logf("moveargs %v", moveargs)

}
