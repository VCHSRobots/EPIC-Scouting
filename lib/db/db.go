/*
Package db provides tools for creating and interacting with the SQLite databases.
*/
package db

import (
	"EPIC-Scouting/lib/lumberjack"
	"database/sql"
	"os"

	"github.com/google/uuid"
	// TODO: Golint needs to stop complaining about this import. >:[
	_ "github.com/mattn/go-sqlite3"
)

var log = lumberjack.New("DB")

/*
TouchBase creates all databases used by the server if they do not exist.
Its name is a play on the GNU program "touch", the idiom "[to] touch base", and the word "database". The author is rather proud of this.
*/
func TouchBase(databasePath string) {
	newDatabase := func(databaseName string) (database *sql.DB) {
		var err error
		db, err := sql.Open("sqlite3", databasePath+databaseName+".db")
		if err != nil {
			log.Fatal("Unable to open or create database: " + err.Error())
		}
		return db
	}

	err := os.MkdirAll(databasePath, 0755)
	if err != nil {
		log.Fatal("Unable to create database: " + err.Error())
	}

	// Users
	// This database stores all users.

	users := newDatabase("users")
	users.Exec("CREATE TABLE IF NOT EXISTS users ( userid TEXT PRIMARY KEY UNIQUE NOT NULL, username TEXT NOT NULL UNIQUE, password TEXT NOT NULL, firstname TEXT, lastname TEXT, email TEXT UNIQUE, phone TEXT UNIQUE, usertype TEXT)") // TODO: Add support for N+ contact options; via linked table?

	// TODO: Create a SYSTEM team which makes the default public campaigns each season.

	// Scouting teams

	teams := newDatabase("teams")
	teams.Exec("CREATE TABLE IF NOT EXISTS teams ( teamid TEXT PRIMARY KEY UNIQUE NOT NULL, number TEXT UNIQUE, name TEXT NOT NULL, currentcampaign TEXT NOT NULL )")                                                                              // A team.
	teams.Exec("CREATE TABLE IF NOT EXISTS members ( teamid TEXT PRIMARY KEY NOT NULL, userid TEXT NOT NULL, usertype TEXT NOT NULL )")                                                                                                            // The members on a team. UserType is either member or admin.
	teams.Exec("CREATE TABLE IF NOT EXISTS participating ( teamid TEXT PRIMARY KEY NOT NULL, eventid TEXT NOT NULL )")                                                                                                                             // What events a team is participating in. If a team is currently running a campaign, they must have *some* event they are participating in. A team is scouting all matches during an event, of course.
	teams.Exec("CREATE TABLE IF NOT EXISTS results ( campaignid TEXT PRIMARY KEY NOT NULL, eventid TEXT NOT NULL, matchid TEXT NOT NULL, competitorid TEXT NOT NULL, teamid TEXT NOT NULL, userid TEXT NOT NULL, datetime TEXT NOT NULL, stats )") // A team's scouted results. Any number of teams may scout for the same campaign / event / match at the same time.

	// Campaigns (game seasons / years)
	// This database stores the expected campaign / event / match schedule and data, and expected competing teams. Data pulled from TBA.

	campaigns := newDatabase("campaigns")
	campaigns.Exec("CREATE TABLE IF NOT EXISTS campaigns ( campaignid TEXT PRIMARY KEY UNIQUE NOT NULL, owner TEXT NOT NULL, name TEXT NOT NULL )")                                      // TODO: Add more information about each campaign. Campaign owner is a teamid. If campaign owner is all zeros, campaign is global.
	campaigns.Exec("CREATE TABLE IF NOT EXISTS events ( campaignid TEXT PRIMARY KEY NOT NULL, eventid TEXT NOT NULL, name TEXT NOT NULL, location TEXT, starttime TEXT, endtime TEXT )") // TODO: Add more information about each event.
	campaigns.Exec("CREATE TABLE IF NOT EXISTS matches ( eventid TEXT PRIMARY KEY NOT NULL, matchid TEXT UNIQUE NOT NULL, matchnumber INTEGER NOT NULL, starttime TEXT, endtime TEXT )") // TODO: Add more information about each match.
	campaigns.Exec("CREATE TABLE IF NOT EXISTS participants ( matchid TEXT PRIMARY KEY NOT NULL, competitorid TEXT UNIQUE NOT NULL) ")                                                   // The participants in each match.

	campaigns.Exec("CREATE TABLE IF NOT EXISTS competitors ( competitorid TEXT PRIMARY KEY UNIQUE NOT NULL, number TEXT UNIQUE, name TEXT NOT NULL )") // TODO: Add more information about each competing team.

}

/*
CreateUser creates a new user.
*/
func CreateUser(username, password, firstname, lastname, email, phone string) {
	id := uuid.New()
	log.Debugf("Created user: %s", id.String())
	// Throw if overwriting existing
}

/*
CreateTeam TODO
*/
func CreateTeam() {
	// TODO
	// Req: TeamNumber, TeamName, etc
	// Throw if existing. If a team's number exists in campaigns/competitors, their competitorid is their new teamid. If a team's number exists in teams/teams, any reference to their competitorid in campaigns/competitors is their teamid.
}

/*
Results TODO
*/
func WriteResults() {
	// TOD	O
	// Throw if overwriting existing
}

/*
GetCampaigns global or team.
*/
func GetCampaigns() {
	// TODO: Load a list of global or team-specific campaigns. Check perms for the latter.
}

/*
CreateCampaign TODO.
*/
func CreateCampaign() {
	// TODO: Clone global campaigns to team-specific campaign if requested. Only sysadmin can create global campaigns.
	id := uuid.New()
	log.Debugf("Created campaign %s", id)
}

/*
WorkCampaign TODO.
*/
func WorkCampaign() {
	// TODO: Set a team to work on a campaign. Check perms.
}

/*
CreateCompetitor TODO
*/
func CreateCompetitor() {
	// TODO: Create a competitor in campaigns.
}
