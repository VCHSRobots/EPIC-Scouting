package routes

import (
	"github.com/gin-gonic/gin"
)

/*
About shows the About page.
*/
func About(c *gin.Context) {
	c.HTML(200, "about.tmpl", nil)
}
