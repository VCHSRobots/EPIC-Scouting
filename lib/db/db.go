/*Package db provides tools for creating and interacting with the SQLite databases.*/
package db

import (
	"EPIC-Scouting/lib/lumberjack"
	"errors"
	"fmt"
	"time"

	"database/sql"
	"os"

	"github.com/google/uuid"
	"github.com/raja/argon2pw"

	// Comment to make golint happy.
	_ "github.com/mattn/go-sqlite3"
)

/*
VARIABLES
*/

var log = lumberjack.New("DB")

/*
DatabasePath is the path to the directory which holds the databases.
*/
var DatabasePath string

var dbUsers *sql.DB
var dbTeams *sql.DB
var dbCampaigns *sql.DB

/*
Schedule describes the current Campaign / Event / Match a team is contributing to.
*/
type Schedule struct {
	matchID string
}

/*
TeamData describes the elements which make up a team.
*/
type TeamData struct {
	TeamID      string
	TeamName    string
	TeamMembers map[string]string
	Schedule    string
}

/*
UserData describes all of the elements which describe a user of the scouting system.
*/
type UserData struct {
	UserID    string
	UserName  string
	Password  string
	FirstName string
	LastName  string
	Email     string // TODO: Add multiple contact options.
	SysAdmin  bool   // This is the only variable here which is NOT stored in the users/users table -- it comes from the users/sysadmin table
	LastSeen  string
}

/*
GENERAL FUNCTIONS
*/

/*
accessCheck determines whether a database was read or written to properly. If not, it reports the error via log.Fatalf
*/
func accessCheck(err error) {
	if err != nil {
		log.Fatalf("Unable to access database: %s", err.Error())
	}
}

/*
encryptPassword encrypts the given plaintext password with the given salt, and then escapes characters that might cause SQL some issue. Returns a string. Returns error if this fails for some reason.
*/
func encryptPassword(password string) (string, error) {
	hashedPassword, err := argon2pw.GenerateSaltedHash(password)
	if err != nil {
		log.Error("Password hash failed.")
		return "", errors.New("unable to hash password")
	}
	return hashedPassword, nil
}

/*
NullifyString makes empty strings into sql.NullStrings, and returns the original string if it isn't empty.
*/
func NullifyString(s string) interface{} {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return s
}

/*
TouchBase creates all databases used by the server if they do not exist, along with the default SysAdmin team and user.
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

	// Users.
	dbUsers = newDatabase("users")
	dbUsers.Exec("CREATE TABLE IF NOT EXISTS users ( userid TEXT PRIMARY KEY UNIQUE NOT NULL, username TEXT NOT NULL UNIQUE, password TEXT NOT NULL, firstname TEXT, lastname TEXT, email TEXT, lastseen TEXT )") // TODO: Add support for N+ contact options; via linked table?
	dbUsers.Exec("CREATE TABLE IF NOT EXISTS sysadmins ( userid TEXT PRIMARY KEY UNIQUE NOT NULL )")                                                                                                              // List of users which are SysAdmins.

	// Create a default SysAdmin user if it does not exist.
	_, errQuery := UserQuery("00000000-0000-0000-0000-000000000000") // TODO: handle the error UserCreate returns here.
	if errQuery != nil {
		result, _ := UserCreate(&UserData{UserID: "00000000-0000-0000-0000-000000000000", UserName: "SysAdmin", Password: "root", FirstName: "", LastName: "", Email: "", SysAdmin: true, LastSeen: ""})
		if result {
			log.Warn("Default system administrator account created. Username: \"SysAdmin\". Password: \"root\".")
			log.Warn("IMPORTANT: Change the password for this account before making your production server public!!!")
		}
		SysAdminPromote("00000000-0000-0000-0000-000000000000")
	}

	// Scouting teams.
	dbTeams = newDatabase("teams")
	dbTeams.Exec("CREATE TABLE IF NOT EXISTS teams ( teamid TEXT PRIMARY KEY UNIQUE NOT NULL, number TEXT UNIQUE, name TEXT NOT NULL, schedule TEXT NOT NULL )")                                                                                     // A team.
	dbTeams.Exec("CREATE TABLE IF NOT EXISTS members ( teamid TEXT PRIMARY KEY NOT NULL, userid TEXT NOT NULL, usertype TEXT NOT NULL )")                                                                                                            // The members on a team. UserType is either member or admin.
	dbTeams.Exec("CREATE TABLE IF NOT EXISTS participating ( teamid TEXT PRIMARY KEY NOT NULL, eventid TEXT NOT NULL, schedule TEXT )")                                                                                                              // What events a team is participating in. If a team is currently running a campaign, they must have *some* event they are participating in. A team is scouting all matches during an event, of course.
	dbTeams.Exec("CREATE TABLE IF NOT EXISTS results ( campaignid TEXT PRIMARY KEY NOT NULL, eventid TEXT NOT NULL, matchid TEXT NOT NULL, competitorid TEXT NOT NULL, teamid TEXT NOT NULL, userid TEXT NOT NULL, datetime TEXT NOT NULL, stats )") // A team's scouted results. Any number of teams may scout for the same campaign / event / match at the same time.

	// Create a default SysAdmin team if it does not exist.

	// Campaigns.
	dbCampaigns = newDatabase("campaigns")
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS campaigns ( campaignid TEXT PRIMARY KEY UNIQUE NOT NULL, owner TEXT NOT NULL, name TEXT NOT NULL )")                                      // TODO: Add more information about each campaign. Campaign owner is a teamid. If campaign owner is all zeros, campaign is global.
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS events ( campaignid TEXT PRIMARY KEY NOT NULL, eventid TEXT NOT NULL, name TEXT NOT NULL, location TEXT, starttime TEXT, endtime TEXT )") // TODO: Add more information about each event.
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS matches ( eventid TEXT PRIMARY KEY NOT NULL, matchid TEXT UNIQUE NOT NULL, matchnumber INTEGER NOT NULL, starttime TEXT, endtime TEXT )") // TODO: Add more information about each match.
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS participants ( matchid TEXT PRIMARY KEY NOT NULL, competitorid TEXT UNIQUE NOT NULL) ")                                                   // The participants in each match.
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS competitors ( competitorid TEXT PRIMARY KEY UNIQUE NOT NULL, number TEXT UNIQUE, name TEXT NOT NULL )")                                   // TODO: Add more information about each competing team.
}

/*
TEAM ADMIN FUNCTIONS
*/

