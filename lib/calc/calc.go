/*
Package calc provides functions for calculating various game statistics.
*/
package calc

import (
	"fmt"
	"math"
	"sort"

	"EPIC-Scouting/lib/db"
)

/*
MatchResults summary of match results. A culmination of all data on a given match
*/
type MatchResults struct {
	MatchNum            int
	RedParticipants     []int
	BlueParticipants    []int
	RedPoints           int
	BluePoints          int
	Winner              string
	RedRankingPoints    int
	BlueRankingPoints   int
	RedAutoLineCrosses  int
	BlueAutoLineCrosses int
	RedAutoPoints       int
	BlueAutoPoints      int
	RedAutoBalls        int
	BlueAutoBalls       int
	RedShootingPoints   int
	BlueShootingPoints  int
	RedTeleopShots      int
	BlueTeleopShots     int
	RedLowShots         int
	BlueLowShots        int
	RedHighShots        int
	BlueHighShots       int
	RedBackShots        int
	BlueBackShots       int
	RedShieldStage      int
	BlueShieldStage     int
	RedClimbStatus      []int
	BlueClimbStatus     []int
	RedBalanced         bool
	BlueBalanced        bool
	RedClimbPoints      int
	BlueClimbPoints     int
}

/*
GetTeamScores gets all team breakdown scores from the active event in a campaign
*/
func GetTeamScores(campaignid string) [][]int {
	scores := make([][]int, 0)
	teamData := make(map[int][]db.MatchData, 0)
	data, _ := db.GetCurrentEventResults(campaignid)
	for _, match := range *data {
		_, ok := teamData[match.Team]
		if !ok {
			teamData[match.Team] = make([]db.MatchData, 0)
		}
		teamData[match.Team] = append(teamData[match.Team], match)
	}
	for team, matches := range teamData {
		scores = append(scores, []int{team, Overall(matches), Auto(matches), Shooting(matches), ColorWheel(matches), Climbing(matches), Foul(matches)})
	}
	return scores
}

/*
TODO these should be match summary functions
*/

/*
GetMatchData gets a summary of match scores for the red and blue alliances respectively
*/
func GetMatchData(matchid string) (MatchResults, error) {
	var matchParticipants [][]int
	var results MatchResults
	var participantScores []db.MatchData
	var scores [][]db.MatchData
	var err error
	matchParticipants = db.GetMatchParticipants(matchid)
	for alliance := range matchParticipants {
		participantScores = make([]db.MatchData, 0)
		for _, team := range matchParticipants[alliance] {
			participantScores = append(participantScores, ResolveMatchConflicts(team, matchid))
		}
		scores = append(scores, participantScores)
	}
	results, err = deriveMatchScores(scores[0], scores[1])
	return results, err
}

