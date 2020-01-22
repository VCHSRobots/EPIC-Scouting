package main

import (
	"EPIC-Scouting/lib/config"
	"EPIC-Scouting/lib/db"
	"EPIC-Scouting/lib/lumberjack"
	"EPIC-Scouting/pages"
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
var buildName string = "Prerelease 0.1" // TODO: Update this with each release. "Prerelease" for development versions.
var buildDate string = "20.021.1"       // TODO: Update this with each release. External script updates number and it is sourced from an external file. Format: YY.DDD.N, where YY is year, DDD is day of year, and N is build number in day.

func main() {
	configuration = config.Load()
	db.TouchBase(configuration.DatabasePath)
	lumberjack.Start(configuration.LogPath, configuration.Verbosity)
	log := lumberjack.New("Main")
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
	router.Use(gin.Logger())
	router.Static("/css", "./static/css") // TODO: Make URL access to /css/, /js/, /media/, and /templates/ use the proper NoRoute() handler and NOT http.404, as it currently just returns a blank page.
	router.Static("/js", "./static/js")
	router.Static("/media", "./static/media")
	router.Static("/templates", "./static/templates")
	router.LoadHTMLGlob("./static/templates/*") // Load templates
	router.NoRoute(func(c *gin.Context) { c.HTML(http.StatusNotFound, "404.tmpl", nil) })
	router.NoMethod(func(c *gin.Context) { c.HTML(http.StatusMethodNotAllowed, "405.tmpl", nil) })
	pageList := pages.GetPages()
	for _, page := range pageList {
		if page.Verb == pages.VerbGET {
			router.GET(page.Route, page.Handlers...)
		}
		if page.Verb == pages.VerbPOST {
			router.POST(page.Route, page.Handlers...)
		}
	}
	log.Debugf("Loaded %d pages.", len(pageList))
	log.Debugf("Serving on port %d.", port)
	log.Fatal(router.Run(address))
}
