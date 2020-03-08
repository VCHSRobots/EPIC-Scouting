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
TeamData describes the most of the data regarding a team.
*/
type TeamData struct {
	TeamID             string
	TeamName           string
	TeamMembers        map[string]string // UserID and UserType.
	AvaliableCampaigns map[string]bool   // List of CampaignIDs a team may write to. Bool indicates if team has write access, FALSE = read only.
	Schedule           []string          // CampaignID, EventID, and MatchID for the team's current scouting.
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
	Alliance         string
	AutoLineCross    bool
	AutoLowBalls     int
	AutoHighBalls    int
	AutoBackBalls    int
	AutoShots        int
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
accessCheck determines whether a database was read or written to properly. If not, it reports the error via log.Fatalf().
*/
func accessCheck(err error) {
	if err != nil {
		log.Fatalf("Unable to access database: %s", err.Error())
	}
}

/*
encryptPassword encrypts the given plaintext. Returns a string. Returns error if this fails for some reason.
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
	dbTeams.Exec("CREATE TABLE IF NOT EXISTS teams ( teamid TEXT PRIMARY KEY UNIQUE NOT NULL, number TEXT UNIQUE, name TEXT NOT NULL, schedule TEXT NOT NULL )") // A team.
	dbTeams.Exec("CREATE TABLE IF NOT EXISTS members ( userid TEXT, teamid TEXT NOT NULL, usertype TEXT NOT NULL )")                                             // The members on a team. UserType is either member or admin.
	dbTeams.Exec("CREATE TABLE IF NOT EXISTS requestMembers ( userid TEXT, teamid TEXT NOT NULL )")                                                              // Membership requests for teams
	dbTeams.Exec("CREATE TABLE IF NOT EXISTS participating ( teamid TEXT PRIMARY KEY NOT NULL, eventid TEXT NOT NULL, schedule TEXT )")                          // What events a team is participating in. If a team is currently running a campaign, they must have *some* event they are participating in. A team is scouting all matches during an event, of course.

	dbTeams.Exec("CREATE TABLE IF NOT EXISTS results ( scoutid TEXT PRIMARY KEY, campaignid TEXT NOT NULL, eventid TEXT NOT NULL, matchid TEXT NOT NULL, userid TEXT NOT NULL, competitorid TEXT NOT NULL, matchnumber INTEGER NOT NULL, alliance STRING, autoLineCross BIT, autoLowBalls INTEGER, autoHighBalls INTEGER, autoBackBalls INTEGER, autoShots, autoPickups INTEGER, shotQuantity INTEGER, lowFuel INTEGER, highFuel INTEGER, backFuel INTEGER, stageOneComplete BIT, stageOneTime INTEGER, stageTwoComplete BIT, stageTwoTime INTEGER, fouls INTEGER, techFouls INTEGER, card TEXT, climbed TEXT, balanced BIT, climbtime INTEGER, comments TEXT )") // A team's scouted results. Any number of teams may scout for the same campaign / event / match at the same time.

	// Create a default SysAdmin team if it does not exist.

	// Campaigns. Stores information about campaigns but does not store the results associated with them.
	dbCampaigns = newDatabase("campaigns")
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS campaigns ( campaignid TEXT PRIMARY KEY UNIQUE NOT NULL, owner TEXT NOT NULL, name TEXT NOT NULL )")                                            // TODO: Add more information about each campaign. Campaign owner is a teamid. If campaign owner is all zeros, campaign is global.
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS events ( eventid TEXT PRIMARY KEY NOT NULL, campaignid TEXT NOT NULL, name TEXT NOT NULL, location TEXT, starttime INTEGER, endtime INTEGER )") // TODO: Add more information about each event.
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS matches ( matchid TEXT PRIMARY KEY NOT NULL, eventid TEXT NOT NULL, matchnumber INTEGER NOT NULL, active BIT )")                                // TODO: Add more information about each match.
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS pitscout ( pitscoutid TEXT PRIMARY KEY NOT NULL, competitorid TEXT NOT NULL, campaignid TEXT NOT NULL, teamname TEXT, cycletime INTEGER NOT NULL, comments TEXT )")
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS images ( imageid TEXT PRIMARY KEY, campaignid TEXT NOT NULL, competitorid TEXT NOT NULL, image TEXT NOT NULL )")
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS participants ( matchid TEXT PRIMARY KEY NOT NULL, competitorid TEXT UNIQUE NOT NULL) ")             // The participants in each match.
	dbCampaigns.Exec("CREATE TABLE IF NOT EXISTS competitors ( competitorid TEXT PRIMARY KEY NOT NULL, number INTEGER UNIQUE, name TEXT NOT NULL )") // TODO: Add more information about each competing team.

	//Reusing indicator for whether database was just made after all databases are written
	//TODO these are for testing
	if errQuery != nil {
		var teamID, campaignID, eventID string
		TeamCreate(4415, "epic robotz", "nothing")
		dbTeams.QueryRow("SELECT teamid FROM teams").Scan(&teamID)
		CampaignCreate(teamID, "00000000-0000-0000-0000-000000000000", "test")
		dbCampaigns.QueryRow("SELECT campaignid FROM campaigns WHERE owner='00000000-0000-0000-0000-000000000000'").Scan(&campaignID)
		CreateEvent(campaignID, "00000000-0000-0000-0000-000000000000", "event", "nowhere", 0, 900000000000)
		dbCampaigns.QueryRow("SELECT eventid FROM events").Scan(&eventID)
		CreateMatch(eventID, "00000000-0000-0000-0000-000000000000", 1, true)
		dbTeams.Exec(fmt.Sprintf("UPDATE teams SET schedule='%s' WHERE teamid='%s'", campaignID, teamID))
	}
}

/*
TEAM ADMIN FUNCTIONS
*/

// TODO: Teams need a way to weight each of their competitors when assigning teams via the scheduler.

/*
TeamCreate creates a new team from a TeamData struct. Returns bool false and an error if unable to create team.
*/
func TeamCreate(number int, name, schedule string) error {
	teamID := uuid.New().String()
	_, err := dbTeams.Exec(fmt.Sprintf("INSERT INTO teams VALUES ( '%s', '%v', '%s', '%s' )", teamID, number, name, schedule))
	return err
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
TeamListFull returns teamID, teamNumber, teamName, and schedule for every team.
*/
func TeamListFull() map[string][]string {
	rows, err := dbTeams.Query("SELECT teamid, number, name, schedule from teams")
	defer rows.Close()
	accessCheck(err)
	results := make(map[string][]string)
	var id, number, name, schedule string
	for rows.Next() {
		rows.Scan(&id, &number, &name, &schedule)
		results[id] = append(results[id], number, name, schedule)
	}
	return results
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
func UserModify(userID string, data UserData) {

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
UserList returns a list of users as userid and username.
*/
func UserList() map[string]string {
	rows, err := dbUsers.Query("SELECT userid, username FROM users")
	accessCheck(err)
	defer rows.Close()
	results := make(map[string]string)
	var id string
	var uname string
	for rows.Next() {
		rows.Scan(&id, &uname)
		results[id] = uname
	}
	return results
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
GetTeamSchedule gets the event and campaign in which a team is currently participating
*/
func GetTeamSchedule(teamID string) (string, string, error) {
	campaignid, err := GetTeamCampaign(teamID)
	if err != nil {
		return "", "", err
	}
	eventid, err := GetActiveCampaignEvent(campaignid)
	if err != nil {
		return "", "", err
	}
	return campaignid, eventid, nil
}

/*
DATA STORAGE FUNCTIONS
*/

/*
StoreMatch takes the array of data from the form and stores it in the database
*/
func StoreMatch(arr []string, agentid, teamid string) error {
	campaignid, eventid, err := GetTeamSchedule(teamid)
	if err != nil {
		return err
	}
	data, err := arrToMatchStruct(arr, eventid, agentid)
	//don't log these errors here since the function calls should log them on their own
	if err != nil {
		return err
	}
	competitorid := GetCompetitorID(data.Team)
	if competitorid == "" {
		CreateCompetitor(data.Team, "")
		competitorid = GetCompetitorID(data.Team)
	}
	scoutid := uuid.New().String()
	result, errExec := dbTeams.Exec(fmt.Sprintf("INSERT INTO results VALUES ( '%s', '%s', '%s', '%s', '%s', '%s', '%v', '%s', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%s' )", scoutid, campaignid, eventid, data.MatchID, agentid, competitorid, data.MatchNum, data.Alliance, data.AutoLineCross, data.AutoLowBalls, data.AutoHighBalls, data.AutoBackBalls, data.AutoShots, data.AutoPickups, data.ShotQuantity, data.LowFuel, data.HighFuel, data.BackFuel, data.StageOneComplete, data.StageOneTime, data.StageTwoComplete, data.StageTwoTime, data.Fouls, data.TechFouls, data.Card, data.Climbed, data.Balanced, data.ClimbTime, escapeText(data.Comments)))
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
	var alliance, card, climbed string
	//cuts off the last value since comments can't be converted into integers
	intarr, err := convertArrToInts(arr[:len(arr)-1])
	if err != nil {
		return nil, err
	}
	matchid, err := matchIDFromNum(intarr[0], eventid)
	if matchid == "" {
		//TODO figure out if true value on matches is uselful under current system
		err = CreateMatch(eventid, agentid, intarr[0], true)
		matchid, err = matchIDFromNum(intarr[0], eventid)
	}
	autoLineCross = intarr[3] == 1
	stageOneComplete = intarr[13] == 1
	stageTwoComplete = intarr[15] == 1
	balanced := intarr[21] == 1
	if intarr[2] == 1 {
		alliance = "red"
	} else {
		alliance = "blue"
	}
	if intarr[19] == 1 {
		card = "yellow"
	} else if intarr[19] == 2 {
		card = "red"
	} else {
		card = "none"
	}
	if intarr[20] == 2 {
		climbed = "climbed"
	} else if intarr[20] == 1 {
		climbed = "platform"
	} else {
		climbed = "none"
	}
	return &MatchData{MatchID: matchid, MatchNum: intarr[0], Team: intarr[1], Alliance: alliance, AutoLineCross: autoLineCross, AutoLowBalls: intarr[4], AutoHighBalls: intarr[5], AutoBackBalls: intarr[6], AutoShots: intarr[7], AutoPickups: intarr[8], ShotQuantity: intarr[9], LowFuel: intarr[10], HighFuel: intarr[11], BackFuel: intarr[12], StageOneComplete: stageOneComplete, StageOneTime: intarr[14], StageTwoComplete: stageTwoComplete, StageTwoTime: intarr[16], Fouls: intarr[17], TechFouls: intarr[18], Card: card, Climbed: climbed, Balanced: balanced, ClimbTime: intarr[22], Comments: arr[23]}, nil
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
WritePitData writes a pit data entry along with the images that come with it
*/
func WritePitData(arr []string, userID, campaignID string) error {
	teamNum, err := strconv.Atoi(arr[0])
	if err != nil {
		log.Warn(err)
		return err
	}
	competitorID := GetCompetitorID(teamNum)
	//make new competitor if competitor id doesn't exist
	if competitorID == "" {
		CreateCompetitor(teamNum, "")
		competitorID = GetCompetitorID(teamNum)
	}
	pitscoutID := uuid.New().String()
	cycletime, err := strconv.Atoi(arr[2])
	if err != nil {
		log.Warn(err)
		return err
	}
	_, err = dbCampaigns.Exec(fmt.Sprintf("INSERT INTO pitscout VALUES ('%s', '%s', '%s', '%s', '%v', '%s' )", pitscoutID, competitorID, campaignID, arr[1], cycletime, arr[3]))
	if err != nil {
		log.Warn(err)
		return err
	}
	//inserts images into seperate table
	writeImage(competitorID, campaignID, arr[4])
	return nil
}

/*
WriteImage writes a robot image to the database
*/
func writeImage(competitorID, campaignid, image string) error {
	imageid := uuid.New().String()
	_, err := dbCampaigns.Exec(fmt.Sprintf("INSERT INTO images VALUES ('%s', '%s', '%s', '%s' )", imageid, campaignid, competitorID, image))
	return err
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
GetMatchResults gets scouter's data based on a match id
*/
func GetMatchResults(matchID, campaignID string) (*[]MatchData, error) {
	var competitorID string
	data := make([]MatchData, 0)
	rows, err := dbTeams.Query(fmt.Sprintf("SELECT matchnumber, competitorid, autoLineCross, autoLowBalls, autoHighBalls, autoBackBalls, autoPickups, shotQuantity, lowFuel, highFuel, backFuel, stageOneComplete, stageOneTime, stageTwoComplete, stageTwoTime, fouls, techFouls, card, climbed, balanced, climbtime, comments FROM results WHERE matchID='%s'", matchID))
	defer rows.Close()
	for rows.Next() {
		var d MatchData
		err = rows.Scan(&d.MatchNum, &competitorID, &d.AutoLineCross, &d.AutoLowBalls, &d.AutoHighBalls, &d.AutoBackBalls, &d.AutoPickups, &d.ShotQuantity, &d.LowFuel, &d.HighFuel, &d.BackFuel, &d.StageOneComplete, &d.StageOneTime, &d.StageTwoComplete, &d.StageTwoTime, &d.Fouls, &d.TechFouls, &d.Card, &d.Climbed, &d.Balanced, &d.ClimbTime, &d.Comments)
		if err != nil {
			return nil, err
		}
		d.MatchID = matchID
		d.Team = GetCompetitorNumberFromID(competitorID)
	}
	return &data, nil
}

/*
GetTeamResults gets scouter's data based on a team id
*/
func GetTeamResults(teamNum int, campaignID string) (*[]MatchData, error) {
	var matchEvent string
	campaignEvent, _ := GetActiveCampaignEvent(campaignID)
	competitorID := GetCompetitorID(teamNum)
	data := make([]MatchData, 0)
	rows, err := dbTeams.Query(fmt.Sprintf("SELECT matchid, matchnumber, autoLineCross, autoLowBalls, autoHighBalls, autoBackBalls, autoPickups, shotQuantity, lowFuel, highFuel, backFuel, stageOneComplete, stageOneTime, stageTwoComplete, stageTwoTime, fouls, techFouls, card, climbed, balanced, climbtime, comments FROM results WHERE competitorid='%s' AND eventid='%s'", competitorID, campaignEvent))
	defer rows.Close()
	for rows.Next() {
		var d MatchData
		err = rows.Scan(&d.MatchID, &d.MatchNum, &d.AutoLineCross, &d.AutoLowBalls, &d.AutoHighBalls, &d.AutoBackBalls, &d.AutoPickups, &d.ShotQuantity, &d.LowFuel, &d.HighFuel, &d.BackFuel, &d.StageOneComplete, &d.StageOneTime, &d.StageTwoComplete, &d.StageTwoTime, &d.Fouls, &d.TechFouls, &d.Card, &d.Climbed, &d.Balanced, &d.ClimbTime, &d.Comments)
		if err != nil {
			return nil, err
		}
		d.Team = teamNum
		matchEvent, _ = getMatchEvent(d.MatchID)
		if matchEvent == campaignEvent {
			data = append(data, d)
		}
	}
	return &data, nil
}

/*
GetTeamComments gets all comments for a team at a certain event
*/
func GetTeamComments(teamNum int, campaignID string) ([]string, error) {
	var comment string
	comments := make([]string, 0)
	teamID := GetCompetitorID(teamNum)
	eventID, _ := GetActiveCampaignEvent(campaignID)
	row, err := dbTeams.Query(fmt.Sprintf("SELECT comments FROM results WHERE competitorid='%s' AND eventid='%s'", teamID, eventID))
	if err != nil {
		return comments, err
	}
	for row.Next() {
		row.Scan(&comment)
		comments = append(comments, comment)
	}
	return comments, err
}

/*
GetTeamMatchResults gets scouter's data based on a team id for a given match
*/
func GetTeamMatchResults(teamNum int, matchID string) (*[]MatchData, error) {
	data := make([]MatchData, 0)
	competitorID := GetCompetitorID(teamNum)
	rows, err := dbTeams.Query(fmt.Sprintf("SELECT matchnumber, autoLineCross, autoLowBalls, autoHighBalls, autoBackBalls, autoPickups, shotQuantity, lowFuel, highFuel, backFuel, stageOneComplete, stageOneTime, stageTwoComplete, stageTwoTime, fouls, techFouls, card, climbed, balanced, climbtime, comments FROM results WHERE competitorid='%s' AND matchid='%s'", competitorID, matchID))
	defer rows.Close()
	for rows.Next() {
		var d MatchData
		err = rows.Scan(&d.MatchNum, &d.AutoLineCross, &d.AutoLowBalls, &d.AutoHighBalls, &d.AutoBackBalls, &d.AutoPickups, &d.ShotQuantity, &d.LowFuel, &d.HighFuel, &d.BackFuel, &d.StageOneComplete, &d.StageOneTime, &d.StageTwoComplete, &d.StageTwoTime, &d.Fouls, &d.TechFouls, &d.Card, &d.Climbed, &d.Balanced, &d.ClimbTime, &d.Comments)
		if err != nil {
			return nil, err
		}
		d.MatchID = matchID
		d.Team = teamNum
		data = append(data, d)
	}
	return &data, nil
}

/*
GetTeamMatches gets scouter's data based on a team id for a given match
*/
func GetTeamMatches(teamNum int, campaignid string) (*[]MatchData, error) {
	data := make([]MatchData, 0)
	competitorID := GetCompetitorID(teamNum)
	event, _ := GetActiveCampaignEvent(campaignid)
	rows, err := dbTeams.Query(fmt.Sprintf("SELECT matchnumber, matchid, autoLineCross, autoLowBalls, autoHighBalls, autoBackBalls, autoPickups, shotQuantity, lowFuel, highFuel, backFuel, stageOneComplete, stageOneTime, stageTwoComplete, stageTwoTime, fouls, techFouls, card, climbed, balanced, climbtime, comments FROM results WHERE competitorid='%s' AND eventid='%s'", competitorID, event))
	defer rows.Close()
	for rows.Next() {
		var d MatchData
		err = rows.Scan(&d.MatchNum, &d.MatchID, &d.AutoLineCross, &d.AutoLowBalls, &d.AutoHighBalls, &d.AutoBackBalls, &d.AutoPickups, &d.ShotQuantity, &d.LowFuel, &d.HighFuel, &d.BackFuel, &d.StageOneComplete, &d.StageOneTime, &d.StageTwoComplete, &d.StageTwoTime, &d.Fouls, &d.TechFouls, &d.Card, &d.Climbed, &d.Balanced, &d.ClimbTime, &d.Comments)
		if err != nil {
			return nil, err
		}
		d.Team = teamNum
		data = append(data, d)
	}
	return &data, nil
}

/*
GetTeamImages reads all images for a given team
*/
func GetTeamImages(teamNum int, campaignID string) ([]string, error) {
	var image string
	images := make([]string, 0)
	competitorID := GetCompetitorID(teamNum)
	rows, err := dbCampaigns.Query(fmt.Sprintf("SELECT image FROM images WHERE competitorid='%s' AND campaignid='%s'", competitorID, campaignID))
	if err != nil {
		return images, err
	}
	defer rows.Close()
	if rows.Next() {
		rows.Scan(&image)
		images = append(images, image)
	}
	return images, err
}

/*
GetEventResults gets results from all matches in an event
*/
func GetEventResults(event string) (*[]MatchData, error) {
	var matchEvent, competitorid string
	data := make([]MatchData, 0)
	rows, err := dbTeams.Query(fmt.Sprintf("SELECT matchid, matchnumber, competitorid, autoLineCross, autoLowBalls, autoHighBalls, autoBackBalls, autoPickups, shotQuantity, lowFuel, highFuel, backFuel, stageOneComplete, stageOneTime, stageTwoComplete, stageTwoTime, fouls, techFouls, card, climbed, balanced, climbtime, comments FROM results WHERE eventid='%s'", event))
	defer rows.Close()
	for rows.Next() {
		var d MatchData
		err = rows.Scan(&d.MatchID, &d.MatchNum, &competitorid, &d.AutoLineCross, &d.AutoLowBalls, &d.AutoHighBalls, &d.AutoBackBalls, &d.AutoPickups, &d.ShotQuantity, &d.LowFuel, &d.HighFuel, &d.BackFuel, &d.StageOneComplete, &d.StageOneTime, &d.StageTwoComplete, &d.StageTwoTime, &d.Fouls, &d.TechFouls, &d.Card, &d.Climbed, &d.Balanced, &d.ClimbTime, &d.Comments)
		if err != nil {
			return nil, err
		}
		d.Team, _ = GetTeamNumberFromID(competitorid)
		matchEvent, _ = getMatchEvent(d.MatchID)
		if matchEvent == event {
			data = append(data, d)
		}
	}
	return &data, nil
}

/*
GetCurrentEventResults gets results from all matches in an event
*/
func GetCurrentEventResults(campaignid string) (*[]MatchData, error) {
	var matchEvent, competitorid string
	event, _ := GetActiveCampaignEvent(campaignid)
	data := make([]MatchData, 0)
	rows, err := dbTeams.Query(fmt.Sprintf("SELECT matchid, matchnumber, competitorid, autoLineCross, autoLowBalls, autoHighBalls, autoBackBalls, autoPickups, shotQuantity, lowFuel, highFuel, backFuel, stageOneComplete, stageOneTime, stageTwoComplete, stageTwoTime, fouls, techFouls, card, climbed, balanced, climbtime, comments FROM results WHERE eventid='%s'", event))
	defer rows.Close()
	for rows.Next() {
		var d MatchData
		err = rows.Scan(&d.MatchID, &d.MatchNum, &competitorid, &d.AutoLineCross, &d.AutoLowBalls, &d.AutoHighBalls, &d.AutoBackBalls, &d.AutoPickups, &d.ShotQuantity, &d.LowFuel, &d.HighFuel, &d.BackFuel, &d.StageOneComplete, &d.StageOneTime, &d.StageTwoComplete, &d.StageTwoTime, &d.Fouls, &d.TechFouls, &d.Card, &d.Climbed, &d.Balanced, &d.ClimbTime, &d.Comments)
		if err != nil {
			return nil, err
		}
		d.Team, _ = GetTeamNumberFromID(competitorid)
		matchEvent, _ = getMatchEvent(d.MatchID)
		if matchEvent == event {
			data = append(data, d)
		}
	}
	return &data, nil
}

/*
GetCampaignResults gets results from all matches in a campaign
*/
func GetCampaignResults(campaignid string) (*[]MatchData, error) {
	var competitorid string
	data := make([]MatchData, 0)
	rows, err := dbTeams.Query(fmt.Sprintf("SELECT matchid, matchnumber, competitorid, autoLineCross, autoLowBalls, autoHighBalls, autoBackBalls, autoPickups, shotQuantity, lowFuel, highFuel, backFuel, stageOneComplete, stageOneTime, stageTwoComplete, stageTwoTime, fouls, techFouls, card, climbed, balanced, climbtime, comments FROM results WHERE campaignid='%s'", campaignid))
	defer rows.Close()
	for rows.Next() {
		var d MatchData
		err = rows.Scan(&d.MatchID, &d.MatchNum, &competitorid, &d.AutoLineCross, &d.AutoLowBalls, &d.AutoHighBalls, &d.AutoBackBalls, &d.AutoPickups, &d.ShotQuantity, &d.LowFuel, &d.HighFuel, &d.BackFuel, &d.StageOneComplete, &d.StageOneTime, &d.StageTwoComplete, &d.StageTwoTime, &d.Fouls, &d.TechFouls, &d.Card, &d.Climbed, &d.Balanced, &d.ClimbTime, &d.Comments)
		if err != nil {
			return nil, err
		}
		d.Team, _ = GetTeamNumberFromID(competitorid)
		data = append(data, d)
	}
	return &data, nil
}

/*
ListMatchIDs gets a list of match ids from a certain campaign
*/
func ListMatchIDs(campaignID string) []string {
	var matchID string
	matchIDs := make([]string, 0)
	eventID, _ := GetActiveCampaignEvent(campaignID)
	rows, _ := dbCampaigns.Query(fmt.Sprintf("SELECT matchid FROM matches WHERE eventid='%s'", eventID))
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&matchID)
		matchIDs = append(matchIDs, matchID)
	}
	return matchIDs
}

/*
GetTeamNumberFromID gets a teams number from their competitor id
*/
func GetTeamNumberFromID(teamID string) (int, error) {
	var number int
	err := dbCampaigns.QueryRow(fmt.Sprintf("SELECT number FROM competitors WHERE competitorid='%s'", teamID)).Scan(&number)
	return number, err
}

/*
GetTeamID gets team uuid
*/
func GetTeamID(number int) (string, error) {
	var teamID string
	err := dbTeams.QueryRow(fmt.Sprintf("SELECT teamid FROM teams WHERE number='%v'", number)).Scan(&teamID)
	return teamID, err
}

func getMatchEvent(matchID string) (string, error) {
	var eventID string
	err := dbCampaigns.QueryRow(fmt.Sprintf("SELECT eventid FROM matches WHERE matchid='%s'", matchID)).Scan(&eventID)
	return eventID, err
}

/*
GetMatchParticipants gets teams participating on each alliance in a match
*/
func GetMatchParticipants(matchID string) [][]int {
	participants := make([][]int, 2)
	participants[0] = GetAllianceParticipants(matchID, "red")
	participants[1] = GetAllianceParticipants(matchID, "blue")
	return participants
}

/*
GetAllianceParticipants gets teams participating in a particular alliance in a match
*/
func GetAllianceParticipants(matchID, alliance string) []int {
	var teamID string
	var teamNum int
	allies := make([]int, 0)
	rows, _ := dbTeams.Query(fmt.Sprintf("SELECT competitorid FROM results WHERE matchid='%s' AND alliance='%s'", matchID, alliance))
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&teamID)
		teamNum, _ = GetTeamNumberFromID(teamID)
		if !contains(allies, teamNum) {
			allies = append(allies, teamNum)
		}
	}
	return allies
}

func contains(arr []int, val int) bool {
	for _, x := range arr {
		if x == val {
			return true
		}
	}
	return false
}

/*
ListAllCompetitors returns a list of all competitor ids
*/
func ListAllCompetitors() ([]string, error) {
	var row string
	var err error
	competitors := make([]string, 0)
	rows, _ := dbTeams.Query(fmt.Sprintf("SELECT competitorid FROM competitors"))
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&row)
		if err != nil {
			return competitors, err
		}
		competitors = append(competitors, row)
	}
	return competitors, nil
}

/*
TEAM ADMIN FUNCTIONS
*/

/*
CAMPAIGN FUCTIONS
*/

/*
GetActiveCampaignEvent gets the eventid of the active event in the given campaign
*/
func GetActiveCampaignEvent(campaignid string) (string, error) {
	var eventid, cid string
	var starttime, endtime int64
	rows, err := dbCampaigns.Query("SELECT eventid, campaignid, starttime, endtime FROM events")
	defer rows.Close()
	//checks if event is currently taking place
	for rows.Next() {
		err = rows.Scan(&eventid, &cid, &starttime, &endtime)
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
GetEventMatchIDs gets a list of matchids from a given event
*/
func GetEventMatchIDs(eventid string) []string {
	var matchid string
	matchids := make([]string, 0)
	rows, _ := dbCampaigns.Query(fmt.Sprintf("SELECT matchid FROM matches WHERE eventid='%s'", eventid))
	defer rows.Close()
	for rows.Next() {
		rows.Scan(&matchid)
		matchids = append(matchids, matchid)
	}
	return matchids
}

/*
CampaignCreate TODO.
*/
func CampaignCreate(agentid, owner, name string) {
	// TODO: Clone global campaigns to team-specific campaign if requested.
	// Only sysadmin can create global campaigns.
	uuid := uuid.New().String()
	dbCampaigns.Exec(fmt.Sprintf("INSERT INTO campaigns VALUES ( '%s', '%s', '%s' )", uuid, owner, name))
}

/*
CampaignClone TODO
*/
func CampaignClone(agentID, clonedID, teamID string) {

}

/*
CampaignList returns list of campaigns as id: owner, name
*/
func CampaignList() map[string][]string {
	rows, err := dbCampaigns.Query("SELECT campaignid, owner, name FROM campaigns")
	defer rows.Close()
	accessCheck(err)
	results := make(map[string][]string)
	var id, owner, name string
	for rows.Next() {
		rows.Scan(&id, &owner, &name)
		results[id] = append(results[id], owner, name)
	}
	return results
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
	_, err := dbCampaigns.Exec(fmt.Sprintf("INSERT INTO matches VALUES ( '%s', '%s', '%v', '%v' )", matchid, eventid, num, active))
	if err == nil {
		log.Infof("Created match #%v for event %s", num, eventid)
	}
	return err
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

/*
ResultsWrite writes
*/
func ResultsWrite(useriD, teamID string, data MatchData) {
	// TODO
	// Throw if overwriting existing
}

/*
ResultsRead reads match results
*/
func ResultsRead() {

}

/*
CreateCompetitor creates a competitor
*/
func CreateCompetitor(teamNumber int, name string) {
	competitorID := uuid.New().String()
	dbCampaigns.Exec(fmt.Sprintf("INSERT INTO competitors VALUES ( '%s', '%v', '%s' )", competitorID, teamNumber, name))
}

/*
GetCompetitorID gets competitor id for team number
*/
func GetCompetitorID(teamNumber int) string {
	var competitorID string
	dbCampaigns.QueryRow(fmt.Sprintf("SELECT competitorid FROM competitors WHERE number='%v'", teamNumber)).Scan(&competitorID)
	return competitorID
}

/*
GetCompetitorNumberFromID gets competitor team number from id
*/
func GetCompetitorNumberFromID(competitorID string) int {
	var teamNum int
	dbCampaigns.QueryRow(fmt.Sprintf("SELECT number FROM competitors WHERE competitorid='%s'", competitorID)).Scan(&teamNum)
	return teamNum
}

/*
competitorIDNumber gets a team number from competitor id
*/
func competitorIDNumber(competitorID string) int {
	var teamNumber int
	dbTeams.QueryRow(fmt.Sprintf("SELECT number FROM competitors WHERE competitorid='%v'", competitorID)).Scan(&teamNumber)
	return teamNumber
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
SysAdminList returns a list of sysadmins as userIDs and usernames.
*/
func SysAdminList() map[string]string {
	rows, err := dbUsers.Query("SELECT * from sysadmins")
	accessCheck(err)
	defer rows.Close()
	var id string
	sysAdmins := make(map[string]string)
	for rows.Next() {
		rows.Scan(&id)
		d, err := UserQuery(id)
		accessCheck(err)
		sysAdmins[id] = d.UserName
	}
	return sysAdmins
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
