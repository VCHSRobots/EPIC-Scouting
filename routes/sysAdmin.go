package routes

import (
	"EPIC-Scouting/lib/auth"
	"EPIC-Scouting/lib/db"
	"EPIC-Scouting/lib/web"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
SysAdmin shows the SysAdmin page.
*/
func SysAdmin(c *gin.Context) {
	print(auth.GetUserMode(c))
	if auth.GetUserMode(c) != "sysadmin" {
		Forbidden(c)
		return
	}
	HeaderData := &web.HeaderData{Title: "Super Secret Sysadmin Bunker", StyleSheets: []string{"global"}}
	sizes := db.GetDatabaseSize()
	var DatabaseSizes []string
	var totalSize float64 = 0
	for dbName, dbSize := range sizes { // TODO: Sort database names alphabetically.
		size := float64(dbSize / 1000) // Convert from bytes to kilobytes
		totalSize += size
		DatabaseSizes = append(DatabaseSizes, fmt.Sprintf("%s: %v KB", dbName, size))
	}
	DatabaseSizes = append(DatabaseSizes, fmt.Sprintf("%s: %v KB", "Total", totalSize))
	var SysAdmins []string
	adminlist := db.SysAdminList()
	for id, name := range adminlist {
		SysAdmins = append(SysAdmins, fmt.Sprintf("%s - %s", id, name))
	}
	userlist := db.UserList()
	var Users []string
	for id, name := range userlist {
		Users = append(Users, fmt.Sprintf("%s - %s", id, name))
	}
	var Campaigns []string
	campaignList := db.CampaignList()
	for id, details := range campaignList {
		Campaigns = append(Campaigns, fmt.Sprintf("%s - %s (Owned by team %s)", id, details[1], details[0]))
	}
	var Teams []string
	teamList := db.TeamListFull()
	for id, details := range teamList {
		Teams = append(Teams, fmt.Sprintf("%s - %s - %s (Scouting match %s at event TODO for campaign TODO)", id, details[0], details[1], details[2]))
	}
	c.HTML(200, "sysAdmin.tmpl", gin.H{"DatabaseSizes": DatabaseSizes, "SysAdmins": SysAdmins, "Users": Users, "Campaigns": Campaigns, "Teams": Teams, "HeaderData": HeaderData})
}

/*
SysAdminToggle toggles a user's status as SysAdmin
*/
func SysAdminToggle(c *gin.Context) {
	print(auth.GetUserMode(c))
	if auth.GetUserMode(c) != "sysadmin" {
		Forbidden(c)
		return
	}
	c.Request.ParseForm()
	id := c.PostForm("toggleSysAdmin")
	user, _ := db.UserQuery(id)
	if user.SysAdmin {
		db.SysAdminDemote(id)
	} else {
		db.SysAdminPromote(id)
	}
	c.Redirect(http.StatusSeeOther, "/sysadmin")
}