func deriveMatchScores(red, blue []db.MatchData) (MatchResults, error) {
	var summary MatchResults
	var count, redPoints, bluePoints, redRP, blueRP int
	//TODO this is only commented for testing
	// if len(red) == 0 || len(blue) == 0 {
	// 	return summary, errors.New("Unable to summarize match: no data provided for one or more alliances")
	// }
	fmt.Println(red, blue)
	summary.MatchNum = red[0].MatchNum
	participants := db.GetMatchParticipants(red[0].MatchID)
	summary.RedParticipants = participants[0]
	summary.BlueParticipants = participants[1]
	for _, teamdata := range red {
		if teamdata.AutoLineCross {
			count++
		}
	}
	summary.RedAutoLineCrosses = count
	count = 0
	for _, teamdata := range blue {
		if teamdata.AutoLineCross {
			count++
		}
	}
	summary.BlueAutoLineCrosses = count
	count = 0
	for _, teamdata := range red {
		if teamdata.AutoLineCross {
			//add points for auto line cross
			count += 5
		}
		count += 2*teamdata.AutoLowBalls + 4*teamdata.AutoHighBalls + 6*teamdata.AutoBackBalls
	}
	summary.RedAutoPoints = count
	redPoints += count
	count = 0
	for _, teamdata := range blue {
		if teamdata.AutoLineCross {
			//add points for auto line cross
			count += 5
		}
		count += 2*teamdata.AutoLowBalls + 4*teamdata.AutoHighBalls + 6*teamdata.AutoBackBalls
	}
	summary.BlueAutoPoints = count
	bluePoints += count
	count = 0
	for _, teamdata := range red {
		count += teamdata.AutoLowBalls + teamdata.AutoHighBalls + teamdata.AutoBackBalls
	}
	summary.RedAutoBalls = count
	count = 0
	for _, teamdata := range blue {
		count += teamdata.AutoLowBalls + teamdata.AutoHighBalls + teamdata.AutoBackBalls
	}
	summary.BlueAutoBalls = count
	count = 0
	for _, teamdata := range red {
		count += teamdata.LowFuel + teamdata.HighFuel*2 + teamdata.BackFuel*3
	}
	summary.RedShootingPoints = count
	redPoints += count
	count = 0
	for _, teamdata := range blue {
		count += teamdata.LowFuel + teamdata.HighFuel*2 + teamdata.BackFuel*3
	}
	summary.BlueShootingPoints = count
	bluePoints += count
	count = 0
	for _, teamdata := range red {
		count += teamdata.LowFuel + teamdata.HighFuel + teamdata.BackFuel
	}
	summary.RedTeleopShots = count
	count = 0
	for _, teamdata := range blue {
		count += teamdata.LowFuel + teamdata.HighFuel + teamdata.BackFuel
	}
	summary.BlueTeleopShots = count
	count = 0
	for _, teamdata := range red {
		count += teamdata.LowFuel
	}
	summary.RedLowShots = count
	count = 0
	for _, teamdata := range blue {
		count += teamdata.LowFuel
	}
	summary.BlueLowShots = count
	count = 0
	for _, teamdata := range red {
		count += teamdata.HighFuel
	}
	summary.RedHighShots = count
	count = 0
	for _, teamdata := range blue {
		count += teamdata.HighFuel
	}
	summary.BlueHighShots = count
	count = 0
	for _, teamdata := range red {
		count += teamdata.BackFuel
	}
	summary.RedBackShots = count
	count = 0
	for _, teamdata := range blue {
		count += teamdata.BackFuel
	}
	summary.BlueBackShots = count
	count = 0
	//TODO stage one and two should be stages two and three
	//TODO update point values for each stage complete
	if red[0].StageOneComplete {
		if red[0].StageTwoComplete {
			summary.RedShieldStage = 3
			redPoints += 50
			redRP++
		} else {
			summary.RedShieldStage = 2
			redPoints += 30
		}
	} else if summary.RedAutoBalls+summary.RedTeleopShots >= 20 {
		summary.RedShieldStage = 1
		redPoints += 10 //???
	}
	//TODO uncomment below
	// if blue[0].StageOneComplete {
	// 	if blue[0].StageTwoComplete {
	// 		summary.BlueShieldStage = 3
	// 		bluePoints += 50
	// 		blueRP++
	// 	} else {
	// 		summary.BlueShieldStage = 2
	// 		bluePoints += 30
	// 	}
	// } else if summary.BlueAutoBalls+summary.BlueTeleopShots >= 20 {
	// 	summary.BlueShieldStage = 1
	// 	bluePoints += 10 //???
	// }
	for _, teamdata := range red {
		if teamdata.Climbed == "climbed" {
			summary.RedClimbStatus = append(summary.RedClimbStatus, 2)
			count += 25
		} else if teamdata.Climbed == "platform" {
			summary.RedClimbStatus = append(summary.RedClimbStatus, 1)
			count += 5
		} else {
			summary.RedClimbStatus = append(summary.RedClimbStatus, 0)
		}
	}
	if red[0].Balanced {
		summary.RedBalanced = true
		count += 15
	}
	//award climbing ranking point if condition met
	if count >= 60 {
		redRP++
	}
	summary.RedClimbPoints = count
	redPoints += count
	count = 0
	for _, teamdata := range blue {
		if teamdata.Climbed == "climbed" {
			summary.BlueClimbStatus = append(summary.RedClimbStatus, 2)
			count += 25
		} else if teamdata.Climbed == "platform" {
			summary.BlueClimbStatus = append(summary.RedClimbStatus, 1)
			count += 5
		} else {
			summary.BlueClimbStatus = append(summary.RedClimbStatus, 0)
		}
	}
	//TODO uncomment below
	// if blue[0].Balanced {
	// 	summary.BlueBalanced = true
	// 	count += 15
	// }
	//award climbing ranking point if condition met
	if count >= 60 {
		blueRP++
	}
	summary.BlueClimbPoints = count
	bluePoints += count
	count = 0
	summary.RedPoints = redPoints
	summary.BluePoints = bluePoints
	if redPoints > bluePoints {
		summary.Winner = "red"
		redRP += 2
	} else {
		summary.Winner = "blue"
		blueRP += 2
	}
	summary.RedRankingPoints = redRP
	summary.BlueRankingPoints = blueRP
	return summary, nil
}

