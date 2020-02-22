package routes

import (
	"EPIC-Scouting/lib/auth"
	"EPIC-Scouting/lib/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Profile shows the profile page.
*/
func Profile(c *gin.Context) {
	uuid, _ := auth.LoginCookie(c)
	//checks if login uuid was valid
	if uuid == "" {
		Forbidden(c)
		return
	}
	d, _ := db.UserQuery(uuid) // TODO: Use the proper UserID
	NullString := "{ false}"
	Email := ""
	FirstName := ""
	LastName := ""
	if d.Email != NullString {
		Email = d.Email
	}
	if d.FirstName != NullString {
		FirstName = d.FirstName
	}
	if d.LastName != NullString {
		LastName = d.LastName
	}
	c.HTML(http.StatusOK, "profile.tmpl", gin.H{"Username": d.UserName, "Email": Email, "FirstName": FirstName, "LastName": LastName, "UserID": d.UserID})
}

/*
ProfilePOST processes the user profile form.
*/
func ProfilePOST(c *gin.Context) {
	// TODO: Update profile based on user input. Very similar to the function in route "register" but calls db.UserModify() instead of db.UserCreate().
}
