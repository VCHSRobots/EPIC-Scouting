package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//MatchData struct for recieving json data for matches
type MatchData struct {
	Data [][]string `json:"data"`
}

/*
Scout shows the scout page.
*/
func Scout(c *gin.Context) {
	c.HTML(http.StatusOK, "scout.tmpl", nil)
}

//MatchPOST processes and stores scouting data from a match
func MatchPOST(c *gin.Context) {
	var data MatchData
	c.ShouldBindJSON(&data)
	postData := data.Data[0]
	for _, str := range postData {
		fmt.Println(str)
	}
}

//PitPOST processes and stores/updates pit scouting data for a team
