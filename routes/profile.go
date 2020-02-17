package routes

import (
	"EPIC-Scouting/lib/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
Profile shows the profile page.
*/
func Profile(c *gin.Context) {
	d, _ := db.UserQuery("00000000-0000-0000-0000-000000000000") // TODO: Use the proper UserID
	c.HTML(http.StatusOK, "profile.tmpl", gin.H{"UserName": &d.UserName, "Email": &d.Email, "FirstName": &d.FirstName, "LastName": &d.LastName, "UserID": &d.UserID})
}

/*
ProfilePOST processes the user profile form.
*/
func ProfilePOST(c *gin.Context) {
	// TODO: Update profile based on user input. Very similar to the function in route "register" but calls db.UserModify() instead of db.UserCreate().
}