// TODO: Teams need a way to weight each of their competitors when assigning teams via the scheduler.

/*
TeamCreate creates a new team from a TeamData struct. Returns bool false and an error if unable to create team.
*/
func TeamCreate() {

}

/*
TeamDelete deletes a team from the database. Returns false and an error if unable to delete.
Note that information recorded by the into a shared campaign is not deleted, but their account and private campaigns are deactivated (password is set to null.)
*/
func TeamDelete(teamID string) {
	// TODO.
}

/*
TeamList returns the teamID for every team in the system.
*/
func TeamList() (teams []string) {
	rows, err := dbTeams.Query("SELECT teamid FROM teams")
	accessCheck(err)
	defer rows.Close()
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		teams = append(teams, id)
		accessCheck(err)
	}
	return teams
}

/*
TeamQuery returns a team's information as a TeamData struct.
*/
func TeamQuery() {

}

/*
USER FUNCTIONS
*/

/*
UserCreate creates a new user from a UserData struct. Returns bool false and an error if unable to create user.
UserData.password is hashed.
If UserData.userID is null, a new and random ID is assigned.
UserData.lastSeen is set to the current time.
*/
func UserCreate(d *UserData) (bool, error) {
	data, err := UserQuery(d.UserID) // Check if user ID exists.
	if data != nil {
		err = errors.New("UserID already exists")
		log.Warnf("Unable to create user %q [%s]: %s", d.UserName, d.UserID, err.Error())
		return false, err
	}
	if err != nil {
		if err != sql.ErrNoRows {
			log.Errorf("Unable to create user %q [%s]: %s", d.UserName, d.UserID, err.Error())
			return false, err
		}
	}
	if d.UserID == "" {
		d.UserID = uuid.New().String()
	}
	hash, _ := encryptPassword(d.Password) // TODO: Handle error
	d.Password = hash
	d.LastSeen = time.Now().Format("2006-01-02 15:04:05")
	result, errExec := dbUsers.Exec(fmt.Sprintf("INSERT INTO users VALUES ( '%s', '%s', '%s', '%v', '%v', '%v', '%s' )", d.UserID, d.UserName, d.Password, NullifyString(d.FirstName), NullifyString(d.LastName), NullifyString(d.Email), d.LastSeen))
	if errExec != nil {
		if errExec != sql.ErrNoRows {
			log.Errorf("Unable to create user %q [%s]: %s", d.UserName, d.UserID, errExec.Error())
			return false, errExec
		}
		log.Warnf("Unable to create user %q [%s]: %s", d.UserName, d.UserID, errExec.Error())
		return false, errExec
	}
	log.Debugf("Created user %s: %q", d.UserID, d.UserName)
	log.Debugf("Result: %q", result) // TODO: TEMP
	return true, nil
}

/*
UserDelete deletes a user's login information. Returns false and an error if unable to delete user.
Note that information recorded by the user into a team's results is not deleted, but their account is deactivated (password is set to null.)
*/
func UserDelete(userID string) {
	// TODO
}