//RawTeamEventData gets a team's raw statistics for an event - best for putting on spreadsheets for raw comparison/printout

/*All calculation functions below can be set to include or exclude certain data based on time to allow display of development of scores over time
Team scoring devices never affect each other and are measured against an ideal target. They are then used for computing a team's overall rank
Relative category scores calculate a robot's score compared to the best preformer in that category*/

//TODO: make an external reference to the weight of each element on the composite scores

//TeamOverall gets a teams overall score based off a weight table yet to be implemented
func TeamOverall(teamNum int, campaignid string) int {
	auto := TeamAuto(teamNum, campaignid)
	shooting := TeamShooting(teamNum, campaignid)
	climbing := TeamClimbing(teamNum, campaignid)
	colorWheel := TeamColorWheel(teamNum, campaignid)
	foul := TeamFoul(teamNum, campaignid)
	overall := auto + shooting + climbing + colorWheel - foul
	return overall
}

//TeamAuto gets a team's autonomous rating
func TeamAuto(teamNum int, campaignid string) int {
	teamID := db.GetCompetitorID(teamNum)
	if teamID == "" {
		return 0
	}
	breakdown := TeamAutoBreakdown(teamNum, campaignid)
	weights := []int{5, 4, 2, 1, 1, 1, 1}
	score := 0
	for ind, weight := range weights {
		score += breakdown[ind] * weight
	}
	return score
}

//TeamShooting gets a team's overall shooting score
func TeamShooting(teamNum int, campaignid string) int {
	teamID := db.GetCompetitorID(teamNum)
	if teamID == "" {
		return 0
	}
	breakdown := TeamShootingBreakdown(teamNum, campaignid)
	weights := []int{1, 2, 3, 5, 3, 2, 1}
	score := 0
	for ind, weight := range weights {
		score += breakdown[ind] * weight
	}
	return score
}

//TeamClimbing gets a team's score for climbing
func TeamClimbing(teamNum int, campaignid string) int {
	teamID := db.GetCompetitorID(teamNum)
	if teamID == "" {
		return 0
	}
	breakdown := TeamClimbingBreakdown(teamNum, campaignid)
	weights := []int{2, 1, 1}
	score := 0
	for ind, weight := range weights {
		score += breakdown[ind] * weight
	}
	return score
}

//TeamColorWheel gets how good a team is at manipulating the color wheel
func TeamColorWheel(teamNum int, campaignid string) int {
	teamID := db.GetCompetitorID(teamNum)
	if teamID == "" {
		return 0
	}
	breakdown := TeamColorWheelBreakdown(teamNum, campaignid)
	weights := []int{1, 1}
	score := 0
	for ind, weight := range weights {
		score += breakdown[ind] * weight
	}
	return score
}

//TeamFoul gets how many penalties a team accrues. Extra weight to yellow cards and tech fouls
func TeamFoul(teamNum int, campaignid string) int {
	teamID := db.GetCompetitorID(teamNum)
	if teamID == "" {
		return 0
	}
	breakdown := TeamFoulBreakdown(teamNum, campaignid)
	weights := []int{1, 3, 2, 2}
	score := 0
	for ind, weight := range weights {
		score += breakdown[ind] * weight
	}
	return score
}

//RelativeAuto gets a team's autonomous rating relative to the highest scoring contestant

//RelativeShooting gets a team's overall shooting score relative to the highest scoring contestant

//RelativeClimbing gets a team's score for climbing relative to the highest scoring contestant

//RelativeColorWheel gets how good a team is at manipulating the color wheel relative to the highest scoring contestant

//RelativePenalty gets how many penalties a team accrues relative to the lowest scoring contestant

/*Team scoring breakdowns return the numbers that go into calculating the factors above. They don't affect the team's overall score directly. They should be used to get a better idea of why a score is a certain value and what a robot is actually good at.
They are also what the above functions use to get their data*/

