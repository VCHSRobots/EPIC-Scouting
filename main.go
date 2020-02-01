package main

import (
	"EPIC-Scouting/lib/config"
	"EPIC-Scouting/lib/db"
	"EPIC-Scouting/lib/lumberjack"
	"EPIC-Scouting/routes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

var configuration config.YAML
var router *gin.Engine

func main() {
	configuration = config.Load()
	db.TouchBase(configuration.DatabasePath)
	lumberjack.Start(configuration.LogPath, configuration.Verbosity)
	log := lumberjack.New("Main")
	buildName, buildDate := config.BuildInformation()
	log.Infof("Scouting system started. Version: %s (%s)", buildName, buildDate)
	go start(configuration.Port)
	// Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Warn("Scouting system shutting down.")
}

func start(port int) {
	log := lumberjack.New("Router")
	if port <= 0 {
		port = 443
	}
	address := fmt.Sprintf(":%d", port)
	if configuration.Verbosity < 1 {
		gin.SetMode(gin.ReleaseMode)
	}

	// TODO BEGIN: Integrate Gin log messages under their appropriate logLevels with *lumberjack*

	ginLogFile, errOpen := os.OpenFile(configuration.LogPath+"access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if errOpen != nil {
		fmt.Fprintf(os.Stderr, "Unable to create log file: %s", errOpen.Error())
		os.Exit(1)
	}
	defer ginLogFile.Close()
	outlog := io.MultiWriter(ginLogFile, os.Stdout)
	errlog := io.MultiWriter(ginLogFile, os.Stderr)
	gin.DefaultWriter = outlog
	gin.DefaultErrorWriter = errlog

	// TODO END

	router = gin.New()
	router.Use(gin.Logger(), gin.Recovery()) // TODO: Add recovery middleware and authentication middle which refreshes session tokens.
	router.Static("/css", "./static/css")    // TODO: Make URL access to /css/, /js/, /media/, and /templates/ use the proper NoRoute() handler and NOT http.404, as it currently just returns a blank page.
	router.Static("/js", "./static/js")
	router.Static("/media", "./static/media")
	router.Static("/templates", "./static/templates")
	router.LoadHTMLGlob("./static/templates/*") // Load templates
	router.NoRoute(func(c *gin.Context) { c.HTML(http.StatusNotFound, "404.tmpl", nil) })
	router.NoMethod(func(c *gin.Context) { c.HTML(http.StatusMethodNotAllowed, "405.tmpl", nil) })
	// TODO: Add handlers for 401, 403, 500 codes.
	// TODO: Dynamically load routes from files in "/routes/" instead of hard-coding them.
	router.GET("/", routes.Index)
	router.GET("/about", routes.About)
	router.GET("/login", routes.Login)
	router.POST("/loginPOST", routes.LoginPOST)
	router.GET("/register", routes.Register)
	router.GET("/scout", routes.Scout)
	router.GET("/dashboard", routes.Dashboard)
	router.GET("/sysadmin", routes.SysAdmin)
	log.Debugf("Serving on port %d.", port)
	log.Fatal(router.Run(address))
}
