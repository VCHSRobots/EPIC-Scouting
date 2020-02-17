/*
Package auth TODO
*/
package auth

import (
	"EPIC-Scouting/lib/db"
	"strings"

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
	uuid, _ := DecodeLoginCookie(c)
	userData, _ := db.UserQuery(uuid)
	if userData == nil {
		return "guest"
	} else if userData.SysAdmin {
		return "sysadmin"
	}
	return "user"
}

//DecodeLoginCookie parses the login cookie into the username and uuid
func DecodeLoginCookie(c *gin.Context) (string, string) {
	cstring, _ := c.Cookie("login")
	if cstring == "" {
		return "", ""
	}
	carr := strings.Fields(cstring)
	return carr[0], carr[1]
}