//TeamAutoBreakdown gets a team's ability to cross the auto line, amount of balls scored in auto, auto accuracy, and ammount of points scored in auto
//TODO: Finish this
func TeamAutoBreakdown(teamNum int, campaignID string) []int {
	breakdown := make([]int, 8)
	//matches is a list of the results struct
	//TODO: Get eventid from active event on the given campaignid
	matches, _ := db.GetTeamResults(teamNum, campaignID)
	//totals the scores the team has accumulated over the matches
	for _, match := range *matches {
		if match.AutoLineCross {
			breakdown[0]++
		}
		breakdown[1] += match.AutoBackBalls
		breakdown[2] += match.AutoHighBalls
		breakdown[3] += match.AutoLowBalls
		breakdown[4] += match.AutoShots
		breakdown[5] += match.AutoPickups
		if match.AutoShots-match.AutoLowBalls != 0 {
			breakdown[6] += (match.AutoBackBalls + match.AutoHighBalls) / (match.AutoShots - match.AutoLowBalls)
		}
		//total auto points
		breakdown[7] += breakdown[0]*15 + match.AutoBackBalls*6 + match.AutoHighBalls*4 + match.AutoLowBalls*2
	}
	for ind, val := range breakdown {
		breakdown[ind] = int(math.Round(float64(val) / float64(len(*matches))))
	}
	return breakdown
}

//TeamShootingBreakdown gets a team's teleop shooting rate, shooting accuracy, ball score rate, and point score rate
func TeamShootingBreakdown(teamNum int, campaignid string) []int {
	breakdown := make([]int, 7)
	//matches is a list of the results struct
	//TODO: Get eventid from active event on the given campaignid
	matches, _ := db.GetTeamResults(teamNum, campaignid)
	//totals the scores the team has accumulated over the matches
	for _, match := range *matches {
		breakdown[0] += match.ShotQuantity
		breakdown[1] += match.LowFuel
		breakdown[2] += match.HighFuel
		breakdown[3] += match.BackFuel
		if match.ShotQuantity-match.LowFuel != 0 {
			breakdown[4] += (match.HighFuel + match.BackFuel) / (match.ShotQuantity - match.LowFuel)
		}
		breakdown[5] += match.LowFuel + match.HighFuel + match.BackFuel
		//total auto points
		breakdown[6] += match.LowFuel*1 + match.HighFuel*2 + match.BackFuel*3
	}
	for ind, val := range breakdown {
		breakdown[ind] = int(math.Round(float64(val) / float64(len(*matches))))
	}
	return breakdown
}

//TeamClimbingBreakdown gets a team's average climbing speed, ability to balance the bar, and average points scored for climbing
func TeamClimbingBreakdown(teamNum int, campaignID string) []int {
	breakdown := make([]int, 3)
	//matches is a list of the results struct
	//TODO: Get eventid from active event on the given campaignid
	matches, _ := db.GetTeamResults(teamNum, campaignID)
	//totals the scores the team has accumulated over the matches
	for _, match := range *matches {
		if match.Climbed == "climbed" {
			breakdown[0] += 2
		} else if match.Climbed == "platform" {
			breakdown[0]++
		}
		if !(match.ClimbTime == 0) {
			breakdown[1] += 100 / match.ClimbTime
		}
		if match.Balanced {
			breakdown[2]++
		}
	}
	for ind, val := range breakdown {
		breakdown[ind] = int(math.Round(float64(val) / float64(len(*matches))))
	}
	return breakdown
}

//TeamColorWheelBreakdown gets how quickly a team can do stage 1 and 2 of the color wheel, along with whether they can do it at all
func TeamColorWheelBreakdown(teamNum int, campaignID string) []int {
	breakdown := make([]int, 2)
	//matches is a list of the results struct
	//TODO: Get eventid from active event on the given campaignid
	matches, _ := db.GetTeamResults(teamNum, campaignID)
	//totals the scores the team has accumulated over the matches
	for _, match := range *matches {
		breakdown[0] += match.StageOneTime
		breakdown[1] += match.StageTwoTime
	}
	for ind, val := range breakdown {
		breakdown[ind] = int(math.Round(float64(val) / float64(len(*matches))))
	}
	return breakdown
}

