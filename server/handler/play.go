package handler

import "github.com/gin-gonic/gin"

// @Summary Start a matchmaking search for a new game
// @Schemes
// @Description Start a matchmaking search for a new game
// @Tags play
// @Accept json
// @Produce json
// @Success 201 {string} Started search
// @Router /play/search [get]
func StartPlaySearch(c *gin.Context) {
}

// @Summary Stop a running matchmaking search
// @Schemes
// @Description Stop a running matchmaking search
// @Tags play
// @Accept json
// @Produce json
// @Success 204 {string} Stopped search
// @Router /play/search [delete]
func StopPlaySearch(c *gin.Context) {
}

// @Summary Get current game
// @Schemes
// @Description Get current game
// @Tags play
// @Accept json
// @Produce json
// @Success 200 {object} types.Game
// @Router /play [get]
func GetCurrentGame(c *gin.Context) {
}

// @Summary Surrend to current game
// @Schemes
// @Description Surrend to current game
// @Tags play
// @Accept json
// @Produce json
// @Success 201 {string} Surrended to game
// @Router /play [delete]
func SurrendToCurrentGame(c *gin.Context) {
}

// @Summary Get possible moves for next turn
// @Schemes
// @Description Get possible moves for next turn
// @Tags play
// @Accept json
// @Produce json
// @Success 200 {string} Possible moves
// @Router /play/moves [get]
func GetPossibleMoves(c *gin.Context) {
}

// @Summary Play all the moves for the current turn
// @Schemes
// @Description Play all the moves for the current turn
// @Tags play
// @Accept json
// @Produce json
// @Success 201 {string} Moves played
// @Router /play/moves [post]
func PlayMoves(c *gin.Context) {
}

// @Summary The player want to double
// @Schemes
// @Description The player want to double
// @Tags play
// @Accept json
// @Produce json
// @Success 201 {string} Double requested
// @Router /play/double [post]
func WantToDouble(c *gin.Context) {
}

// @Summary Refuse the double
// @Schemes
// @Description Refuse the double
// @Tags play
// @Accept json
// @Produce json
// @Success 201 {string} Double refused
// @Router /play/double [delete]
func RefuseDouble(c *gin.Context) {
}

// @Summary Accept the double
// @Schemes
// @Description Accept the double
// @Tags play
// @Accept json
// @Produce json
// @Success 201 {string} Double accepted
// @Router /play/double [put]
func AcceptDouble(c *gin.Context) {
}