/*
UserLogin returns true if the username and password exist in the users database. Otherwise returns false with error.
*/
func UserLogin(username, password string) (loggedIn bool, err error) {
	var storedHash string
	err = dbUsers.QueryRow(fmt.Sprintf("SELECT username, password FROM users WHERE username='%s'", username)).Scan(&username, &storedHash)
	if err != nil {
		log.Debugf("Failed to log in user %q: %s", username, err.Error())
		loggedIn = false
		return
	}
	valid, err := argon2pw.CompareHashWithPassword(storedHash, password)
	if !valid {
		log.Debugf("Failed to log in user %q: %s", username, "password mismatch.")
		loggedIn = false
		return
	}
	log.Debugf("Logged in user %q.", username)
	loggedIn = true
	return
}

/*
UserModify modifies an existing user account. Returns an error if the user could not be found.
*/
func UserModify(userID string, data interface{}) {

}

/*
UserQuery returns the user's information as a UserData struct. Returns an error if the user could not be found.
If bool userName is true searches with userName instead of userID.
UserData.password is the hashed password.
*/
func UserQuery(userID string, userName bool) (*UserData, error) {
	var d UserData // TODO: Tidy this function up.
	if userName {
		errQueryRowUsers := dbUsers.QueryRow(fmt.Sprintf("SELECT * FROM users WHERE username='%s'", userID)).Scan(&d.UserID, &d.UserName, &d.Password, &d.FirstName, &d.LastName, &d.Email, &d.LastSeen) // Load user data.
		if errQueryRowUsers != nil {
			return nil, errQueryRowUsers
		}
	} else {
		errQueryRowUsers := dbUsers.QueryRow(fmt.Sprintf("SELECT * FROM users WHERE userid='%s'", userID)).Scan(&d.UserID, &d.UserName, &d.Password, &d.FirstName, &d.LastName, &d.Email, &d.LastSeen) // Load user data.
		if errQueryRowUsers != nil {
			return nil, errQueryRowUsers
		}
	}
	var foundID string
	errQueryRowSysAdmins := dbUsers.QueryRow(fmt.Sprintf("SELECT userid FROM sysadmins WHERE userid='%s'", userID)).Scan(&foundID) // Check if user is in the SysAdmin list.
	if errQueryRowSysAdmins == sql.ErrNoRows {
		d.SysAdmin = false
	} else if foundID == userID {
		d.SysAdmin = true
	}
	return &d, nil
}

/*
UserQueryTeams returns the list of teams a user is a member of and the user's userType for each team.
*/
func UserQueryTeams(userID string) map[string]string {

}

/*
CreateCampaign TODO.
*/
/*
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
*/

/*
CreateEvent adds an event to the event table in the campaigns database
*/
/*
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
*/

/*
CreateMatch adds a mach to the match table in the campaign database
*/
/*
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
		schedule
		log.Fatal("Unable to open or create database: " + err.Error())
	}
	log.Debugf("Created match: '%s'", id.String())
}
*/

/*
TEAM MATCH DATA READ-WRITE FUNCTIONS
*/

func getCampaignOwner(campaignid string, db *sql.DB) string {
	// TODO
	return ""
}

func retconCompetitorID(teamnumber, teamid, DatabasePath string) {
	// TODO
}

func getCompetitorNumber(competitorid string, db *sql.DB) string {
	// TODO
	return ""
}

/*
WriteResults TODO
*/
func WriteResults() {
	// TODO
	// Throw if overwriting existing
}

/*
SYSADMIN UTILITY FUNCTIONS
*/

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

/*
SysAdminPromote adds a user to the list of SysAdmins. Returns false if user is currently a SysAdmin. Returns error if an error occurs.
*/
func SysAdminPromote(userID string) bool {
	d, err := UserQuery(userID)
	if err != nil { // Error with query.
		log.Errorf("Unable to promote user %s: %s", userID, err.Error())
	}
	if d.SysAdmin { // User is already a SysAdmin.
		log.Infof("Unable to promote user %s: user is already a SysAdmin.", userID)
		return false
	}
	_, errAccess := dbUsers.Exec(fmt.Sprintf("INSERT INTO sysadmins VALUES ('%s')", userID))
	if errAccess != nil {
		log.Errorf("Unable to promote user %s: %s", userID, errAccess.Error())
		return false
	}
	log.Warnf("Promoted user %s to SysAdmin.", userID)
	return true
}

/*
SysAdminDemote removes a user from the list of SysAdmins. Returns false if the user is the only user in the list. There must be one!
*/
func SysAdminDemote(useriD string) bool {
	// TODO
	return false
}
