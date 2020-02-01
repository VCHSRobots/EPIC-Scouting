package routes

import (
	"github.com/gin-gonic/gin"
)

/*
Data shows the data page.
*/
func Data(c *gin.Context) {
	c.HTML(200, "data.tmpl", nil)
}