//TeamFoulBreakdown gets how many times a team has recieved regular fouls, tech fouls, and yellow cards, along with the total amount of points lost by them to fouls
func TeamFoulBreakdown(teamNum int, campaignID string) []int {
	breakdown := make([]int, 4)
	//matches is a list of the results struct
	//TODO: Get eventid from active event on the given campaignid
	matches, _ := db.GetTeamResults(teamNum, campaignID)
	//totals the scores the team has accumulated over the matches
	for _, match := range *matches {
		breakdown[0] += match.Fouls
		breakdown[1] += match.TechFouls
		breakdown[2] += match.Fouls*3 + match.TechFouls*15
		//0=no card, 1=yellow card, 2=red card
		if match.Card == "red" {
			breakdown[3] += 2
		} else if match.Card == "yellow" {
			breakdown[3]++
		}
	}
	for ind, val := range breakdown {
		breakdown[ind] = int(math.Round(float64(val) / float64(len(*matches))))
	}
	return breakdown
}

//Overall gets a teams overall score based off a weight table yet to be implemented
func Overall(matches []db.MatchData) int {
	auto := Auto(matches)
	shooting := Shooting(matches)
	climbing := Climbing(matches)
	colorWheel := ColorWheel(matches)
	foul := Foul(matches)
	overall := auto + shooting + climbing + colorWheel - foul
	return overall
}

//Auto gets a team's autonomous rating
func Auto(matches []db.MatchData) int {
	breakdown := AutoBreakdown(matches)
	weights := []int{5, 4, 2, 1, 1, 1, 1}
	score := 0
	for ind, weight := range weights {
		score += breakdown[ind] * weight
	}
	return score
}

//Shooting gets a team's overall shooting score
func Shooting(matches []db.MatchData) int {
	breakdown := ShootingBreakdown(matches)
	weights := []int{1, 2, 3, 5, 3, 2, 1}
	score := 0
	for ind, weight := range weights {
		score += breakdown[ind] * weight
	}
	return score
}

//Climbing gets a team's score for climbing
func Climbing(matches []db.MatchData) int {
	breakdown := ClimbingBreakdown(matches)
	weights := []int{2, 1, 1}
	score := 0
	for ind, weight := range weights {
		score += breakdown[ind] * weight
	}
	return score
}

//ColorWheel gets how good a team is at manipulating the color wheel
func ColorWheel(matches []db.MatchData) int {
	breakdown := ColorWheelBreakdown(matches)
	weights := []int{1, 1}
	score := 0
	for ind, weight := range weights {
		score += breakdown[ind] * weight
	}
	return score
}

//Foul gets how many penalties a team accrues. Extra weight to yellow cards and tech fouls
func Foul(matches []db.MatchData) int {
	breakdown := FoulBreakdown(matches)
	weights := []int{1, 3, 2, 2}
	score := 0
	for ind, weight := range weights {
		score += breakdown[ind] * weight
	}
	return score
}

//AutoBreakdown gets a team's ability to cross the auto line, amount of balls scored in auto, auto accuracy, and ammount of points scored in auto
//TODO: Finish this
func AutoBreakdown(matches []db.MatchData) []int {
	breakdown := make([]int, 8)
	//matches is a list of the results struct
	//totals the scores the team has accumulated over the matches
	for _, match := range matches {
		if match.AutoLineCross {
			breakdown[0]++
		}
		breakdown[1] += match.AutoBackBalls
		breakdown[2] += match.AutoHighBalls
		breakdown[3] += match.AutoLowBalls
		breakdown[4] += match.AutoShots
		breakdown[5] += match.AutoPickups
		if match.AutoShots-match.AutoLowBalls != 0 {
			breakdown[6] += (match.AutoBackBalls + match.AutoHighBalls) / (match.AutoShots - match.AutoLowBalls)
		}
		//total auto points
		breakdown[7] += breakdown[0]*15 + match.AutoBackBalls*6 + match.AutoHighBalls*4 + match.AutoLowBalls*2
	}
	for ind, val := range breakdown {
		breakdown[ind] = int(math.Round(float64(val) / float64(len(matches))))
	}
	return breakdown
}

