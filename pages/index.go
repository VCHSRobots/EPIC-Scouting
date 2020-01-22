package pages

import (
	"github.com/gin-gonic/gin"
)

func init() {
	RegisterPage("index", VerbGET, showPageIndex)
}

func showPageIndex(c *gin.Context) {
	SendPage(c, nil, "index")
}
