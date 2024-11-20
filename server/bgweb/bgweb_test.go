package bgweb

import (
	"encoding/json"
	"testing"

	"gotest.tools/v3/assert"
)

func TestGetmoves(t *testing.T) {

	url = "localhost:3030/api/v1/"

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

	jsonresp := `
   [
   {
     "evaluation": {
       "diff": 0,
       "eq": 0.159,
       "info": {
         "cubeful": false,
         "plies": 1
       },
       "probability": {
         "lose": 0.449,
         "loseBG": 0.005,
         "loseG": 0.124,
         "win": 0.551,
         "winBG": 0.013,
         "winG": 0.174
       }
     },
     "play": [
       {
         "from": "8",
         "to": "5"
       },
       {
         "from": "6",
         "to": "5"
       }
     ]
   },
   {
     "evaluation": {
       "diff": -0.168,
       "eq": -0.009,
       "info": {
         "cubeful": false,
         "plies": 1
       },
       "probability": {
         "lose": 0.503,
         "loseBG": 0.007,
         "loseG": 0.14,
         "win": 0.497,
         "winBG": 0.008,
         "winG": 0.137
       }
     },
     "play": [
       {
         "from": "13",
         "to": "10"
       },
       {
         "from": "24",
         "to": "23"
       }
     ]
   },
   {
     "evaluation": {
       "diff": -0.175,
       "eq": -0.015,
       "info": {
         "cubeful": false,
         "plies": 1
       },
       "probability": {
         "lose": 0.503,
         "loseBG": 0.004,
         "loseG": 0.135,
         "win": 0.497,
         "winBG": 0.005,
         "winG": 0.125
       }
     },
     "play": [
       {
         "from": "24",
         "to": "21"
       },
       {
         "from": "21",
         "to": "20"
       }
     ]
   }
 ]
   `

	moveargs := DefaultMoveArgs
	err := json.Unmarshal([]byte(jsonArg), &moveargs)
	assert.NilError(t, err)
	t.Logf("moveargs %v", moveargs)

	var respmoves []Move
	err = json.Unmarshal([]byte(jsonresp), &respmoves)
	assert.NilError(t, err)
	t.Logf("respmoves %v", respmoves)

	rp, err := GetMoves(&moveargs)
	assert.NilError(t, err)

	assert.DeepEqual(t, rp, respmoves)
}
