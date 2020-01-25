package routes

import (
	"github.com/gin-gonic/gin"
)

/*
Register shows the register page.
*/
func Register(c *gin.Context) {
	c.HTML(200, "register.tmpl", nil)
}
