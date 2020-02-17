package main

import (
	"EPIC-Scouting/lib/config"
	"EPIC-Scouting/lib/db"
	"EPIC-Scouting/lib/lumberjack"
	"EPIC-Scouting/routes"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	nice "github.com/ekyoung/gin-nice-recovery"
	"github.com/gin-contrib/gzip"

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
	router.Use(gin.Logger())                    // TODO: Add authentication middle which refreshes session tokens.
	router.Use(gzip.Gzip(gzip.BestCompression)) // Gzip compression.
	router.Static("/css", "./static/css")       // TODO: Make URL access to /css/, /js/, /media/, and /templates/ use the proper NoRoute() handler and NOT http.404, as it currently just returns a blank page.
	router.Static("/js", "./static/js")
	router.Static("/media", "./static/media")
	router.Static("/templates", "./static/templates")
	router.LoadHTMLGlob("./static/templates/*") // Load templates.
	// TODO: Dynamically load routes from files in "/routes/" instead of hard-coding them.
	router.NoRoute(routes.NotFound)                       // 404.
	router.NoMethod(routes.MethodNotAllowed)              // 405.
	router.Use(nice.Recovery(routes.InternalServerError)) // 500.
	router.GET("/", routes.Index)
	router.GET("/about", routes.About)
	router.GET("/dashboard", routes.Dashboard)
	router.GET("/help", routes.Help)
	router.GET("/login", routes.Login)
	router.POST("/loginPOST", routes.LoginPOST)
	router.GET("/logout", routes.Logout)
	router.GET("/profile", routes.Profile)
	router.POST("/profilePOST", routes.ProfilePOST)
	router.GET("/register", routes.Register)
	router.POST("/registerPOST", routes.RegisterPOST)
	router.GET("/scout", routes.Scout)
	router.POST("/matchPOST", routes.MatchPOST)
	router.GET("/sysadmin", routes.SysAdmin)
	router.GET("/teamJoin", routes.TeamJoin)
	router.GET("/teamCreate", routes.TeamCreate)
	router.GET("/teamData", routes.TeamData)
	router.GET("/teamAdmin", routes.TeamAdmin)
	router.GET("/data", routes.Data)
	router.GET("/getGraph", routes.GetGraph)
	log.Debugf("Serving on port %d.", port)
	log.Fatal(router.Run(address))
}