//ShootingBreakdown gets a team's teleop shooting rate, shooting accuracy, ball score rate, and point score rate
func ShootingBreakdown(matches []db.MatchData) []int {
	breakdown := make([]int, 7)
	//matches is a list of the results struct
	//totals the scores the team has accumulated over the matches
	for _, match := range matches {
		breakdown[0] += match.ShotQuantity
		breakdown[1] += match.LowFuel
		breakdown[2] += match.HighFuel
		breakdown[3] += match.BackFuel
		if match.ShotQuantity-match.LowFuel != 0 {
			breakdown[4] += (match.HighFuel + match.BackFuel) / (match.ShotQuantity - match.LowFuel)
		}
		breakdown[5] += match.LowFuel + match.HighFuel + match.BackFuel
		//total auto points
		breakdown[6] += match.LowFuel*1 + match.HighFuel*2 + match.BackFuel*3
	}
	for ind, val := range breakdown {
		breakdown[ind] = int(math.Round(float64(val) / float64(len(matches))))
	}
	return breakdown
}

//ClimbingBreakdown gets a team's average climbing speed, ability to balance the bar, and average points scored for climbing
func ClimbingBreakdown(matches []db.MatchData) []int {
	breakdown := make([]int, 3)
	//matches is a list of the results struct
	//totals the scores the team has accumulated over the matches
	for _, match := range matches {
		if match.Climbed == "climbed" {
			breakdown[0] += 2
		} else if match.Climbed == "platform" {
			breakdown[0]++
		}
		if !(match.ClimbTime == 0) {
			breakdown[1] += 100 / match.ClimbTime
		}
		if match.Balanced {
			breakdown[2] += 1 / len(matches)
		}
	}
	for ind, val := range breakdown {
		breakdown[ind] = int(math.Round(float64(val) / float64(len(matches))))
	}
	return breakdown
}

//ColorWheelBreakdown gets how quickly a team can do stage 1 and 2 of the color wheel, along with whether they can do it at all
func ColorWheelBreakdown(matches []db.MatchData) []int {
	breakdown := make([]int, 2)
	//matches is a list of the results struct
	//totals the scores the team has accumulated over the matches
	for _, match := range matches {
		breakdown[0] += match.StageOneTime
		breakdown[1] += match.StageTwoTime
	}
	for ind, val := range breakdown {
		breakdown[ind] = int(math.Round(float64(val) / float64(len(matches))))
	}
	return breakdown
}

//FoulBreakdown gets how many times a team has recieved regular fouls, tech fouls, and yellow cards, along with the total amount of points lost by them to fouls
func FoulBreakdown(matches []db.MatchData) []int {
	breakdown := make([]int, 4)
	//matches is a list of the results struct
	//totals the scores the team has accumulated over the matches
	for _, match := range matches {
		breakdown[0] += match.Fouls
		breakdown[1] += match.TechFouls
		breakdown[2] += match.Fouls*3 + match.TechFouls*15
		//0=no card, 1=yellow card, 2=red card
		if match.Card == "red" {
			breakdown[3] += 2
		} else if match.Card == "yellow" {
			breakdown[3]++
		}
	}
	for ind, val := range breakdown {
		breakdown[ind] = int(math.Round(float64(val) / float64(len(matches))))
	}
	return breakdown
}

//Team Overall Scoring and Ranking functions give teams conglomerate scores such as OPR, DPR, and overall ranking

//TeamOverallEvent gives a team an overall quality score

//RankEventTeams ranks all teams from best to worst based on their overall score

//Match Filtering Functions

//filterTeamMatchesBefore returns all match data from a team at an event before a given match number. Includes the given match number

//filterTeamMatchesAfter returns all match data from a team at an event after a given match number. Includes the given match number

//getTeamMatches returns all data from a team at an event

/*Scouter Ranking Functions rank scouts based on their accuracy*/

//RankScouterEvent

//RankScouterGlobal

/*
Match Summary functions use the scouter data on matches to summarize their results
*/

/*Match Census Functions determine the weight of contradictary data on the same match and return a score useable for the system*/
//Below are differing census methods. They may or may not be used.

/*
ResolveMatchList resolves a list of scouter data on various matches into their resolved versions
*/
func ResolveMatchList(matches []db.MatchData) []db.MatchData {
	resolved := make([]db.MatchData, 0)
	numberedMatches := make(map[int][]db.MatchData)
	for _, data := range matches {
		_, ok := numberedMatches[data.MatchNum]
		if ok {
			numberedMatches[data.MatchNum] = append(numberedMatches[data.MatchNum], data)
		} else {
			numberedMatches[data.MatchNum] = make([]db.MatchData, 1)
			numberedMatches[data.MatchNum][0] = data
		}
	}
	for _, data := range numberedMatches {
		resolved = append(resolved, ResolveDataConflicts(data))
	}
	return resolved
}

