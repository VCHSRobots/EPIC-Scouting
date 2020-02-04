package routes

import (
	"github.com/gin-gonic/gin"
)

/*
SysAdmin shows the SysAdmin page.
*/
func SysAdmin(c *gin.Context) {
	c.HTML(200, "sysAdmin.tmpl", nil)
}
