/*Package db provides tools for creating and interacting with the SQLite databases.*/
package db

import (
	"EPIC-Scouting/lib/lumberjack"
	"crypto/rand"
	"fmt"

	"database/sql"
	"os"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

var log = lumberjack.New("DB")

/*
DatabasePath is the path to the directory which holds the databases.
*/
var DatabasePath string

/*
TouchBase creates all databases used by the server if they do not exist.
Its name is a play on the GNU program "touch", the idiom "[to] touch base", and the word "database". The author is rather proud of this.
*/
func TouchBase(databasePath string) {
	DatabasePath = databasePath
	newDatabase := func(databaseName string) *sql.DB {
		db, err := sql.Open("sqlite3", DatabasePath+databaseName+".db")
		accessCheck(err)
		return db
	}

	err := os.MkdirAll(DatabasePath, 0755)
	accessCheck(err)

	// Users
	// This database stores all users.
	users := newDatabase("users")
	users.Exec("CREATE TABLE IF NOT EXISTS users ( userid TEXT PRIMARY KEY UNIQUE NOT NULL, username TEXT NOT NULL UNIQUE, password TEXT NOT NULL, firstname TEXT, lastname TEXT, email TEXT UNIQUE, phone TEXT, usertype TEXT, salt TEXT)") // TODO: Add support for N+ contact options; via linked table?

	// TODO: Create a SYSTEM team which makes the default public campaigns each season.

	// Scouting teams
	teams := newDatabase("teams")
	teams.Exec("CREATE TABLE IF NOT EXISTS teams ( teamid TEXT PRIMARY KEY UNIQUE NOT NULL, number TEXT UNIQUE, name TEXT NOT NULL, schedule TEXT NOT NULL )")                                                                                     // A team.
	teams.Exec("CREATE TABLE IF NOT EXISTS members ( teamid TEXT PRIMARY KEY NOT NULL, userid TEXT NOT NULL, usertype TEXT NOT NULL )")                                                                                                            // The members on a team. UserType is either member or admin.
	teams.Exec("CREATE TABLE IF NOT EXISTS participating ( teamid TEXT PRIMARY KEY NOT NULL, eventid TEXT NOT NULL, schedule TEXT )")                                                                                                              // What events a team is participating in. If a team is currently running a campaign, they must have *some* event they are participating in. A team is scouting all matches during an event, of course.
	teams.Exec("CREATE TABLE IF NOT EXISTS results ( campaignid TEXT PRIMARY KEY NOT NULL, eventid TEXT NOT NULL, matchid TEXT NOT NULL, competitorid TEXT NOT NULL, teamid TEXT NOT NULL, userid TEXT NOT NULL, datetime TEXT NOT NULL, stats )") // A team's scouted results. Any number of teams may scout for the same campaign / event / match at the same time.

	// Campaigns (game seasons / years)
	// This database stores the expected campaign / event / match schedule and data, and expected competing teams. Data pulled from TBA.
	campaigns := newDatabase("campaigns")
	campaigns.Exec("CREATE TABLE IF NOT EXISTS campaigns ( campaignid TEXT PRIMARY KEY UNIQUE NOT NULL, owner TEXT NOT NULL, name TEXT NOT NULL )")                                      // TODO: Add more information about each campaign. Campaign owner is a teamid. If campaign owner is all zeros, campaign is global.
	campaigns.Exec("CREATE TABLE IF NOT EXISTS events ( campaignid TEXT PRIMARY KEY NOT NULL, eventid TEXT NOT NULL, name TEXT NOT NULL, location TEXT, starttime TEXT, endtime TEXT )") // TODO: Add more information about each event.
	campaigns.Exec("CREATE TABLE IF NOT EXISTS matches ( eventid TEXT PRIMARY KEY NOT NULL, matchid TEXT UNIQUE NOT NULL, matchnumber INTEGER NOT NULL, starttime TEXT, endtime TEXT )") // TODO: Add more information about each match.
	campaigns.Exec("CREATE TABLE IF NOT EXISTS participants ( matchid TEXT PRIMARY KEY NOT NULL, competitorid TEXT UNIQUE NOT NULL) ")                                                   // The participants in each match.
	campaigns.Exec("CREATE TABLE IF NOT EXISTS competitors ( competitorid TEXT PRIMARY KEY UNIQUE NOT NULL, number TEXT UNIQUE, name TEXT NOT NULL )")                                   // TODO: Add more information about each competing team.

	// TODO: Teams need a way to weight each of their competitors.
}

/*
CheckLogin checks if a user has a valid login and returns their data (sans their password) if they do.

func CheckLogin(username, password string) (bool, []string) {
	users, err := sql.Open("sqlite3", DatabasePath+"users.db")
	accessCheck(err)
	usernames := processQuery(users.Query("SELECT username FROM users"))
	passwords := processQuery(users.Query("SELECT password FROM users"))
	userids := processQuery(users.Query("SELECT userid FROM users"))
	//TODO: Make this check hashes for all people with the same username. This may or may not be added.
	checkhash := correlateFields(username, usernames, passwords)
	userid := correlateFields(username, usernames, userids)
	log.Debugf("Login Attempt: Username: %s Password: %s\n", username, password)
	//TODO: Hash password for check
	salted := GetUserSaltedPassword(password, userid, users)
	userinfo := make([]string, 0, 6)
	hashmatched := bcrypt.CompareHashAndPassword([]byte(checkhash), []byte(salted))
	if hashmatched == nil {
		userind := -1
		for ind, user := range usernames {
			if user == username {
				userind = ind
				break
			}
		}
		if userind == -1 {
			//this should only happen if something went terribly wrong - username was already checked at this point
			return false, userinfo
		}
		userids := processQuery(users.Query("SELECT userid FROM users"))
		firstnames := processQuery(users.Query("SELECT firstname FROM users"))
		lastnames := processQuery(users.Query("SELECT lastname FROM users"))
		emails := processQuery(users.Query("SELECT email FROM users"))
		phones := processQuery(users.Query("SELECT phone FROM users"))
		usertypes := processQuery(users.Query("SELECT usertype FROM users"))
		userinfo = append(userinfo, userids[userind], username, firstnames[userind], lastnames[userind], emails[userind], phones[userind], usertypes[userind])
		return true, userinfo
	}
	return false, userinfo
}
*/

/*
CheckLoginII is an attempt at re-creating CheckLogin with different SQL queries.
*/
func CheckLoginII(username, password string) (loggedIn bool) {
	users, err := sql.Open("sqlite3", DatabasePath+"users.db")
	accessCheck(err)
	var uname string
	err = users.QueryRow("SELECT PASSWORD FROM USERS WHERE USERNAME = ?", username).Scan(&uname)
	accessCheck(err)
	log.Debugf("Checking login for username %q", uname)
	return loggedIn
}

/*
CreateUser creates a new user.
*/
func CreateUser(databasePath, username, password, firstname, lastname, email, phone, usertype string) {
	fmt.Printf("Creating user %s\n", username)
	users, err := sql.Open("sqlite3", databasePath+"users.db")
	if err != nil {
		log.Fatal("Unable to open or create database: " + err.Error())
		return
	}
	id := uuid.New()
	//make hash for password with added salt of len 256
	//salt comes out as a printed array of integers because sql doesn't like the random bytes converted to strings. it looks odd when you print it out but it works.
	hash, salt := hashNewPassword(password, fmt.Sprint(id))
	//TODO: Add password hashing with uuid as salt
	_, err = users.Exec(fmt.Sprintf("INSERT INTO users (userid, username, password, firstname, lastname, email, phone, usertype, salt) VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s');", fmt.Sprint(id), username, string(hash), firstname, lastname, email, phone, usertype, salt))
	// Throw if doing something illegal, such as overwriting existing user
	if err != nil {
		log.Info("Unable to create user: " + err.Error())
	}
	log.Debugf("Created user: '%s'", id.String())
}

/*
CreateCampaign TODO.
*/
func CreateCampaign(DatabasePath, agentid, owner, name string) {
	// TODO: Clone global campaigns to team-specific campaign if requested.
	// Only sysadmin can create global campaigns.
	campaigns, err := sql.Open("sqlite3", DatabasePath+"campaigns.db")
	if err != nil {
		log.Fatal("Unable to open or create database: " + err.Error())
		return
	}
	//Check if user is authorized to create campaigns
	usersDb, err := sql.Open("sqlite3", DatabasePath+"users.db")
	if err != nil {
		log.Fatal("Unable to open or create database: " + err.Error())
		return
	}
	usertype := checkUserType(agentid, usersDb)
	if (owner == "globalowner" && usertype != "sysadmin") || (usertype == "invalid") {
		log.Info("Unauthorized attempt to create global database. User ID: " + agentid)
		return
	}
	id := uuid.New()
	_, err = campaigns.Exec(fmt.Sprintf("INSERT INTO campaigns (campaignid, owner, name) VALUES ('%s', '%s', '%s');", id, owner, name))
	// Throw if doing something illegal, such as overwriting existing
	if err != nil {
		log.Fatal("Unable to open or create database: " + err.Error())
	}
	log.Debugf("Created campaign '%s'", id)
}

//CreateMatch adds a mach to the match table in the campaign database
func CreateMatch(DatabasePath, agentid, eventid, matchnumber, starttime, endtime string) {
	//TODO: Adds a match to a campaign
	//TODO: Adds event to campaign table
	campaigns, err := sql.Open("sqlite3", DatabasePath+"campaigns.db")
	if err != nil {
		log.Fatal("Unable to open or create database: " + err.Error())
		return
	}
	//Check to see the team who owns the match's campaign tied to the event the match is connected to. If the campaign is not write accessable to the agent, deny access
	usersDb, err := sql.Open("sqlite3", DatabasePath+"users.db")
	if err != nil {
		log.Fatal("Unable to open or create database: " + err.Error())
		return
	}
	usertype := checkUserType(agentid, usersDb)
	userteam := checkUserTeam(agentid, usersDb)
	campaignid := getEventCampaign(eventid, campaigns)
	campaignowner := getCampaignOwner(campaignid, campaigns)
	//TODO: get real globalowner name here
	if (campaignowner != "globalowner" && campaignowner != userteam) || (usertype == "invalid") {
		log.Info(fmt.Sprintf("Unauthorized attempt to create match in event ID %s in campaign %s. User ID: %s", eventid, campaignid, agentid))
		return
	}
	id := uuid.New()
	_, err = campaigns.Exec(fmt.Sprintf("INSERT INTO matches (matchid, eventid, matchnumber, starttime, endtime) VALUES ('%s', '%s', '%s', '%s', '%s');", id, eventid, matchnumber, starttime, endtime))
	// Throw if doing something illegal, such as overwriting existing
	if err != nil {
		log.Fatal("Unable to open or create database: " + err.Error())
	}
	log.Debugf("Created match: '%s'", id.String())
}

//CreateEvent adds an event to the event table in the campaigns database
func CreateEvent(DatabasePath, agentid, campaignid, name, location, starttime, endtime string) {
	//TODO: Adds event to campaign table
	campaigns, err := sql.Open("sqlite3", DatabasePath+"campaigns.db")
	if err != nil {
		log.Fatal("Unable to open or create database: " + err.Error())
		return
	}
	//Check to see if agent has write access to the campaign. If not, deny access
	usersDb, err := sql.Open("sqlite3", DatabasePath+"users.db")
	if err != nil {
		log.Fatal("Unable to open or create database: " + err.Error())
		return
	}
	usertype := checkUserType(agentid, usersDb)
	userteam := checkUserTeam(agentid, usersDb)
	campaignowner := getCampaignOwner(campaignid, campaigns)
	//TODO: get real globalowner name here
	if (campaignowner != "globalowner" && campaignowner != userteam) || (usertype == "invalid") {
		log.Info(fmt.Sprintf("Unauthorized attempt to add event to campaign ID %s. User ID: %s", campaignid, agentid))
		return
	}
	id := uuid.New()
	_, err = campaigns.Exec(fmt.Sprintf("INSERT INTO events (eventid, campaignid, name, location, starttime, endtime) VALUES ('%s', '%s', '%s', '%s', '%s', '%s');", id, campaignid, name, location, starttime, endtime))
	// Throw if doing something illegal, such as overwriting existing
	if err != nil {
		log.Fatal("Unable to open or create database: " + err.Error())
	}
	log.Debugf("Created event: '%s'", id.String())
}

/*
CreateTeam TODO

func CreateTeam(DatabasePath, agentid, number, name, currentcampaign string) {
	// TODO
	// Req: TeamNumber, TeamName, etc
	//TODO: Adds event to campaign table
	teams, err := sql.Open("sqlite3", DatabasePath+"teams.db")
	if err != nil {
		log.Fatal("Unable to open or create database: " + err.Error())
		return
	}
	//Check to see if agent is a valid user
	usersDb, err := sql.Open("sqlite3", DatabasePath+"users.db")
	if err != nil {
		log.Fatal("Unable to open or create database: " + err.Error())
		return
	}
	usertype := checkUserType(agentid, usersDb)
	if usertype == "invalid" {
		log.Info(fmt.Sprintf("Unauthorized user %s attempted to create team %s: %s", agentid, number, name))
		return
	}
	id := uuid.New()
	_, err = teams.Exec(fmt.Sprintf("INSERT INTO teams (teamid, number, name, currentcampaign) VALUES ('%s', '%s', '%s', '%s');", id, number, name, currentcampaign))
	// Throw if doing something illegal, such as overwriting existing
	if err != nil {
		log.Fatal("Unable to open or create database: " + err.Error())
	}
	log.Debugf("Created team: %s", id.String())
	// If a team's number exists in campaigns/competitors, their competitorid is their new teamid. If a team's number exists in teams/teams, any reference to their competitorid in campaigns/competitors is their teamid.
	retconCompetitorID(number, fmt.Sprint(id), DatabasePath)
}

/*
CreateCompetitor TODO
mplate "footer"}}ccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc
	if err != nil {
		log.Fatal("Unable to open or create database: " + err.Error())
		return
	}
	//Check to see if agent is a valid user
	usersDb, err := sql.Open("sqlite3", DatabasePath+"users.db")
	if err != nil {
		log.Fatal("Unable to open or create database: " + err.Error())
		return
	}
	usertype := checkUserType(agentid, usersDb)
	if usertype == "invalid" {
		log.Info(fmt.Sprintf("Unauthorized user %s attempted to create competitor %s: %s", agentid, number, name))
		return
	}
	id := uuid.New()
	_, err = campaigns.Exec(fmt.Sprintf("INSERT INTO competitors (competitorid, number, name) VALUES ('%s', '%s', '%s');", id, number, name))
	// Throw if doing something illegal, such as overwriting existing
	if err != nil {
		log.Fatal("Unable to open or create database: " + err.Error())
	}
	log.Debugf("Created event: %s", id.String())
}*/

func getCampaignOwner(campaignid string, db *sql.DB) string {
	campaignids := processQuery(db.Query("SELECT campaignid FROM campaigns;"))
	owners := processQuery(db.Query("SELECT owner FROM campaigns;"))
	return correlateFields(campaignid, campaignids, owners)
}

func processQuery(rows *sql.Rows, err error) []string {
	ind := 0
	outputs := make([]string, 0)
	if err != nil {
		return outputs
	}
	for rows.Next() {
		var str string
		rows.Scan(&str)
		outputs = append(outputs, str)
		ind++
	}
	return outputs
}

func checkUserType(userid string, db *sql.DB) string {
	userids := processQuery(db.Query("SELECT userid FROM users;"))
	usertypes := processQuery(db.Query("SELECT usertype FROM users;"))
	return correlateFields(userid, userids, usertypes)
}

func checkUserTeam(userid string, db *sql.DB) string {
	userids := processQuery(db.Query("SELECT userid FROM users;"))
	userteams := processQuery(db.Query("SELECT team FROM users;"))
	return correlateFields(userid, userids, userteams)
}

func getEventCampaign(eventid string, db *sql.DB) string {
	eventids := processQuery(db.Query("SELECT eventid FROM events;"))
	campaignids := processQuery(db.Query("SELECT campaignid FROM events;"))
	return correlateFields(eventid, eventids, campaignids)
}

func retconCompetitorID(teamnumber, teamid, DatabasePath string) {
	//Replaces all references to a competitor team with their proper team id that matches with the teams table
	campaigns, err := sql.Open("sqlite3", DatabasePath+"campaigns.db")
	if err != nil {
		log.Fatal("Unable to open or create database: " + err.Error())
		return
	}
	competitorids := processQuery(campaigns.Query("SELECT competitorid FROM competitors;"))
	numbers := processQuery(campaigns.Query("SELECT number FROM competitors;"))
	oldid := correlateFields(teamnumber, numbers, competitorids)
	//TODO: Determine places to swap out old id for new one
	fmt.Sprint(oldid)
}

func getCompetitorNumber(competitorid string, db *sql.DB) string {
	competitorids := processQuery(db.Query("SELECT competitorid FROM competitors;"))
	numbers := processQuery(db.Query("SELECT number FROM competitors;"))
	return correlateFields(competitorid, competitorids, numbers)
}

func correlateFields(term string, searchfield, resultfield []string) string {
	searchind := -1
	for ind, search := range searchfield {
		if search == term {
			searchind = ind
			break
		}
	}
	if searchind == -1 {
		//If the search did not match anything in the searchfield
		return "invalid"
	}
	return resultfield[searchind]
}

func hashNewPassword(password, id string) ([]byte, string) {
	//always hashes passwords with a 256 byte salt
	salt := [256]byte{}
	_, err := rand.Read(salt[:])
	if err != nil {
		log.Fatal("Unable to generate salt for user password")
	}
	saltstring := fmt.Sprint(salt)
	pstring := fmt.Sprintf("%s_%s%s", password, fmt.Sprint(id), saltstring)
	fmt.Printf("New salted password: %s\n", pstring)
	hash, err := bcrypt.GenerateFromPassword([]byte(pstring), 10)
	if err != nil {
		log.Info("Password hash failed")
	}
	return hash, saltstring
}

//GetUserSaltedPassword gets the salted password string for the given userid
func GetUserSaltedPassword(password, id string, userDb *sql.DB) string {
	userids := processQuery(userDb.Query("SELECT userid FROM users"))
	salts := processQuery(userDb.Query("SELECT salt FROM users"))
	salt := correlateFields(id, userids, salts)
	pstring := fmt.Sprintf("%s_%s%s", password, id, salt)
	fmt.Printf("Got salted password: %s\n", pstring)
	return pstring
}

//GetUserPasswordHash gets the password hash for a specific user's password
func GetUserPasswordHash(id string) []byte {
	userDb, err := sql.Open("sqlite3", DatabasePath+"users.db")
	if err != nil {
		log.Fatal(fmt.Sprintf("Unable to create or open database: %s", DatabasePath+"users.db"))
	}
	userids := processQuery(userDb.Query("SELECT userid FROM users"))
	passwords := processQuery(userDb.Query("SELECT password FROM users"))
	hash := []byte(correlateFields(id, userids, passwords))
	return hash
}

/*
WriteResults TODO
*/
func WriteResults() {
	// TODO
	// Throw if overwriting existing
}

/*
GetCampaigns global or team.

func GetCampaigns(DatabasePath, agentid, campaignid, teamid string) {
	// TODO: Load a list of global or team-specific campaigns. Check perms for the latter.
}

/*
WorkCampaign TODO.

func WorkCampaign(DatabasePath, agentid, teamid, campaignid string) {
	// TODO: Set a team to work on a campaign. Check perms.
}

/*
accessCheck determines whether a database was read or written to properly. If not, it reports the error via log.Fatalf
*/
func accessCheck(err error) {
	if err != nil {
		log.Fatalf("Unable to access database: %s", err.Error())
	}
}

/*
GetDatabaseSize returns the databases' current sizes in bytes as a []string.
*/
func GetDatabaseSize() map[string]int64 {
	results := make(map[string]int64)
	bases := []string{"users", "teams", "campaigns"}
	for _, base := range bases {
		file, error := os.Stat(DatabasePath + base + ".db")
		accessCheck(error)
		size := file.Size()
		results[base] = size
	}
	return results
}