/*
ResolveMatchConflicts takes multiple scouter's data that may contradict and combines it, eliminating outliers
*/
func ResolveMatchConflicts(teamNum int, matchid string) db.MatchData {
	var resolved db.MatchData
	data, _ := db.GetTeamMatchResults(teamNum, matchid)
	resolved = ResolveDataConflicts(*data)
	return resolved
}

/*
ResolveDataConflicts resolves discrepencies between scouting data
*/
func ResolveDataConflicts(data []db.MatchData) db.MatchData {
	var resolved db.MatchData
	if len(data) == 0 {
		return resolved
	}
	autoLowBallsList := make([]int, len(data))
	autoHighBallsList := make([]int, len(data))
	autoBackBallsList := make([]int, len(data))
	autoShotsList := make([]int, len(data))
	autoPickupsList := make([]int, len(data))
	shotQuantityList := make([]int, len(data))
	lowFuelList := make([]int, len(data))
	highFuelList := make([]int, len(data))
	backFuelList := make([]int, len(data))
	stageOneTimeList := make([]int, len(data))
	stageTwoTimeList := make([]int, len(data))
	foulsList := make([]int, len(data))
	techFoulsList := make([]int, len(data))
	climbTimeList := make([]int, len(data))
	autoLineCrossList := make([]bool, len(data))
	stageOneCompleteList := make([]bool, len(data))
	stageTwoCompleteList := make([]bool, len(data))
	balancedList := make([]bool, len(data))
	cardList := make([]string, len(data))
	climbedList := make([]string, len(data))
	for ind, d := range data {
		autoLowBallsList[ind] = d.AutoLowBalls
		autoHighBallsList[ind] = d.AutoHighBalls
		autoBackBallsList[ind] = d.AutoBackBalls
		autoShotsList[ind] = d.AutoShots
		autoPickupsList[ind] = d.AutoPickups
		shotQuantityList[ind] = d.ShotQuantity
		lowFuelList[ind] = d.LowFuel
		highFuelList[ind] = d.HighFuel
		backFuelList[ind] = d.BackFuel
		stageOneTimeList[ind] = d.StageOneTime
		stageTwoTimeList[ind] = d.StageTwoTime
		foulsList[ind] = d.Fouls
		techFoulsList[ind] = d.TechFouls
		climbTimeList[ind] = d.ClimbTime
		autoLineCrossList[ind] = d.AutoLineCross
		stageOneCompleteList[ind] = d.StageOneComplete
		stageTwoCompleteList[ind] = d.StageTwoComplete
		balancedList[ind] = d.Balanced
		cardList[ind] = d.Card
		climbedList[ind] = d.Climbed
	}
	resolved.AutoLowBalls = resolveInt(autoLowBallsList)
	resolved.AutoHighBalls = resolveInt(autoHighBallsList)
	resolved.AutoBackBalls = resolveInt(autoBackBallsList)
	resolved.AutoShots = resolveInt(autoShotsList)
	resolved.AutoPickups = resolveInt(autoPickupsList)
	resolved.ShotQuantity = resolveInt(shotQuantityList)
	resolved.LowFuel = resolveInt(lowFuelList)
	resolved.HighFuel = resolveInt(highFuelList)
	resolved.BackFuel = resolveInt(backFuelList)
	resolved.StageOneTime = resolveInt(stageOneTimeList)
	resolved.StageTwoTime = resolveInt(stageTwoTimeList)
	resolved.Fouls = resolveInt(foulsList)
	resolved.TechFouls = resolveInt(techFoulsList)
	resolved.ClimbTime = resolveIntMean(climbTimeList)
	resolved.AutoLineCross = resolveBool(autoLineCrossList)
	resolved.StageOneComplete = resolveBool(stageOneCompleteList)
	resolved.StageTwoComplete = resolveBool(stageTwoCompleteList)
	resolved.Balanced = resolveBool(balancedList)
	resolved.Card = resolveString(cardList)
	resolved.Climbed = resolveString(climbedList)
	resolved.Team = data[0].Team
	resolved.MatchID = data[0].MatchID
	return resolved
}

