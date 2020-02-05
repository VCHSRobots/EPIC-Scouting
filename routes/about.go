package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

/*
About shows the About page.
*/
func About(c *gin.Context) {
	c.HTML(http.StatusOK, "about.tmpl", nil)
}
