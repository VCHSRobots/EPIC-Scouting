/*
Package auth TODO
*/
package auth

import (
	"EPIC-Scouting/lib/db"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

/*
GetUserMode TODO. For front-end web use.
*/
func GetUserMode(c *gin.Context) string {
	// TODO:
	// This function is used when rendering menu options on the index.
	// 1. Check if user is logged in. If yes, return "user". Otherwise "guest"
	// Basically looks at their cookies / auth token.
	uuid := CheckLogin(c)
	userData, _ := db.UserQuery(uuid)
	if userData == nil {
		return "guest"
	} else if userData.SysAdmin {
		return "sysadmin"
	}
	return "user"
}

/*
SetLogin sets a login cookie
*/
func SetLogin(c *gin.Context, userID string) {
	session := sessions.Default(c)
	session.Set("userID", userID)
	session.Save()
}

/*
CheckLogin gets the currently logged in uuid
*/
func CheckLogin(c *gin.Context) string {
	session := sessions.Default(c)
	userID := session.Get("userID")
	if userID == nil {
		return ""
	}
	return userID.(string)
}

/*
SetTeam sets the team cookie
*/
func SetTeam(c *gin.Context, team string) {
	session := sessions.Default(c)
	session.Set("team", team)
	session.Save()
}

/*
CheckTeam gets the team cookie
*/
func CheckTeam(c *gin.Context) string {
	session := sessions.Default(c)
	team := session.Get("userID")
	if team == nil {
		return ""
	}
	return team.(string)
}