func resolveInt(arr []int) int {
	var resolved, maxCount int
	occuranceCount := make(map[int]int)
	for _, val := range arr {
		_, ok := occuranceCount[val]
		if !ok {
			occuranceCount[val] = 1
		}
		occuranceCount[val]++
	}
	for key, val := range occuranceCount {
		if val > maxCount {
			maxCount = val
			resolved = key
		}
	}
	return resolved
}

func resolveIntMean(arr []int) int {
	var resolved, total, counted int
	outliers := StatisticalOutliers(arr)
	for _, val := range arr {
		if !contains(outliers, val) {
			total += val
			counted++
		}
	}
	if counted == 0 {
		return 0
	}
	resolved = total / counted
	return resolved
}

func resolveBool(arr []bool) bool {
	var resolved bool
	var maxCount int
	occuranceCount := make(map[bool]int)
	for _, val := range arr {
		_, ok := occuranceCount[val]
		if !ok {
			occuranceCount[val] = 1
		}
		occuranceCount[val]++
	}
	for key, val := range occuranceCount {
		if val > maxCount {
			maxCount = val
			resolved = key
		}
	}
	return resolved
}

func resolveString(arr []string) string {
	var resolved string
	maxCount := 0
	occuranceCount := make(map[string]int)
	for _, val := range arr {
		_, ok := occuranceCount[val]
		if !ok {
			occuranceCount[val] = 1
		}
		occuranceCount[val]++
	}
	for key, val := range occuranceCount {
		if val > maxCount {
			maxCount = val
			resolved = key
		}
	}
	return resolved
}

func contains(arr []int, val int) bool {
	for _, x := range arr {
		if x == val {
			return true
		}
	}
	return false
}

//IsUnanimous checks if data has significant disagreements - this usually means non-identical pieces of data unless we had something dealing in decimals/seconds

//DemocraticeCensus tries to find answers which are most common and picks them

//MeanCensus averages the answers

//MedianCensus takes the median of answers

//PruneOutliers takes outliers out of the data

//Calculation functions that operate on raw data

//Most of the below belong to match census functions

//DemocraticOutliers - Finds outliers which are not in the majority opinion i.e. are not the same value as most of the others
func DemocraticOutliers(data []float64) []float64 {
	//This will mark repeated values as outliers if there is another value more common than them
	outliers := make([]float64, 0, len(data))
	maxentries := 0
	counter := make(map[float64]int)
	for _, num := range data {
		_, ok := counter[num]
		if ok != true {
			counter[num] = 1
		} else {
			counter[num]++
		}
	}
	//Gets max number of occurances for any value
	for _, v := range counter {
		if v > maxentries {
			maxentries = v
		}
	}
	//Checks for outliers and appends them to list
	for k, v := range counter {
		if v < maxentries {
			outliers = append(outliers, k)
		}
	}
	return outliers
}

//StatisticalOutliers - Checks for outliers outside +/-1.5 of the interquartile range
func StatisticalOutliers(data []int) []int {
	data64 := make([]float64, len(data))
	for ind, val := range data {
		data64[ind] = float64(val)
	}
	outliers := make([]int, len(data64))
	mean := mean(data64)
	iqr := findIQR(data64)
	for ind, num := range data64 {
		if num > mean+iqr*1.5 || num < mean-iqr*1.5 {
			outliers[ind] = int(num)
		}
	}
	return outliers
}

func findIQR(data []float64) float64 {
	var half1 []float64
	var half2 []float64
	sort.Float64s(data)
	if len(data)%2 == 1 {
		half1 = data[:len(data)/2]
		half2 = data[len(data)/2+1:]
	} else {
		half1 = data[:len(data)/2]
		half2 = data[len(data)/2:]
	}
	return findMedian(half2) - findMedian(half1)
}

func findMedian(data []float64) float64 {
	var median float64
	if len(data) == 0 {
		return 0
	}
	sort.Float64s(data)
	if len(data)%2 == 1 {
		median = data[len(data)-1]
	} else {
		median = (data[len(data)/2] + data[len(data)/2-1]) / 2
	}
	return median
}

func standardDeviation(data []float64) float64 {
	mean := mean(data)
	var total float64 = 0
	for _, num := range data {
		total += math.Pow(num-mean, 2)
	}
	return math.Sqrt(total / float64(len(data)))
}

func mean(data []float64) float64 {
	var total float64 = 0
	for _, num := range data {
		total += num
	}
	return total / float64(len(data))
}

//getCurrentMatch indicates data from the last match in the scouting system, especially the match number
