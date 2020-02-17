package routes

import (
	"EPIC-Scouting/lib/db"
	"EPIC-Scouting/lib/web"
	"fmt"

	"github.com/gin-gonic/gin"
)

/*
SysAdmin shows the SysAdmin page.
*/
func SysAdmin(c *gin.Context) {
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
	c.HTML(200, "sysAdmin.tmpl", gin.H{"DatabaseSizes": DatabaseSizes, "HeaderData": HeaderData})
}
