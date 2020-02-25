/*Package db provides tools for creating and interacting with the SQLite databases.*/
package db

import (
	"EPIC-Scouting/lib/lumberjack"
	"errors"
	"fmt"
	"strconv"
	"strings"
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
UserDataContact describes all of the elements of a user's contact information. TODO.
*/
/*
type UserDataContact struct {
	email string
	phone string
	other string
}
*/

/*
MatchData stores match data for transit to and from database
*/
type MatchData struct {
	MatchID          string
	MatchNum         int
	Team             int
	AutoLineCross    bool
	AutoLowBalls     int
	AutoHighBalls    int
	AutoBackBalls    int
	AutoPickups      int
	ShotQuantity     int
	LowFuel          int
	HighFuel         int
	BackFuel         int
	StageOneComplete bool
	StageOneTime     int
	StageTwoComplete bool
	StageTwoTime     int
	Fouls            int
	TechFouls        int
	Card             string
	Climbed          string
	Balanced         bool
	ClimbTime        int
	Comments         string
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
	dbTeams.Exec("CREATE TABLE IF NOT EXISTS teams ( teamid TEXT PRIMARY KEY UNIQUE NOT NULL, number TEXT UNIQUE, name TEXT NOT NULL, schedule TEXT NOT NULL )")                                                                              // A team.
	dbTeams.Exec("CREATE TABLE IF NOT EXISTS members ( userid TEXT, teamid TEXT NOT NULL, usertype TEXT NOT NULL )")                                                                                                                          // The members on a team. UserType is either member or admin.
	dbTeams.Exec("CREATE TABLE IF NOT EXISTS requestMembers ( userid TEXT, teamid TEXT NOT NULL )")                                                                                                                                           // Membership requests for teams
	dbTeams.Exec("CREATE TABLE IF NOT EXISTS participating ( teamid TEXT PRIMARY KEY NOT NULL, eventid TEXT NOT NULL, schedule TEXT )")                                                                                                       // What events a team is participating in. If a team is currently running a campaign, they must have *some* event they are participating in. A team is scouting all matches during an event, of course.
	dbTeams.Exec("CREATE TABLE IF NOT EXISTS results ( campaignid TEXT PRIMARY KEY NOT NULL, eventid TEXT NOT NULL, matchid TEXT NOT NULL, competitorid TEXT NOT NULL, teamid TEXT NOT NULL, userid TEXT NOT NULL, datetime TEXT NOT NULL )") // A team's scouted results. Any number of teams may scout for the same campaign / event / match at the same time.

	// Create a default SysAdmin team if it does not exist.

	// Campaigns.
	dbCampaigns = newDatabase("campaigns")
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS campaigns ( campaignid TEXT PRIMARY KEY UNIQUE NOT NULL, owner TEXT NOT NULL, name TEXT NOT NULL )")                                            // TODO: Add more information about each campaign. Campaign owner is a teamid. If campaign owner is all zeros, campaign is global.
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS events ( eventid TEXT PRIMARY KEY NOT NULL, campaignid TEXT NOT NULL, name TEXT NOT NULL, location TEXT, starttime INTEGER, endtime INTEGER )") // TODO: Add more information about each event.
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS matches ( eventid TEXT PRIMARY KEY NOT NULL, matchid TEXT UNIQUE NOT NULL, matchnumber INTEGER NOT NULL, active BIT )")                         // TODO: Add more information about each match.
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS matchscout ( scoutid TEXT PRIMARY KEY, matchid TEXT NOT NULL, matchnumber INTEGER NOT NULL, userid TEXT NOT NULL, team INTEGER, autoLineCross BIT, autoLowBalls INTEGER, autoHighBalls INTEGER, autoBackBalls INTEGER, autoPickups INTEGER, shotQuantity INTEGER, lowFuel INTEGER, highFuel INTEGER, backFuel INTEGER, stageOneComplete BIT, stageOneTime INTEGER, stageTwoComplete BIT, stageTwoTime INTEGER, fouls INTEGER, techFouls INTEGER, card TEXT, climbed TEXT, balanced BIT, climbtime INTEGER, comments TEXT )")
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS pitscout ( eventid TEXT PRIMARY KEY NOT NULL, matchid TEXT UNIQUE NOT NULL, matchnumber INTEGER NOT NULL, starttime TEXT, endtime TEXT )")
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS participants ( matchid TEXT PRIMARY KEY NOT NULL, competitorid TEXT UNIQUE NOT NULL) ")                 // The participants in each match.
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS competitors ( competitorid TEXT PRIMARY KEY UNIQUE NOT NULL, number TEXT UNIQUE, name TEXT NOT NULL )") // TODO: Add more information about each competing team.
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
UserData.password is the hashed password.
*/
func UserQuery(userID string) (*UserData, error) {
	var d UserData
	errQueryRowUsers := dbUsers.QueryRow(fmt.Sprintf("SELECT * FROM users WHERE username='%s'", userID)).Scan(&d.UserID, &d.UserName, &d.Password, &d.FirstName, &d.LastName, &d.Email, &d.LastSeen) // Load user data.
	if errQueryRowUsers == sql.ErrNoRows {
		errQueryRowUsers = dbUsers.QueryRow(fmt.Sprintf("SELECT * FROM users WHERE userid='%s'", userID)).Scan(&d.UserID, &d.UserName, &d.Password, &d.FirstName, &d.LastName, &d.Email, &d.LastSeen) // Load user data.
	}
	if errQueryRowUsers != nil {
		return nil, errQueryRowUsers
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
GetTeamCampaign gets the uuid of the campaign with which a team is associated
*/
func GetTeamCampaign(teamID string) (string, error) {
	var campaign string
	//TODO account for schedule as object instead of storing active campaign in it directly
	err := dbTeams.QueryRow(fmt.Sprintf("SELECT schedule FROM teams WHERE teamid='%s'", teamID)).Scan(&campaign)
	if err != nil {
		return "", err
	}
	return campaign, nil
}

/*
GetTeamEvent gets the event in which a team is currently participating
*/
func GetTeamEvent(teamID string) (string, error) {
	//test
	if teamID == "4415epicrobotz" {
		return "epicevent", nil
	}
	campaignid, err := GetTeamCampaign(teamID)
	if err != nil {
		return "", err
	}
	event, err := getActiveCampaignEvent(campaignid)
	if err != nil {
		return "", err
	}
	return event, nil
}

/*
DATA STORAGE FUNCTIONS
*/

/*
StoreMatch takes the array of data from the form and stores it in the database
*/
func StoreMatch(arr []string, agentid, teamid string) error {
	eventid, err := GetTeamEvent(teamid)
	if err != nil {
		return err
	}
	data, err := arrToMatchStruct(arr, eventid, agentid)
	//don't log these errors here since the function calls should log them on their own
	if err != nil {
		return err
	}
	scoutid := uuid.New().String()
	fmt.Println(data.Team)
	result, errExec := dbCampaigns.Exec(fmt.Sprintf("INSERT INTO matchscout VALUES ( '%s', '%s', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%s', '%v', '%v', '%s' )", scoutid, data.MatchID, data.MatchNum, agentid, data.Team, data.AutoLineCross, data.AutoLowBalls, data.AutoHighBalls, data.AutoBackBalls, data.AutoPickups, data.ShotQuantity, data.LowFuel, data.HighFuel, data.BackFuel, data.StageOneComplete, data.StageOneTime, data.StageTwoComplete, data.StageTwoTime, data.Fouls, data.TechFouls, data.Card, data.Climbed, data.Balanced, data.ClimbTime, escapeText(data.Comments)))
	if errExec != nil {
		log.Errorf("Unable to write match scouting data to database: %s", errExec)
		return errExec
	}
	log.Debugf("Match Write Result: %v", result)
	return nil
}

/*
arrToMatchStruct turns the data array into a match struct
*/
func arrToMatchStruct(arr []string, eventid, agentid string) (*MatchData, error) {
	var autoLineCross, stageOneComplete, stageTwoComplete bool
	var card, climbed string
	//cuts off the last value since comments can't be converted into integers
	intarr, err := convertArrToInts(arr[:len(arr)-1])
	if err != nil {
		return nil, err
	}
	matchid, err := matchIDFromNum(intarr[0], eventid)
	if matchid == "" {
		return nil, err
	}
	autoLineCross = intarr[2] == 1
	stageOneComplete = intarr[11] == 1
	stageTwoComplete = intarr[13] == 1
	balanced := intarr[19] == 1
	if intarr[17] == 1 {
		card = "yellow"
	} else if intarr[17] == 2 {
		card = "red"
	} else {
		card = "none"
	}
	if intarr[18] == 2 {
		climbed = "climbed"
	} else if intarr[18] == 1 {
		climbed = "platform"
	} else {
		climbed = "none"
	}
	return &MatchData{MatchID: matchid, MatchNum: intarr[0], Team: intarr[1], AutoLineCross: autoLineCross, AutoLowBalls: intarr[3], AutoHighBalls: intarr[4], AutoBackBalls: intarr[5], AutoPickups: intarr[6], ShotQuantity: intarr[7], LowFuel: intarr[8], HighFuel: intarr[9], BackFuel: intarr[10], StageOneComplete: stageOneComplete, StageOneTime: intarr[12], StageTwoComplete: stageTwoComplete, StageTwoTime: intarr[14], Fouls: intarr[15], TechFouls: intarr[16], Card: card, Climbed: climbed, Balanced: balanced, ClimbTime: intarr[20], Comments: arr[21]}, nil
}

func convertArrToInts(arr []string) ([]int, error) {
	var intval int
	var err error
	ints := make([]int, len(arr))
	for ind, str := range arr {
		intval, err = strconv.Atoi(str)
		ints[ind] = intval
		if err != nil {
			log.Errorf("Unable to convert match data array to integers: %s", err)
			return ints, err
		}
	}
	return ints, nil
}

func escapeText(str string) string {
	return strings.ReplaceAll(str, "'", "{singlequote}")
}

func unescapeText(str string) string {
	return strings.ReplaceAll(str, "{singlequote}", "'")
}

/*
ReadMatch turns database outputs from a given match to match struct
*/
func ReadMatch(num int, eventid, scouterid string) (*MatchData, error) {
	matchid, err := matchIDFromNum(num, eventid)
	if err != nil {
		return nil, err
	}
	var d MatchData
	err = dbCampaigns.QueryRow("SELECT team, autoLineCross, autoLowBalls, autoHighBalls, autoBackBalls, autoPickups, shotQuantity, lowFuel, highFuel, backFuel, stageOneComplete, stageOneTime, stageTwoComplete, stageTwoTime, fouls, techFouls, card, climbed, climbtime, comments FROM matchscout WHERE matchid='%s' AND userid='%s'", matchid, scouterid).Scan(&d.Team, &d.AutoLineCross, &d.AutoLowBalls, &d.AutoHighBalls, &d.AutoBackBalls, &d.AutoPickups, &d.ShotQuantity, &d.LowFuel, &d.HighFuel, &d.BackFuel, &d.StageOneComplete, &d.StageTwoComplete, &d.StageTwoTime, &d.Fouls, &d.TechFouls, &d.Card, &d.Comments)
	d.MatchID = matchid
	d.MatchNum = num
	d.Comments = unescapeText(d.Comments)
	return &d, nil
}

func matchIDFromNum(num int, eventid string) (string, error) {
	var matchid string
	err := dbCampaigns.QueryRow(fmt.Sprintf("SELECT matchid FROM matches WHERE eventid='%s' AND matchnumber='%v'", eventid, num)).Scan(&matchid)
	if err != nil {
		log.Warnf("Failed to retrive match id for match #%v from event %s", num, eventid)
		return "", err
	}
	return matchid, nil
}

/*
GetMatchScoutData gets scouter's data based on a match id
*/
func GetMatchScoutData(matchid string) (*[]MatchData, error) {
	data := make([]MatchData, 0)
	rows, _ := dbCampaigns.Query(fmt.Sprintf("SELECT * FROM matchscout WHERE matchid='%s'", matchid))
	for rows.Next() {
		var d MatchData
		err := rows.Scan(nil, &d.MatchID, &d.MatchNum, nil, &d.Team, &d.AutoLineCross, &d.AutoLowBalls, &d.AutoHighBalls, &d.AutoBackBalls, &d.AutoPickups, &d.ShotQuantity, &d.LowFuel, &d.HighFuel, &d.StageOneComplete, &d.StageOneTime, &d.StageTwoComplete, &d.StageTwoTime, &d.Fouls, &d.TechFouls, &d.Card, &d.Comments)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}
	return &data, nil
}

/*
GetTeamScoutData gets scouter's data based on a team id
*/
func GetTeamScoutData(teamid string) (*[]MatchData, error) {
	data := make([]MatchData, 0)
	rows, _ := dbCampaigns.Query(fmt.Sprintf("SELECT matchid, matchnumber, team, autoLineCross, autoLowBalls, autoHighBalls, autoBackBalls, autoPickups, shotQuantity, lowFuel, highFuel, backFuel, stageOneComplete, stageOneTime, stageTwoComplete, stageTwoTime, fouls, techFouls, card, climbed, balanced, climbtime, comments FROM matchscout"))
	for rows.Next() {
		var d MatchData
		err := rows.Scan(&d.MatchID, &d.MatchNum, &d.Team, &d.AutoLineCross, &d.AutoLowBalls, &d.AutoHighBalls, &d.AutoBackBalls, &d.AutoPickups, &d.ShotQuantity, &d.LowFuel, &d.HighFuel, &d.BackFuel, &d.StageOneComplete, &d.StageOneTime, &d.StageTwoComplete, &d.StageTwoTime, &d.Fouls, &d.TechFouls, &d.Card, &d.Climbed, &d.Balanced, &d.ClimbTime, &d.Comments)
		if err != nil {
			return nil, err
		}
		data = append(data, d)
	}
	return &data, nil
}

/*
TEAM ADMIN FUNCTIONS
*/

/*
CAMPAIGN FUCTIONS
*/

/*
getActiveCampaignEvent gets the eventid of the active event in the given campaign
*/
func getActiveCampaignEvent(campaignid string) (string, error) {
	var eventid string
	var starttime, endtime int64
	rows, err := dbCampaigns.Query(fmt.Sprintf("SELECT eventid, starttime, endtime FROM events WHERE campaignid='%s'", campaignid))
	//checks if event is currently taking place
	for rows.Next() {
		err = rows.Scan(&eventid, &starttime, &endtime)
		if err != nil {
			return "", err
		}
		if starttime < time.Now().Unix() && endtime > time.Now().Unix() {
			break
		}
	}
	return eventid, nil
}

/*
CreateCampaign TODO.
*/
func CreateCampaign(agentid, owner, name string) {
	// TODO: Clone global campaigns to team-specific campaign if requested.
	// Only sysadmin can create global campaigns.
	uuid := uuid.New().String()
	dbCampaigns.Exec(fmt.Sprintf("INSERT INTO campaigns VALUES ( '%s', '%s', '%s' )", uuid, owner, name))
}

/*
CreateEvent adds an event to the event table in the campaigns database
Its starttime and endtime should be Unix time integers of its start and end dates
*/
func CreateEvent(campaignid, agentid, name, location string, starttime, endtime int) error {
	eventid := uuid.New().String()
	dbCampaigns.Exec(fmt.Sprintf("INSERT INTO events VALUES ( '%s', '%s', '%s', '%s', '%v', '%v')", eventid, campaignid, name, location, starttime, endtime))
	return nil
}

/*
CreateMatch adds a mach to the match table in the campaign database
*/
func CreateMatch(eventid, agentid string, num int, active bool) error {
	matchid := uuid.New().String()
	dbCampaigns.Exec(fmt.Sprintf("INSERT INTO matches VALUES ( '%s', '%s', '%v', '%v' )", eventid, matchid, num, active))
	return nil
}

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
