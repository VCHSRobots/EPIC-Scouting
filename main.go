package main

import (
	"EPIC-Scouting/lib/config"
	"EPIC-Scouting/lib/db"
	"EPIC-Scouting/lib/lumberjack"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine
var buildName string = "Prerelease 0.1" // TODO: Update this with each release. "Prerelease" for development versions.
var buildDate string = "20.008.1"       // TODO: Update this with each release. External script updates number and it is sourced from an external file. Format: YY.DDD.N, where YY is year, DDD is day of year, and N is build number in day.

func apiHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello %s.", r.URL.Path[1:])
}

func main() {
	configuration := config.Load()
	db.TouchBase(configuration.DatabasePath)
	lumberjack.Start(configuration.LogPath, configuration.Verbosity)
	log := lumberjack.New("")
	log.Infof("Scouting system started. Version: %s (%s)", buildName, buildDate)

	log.Fatal(http.ListenAndServe(":4415", nil))
	// TODO: Graceful shutdown and restart procedure.
}
