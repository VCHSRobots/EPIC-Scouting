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
	fmt.Println(overall)
	return overall
}

//TeamAuto gets a team's autonomous rating
func TeamAuto(teamNum int, campaignid string) int {
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
			breakdown[0] = 1
		}
		breakdown[1] += match.AutoBackBalls
		breakdown[2] += match.AutoHighBalls
		breakdown[3] += match.AutoLowBalls
		breakdown[4] += match.AutoShots
		breakdown[5] += match.AutoPickups
		if match.AutoShots-match.AutoLowBalls != 0 {
			breakdown[6] += (match.AutoBackBalls + match.AutoHighBalls) / (match.AutoShots - match.AutoLowBalls)
		} else {
			breakdown[6] = 0
		}
		//total auto points
		breakdown[7] += breakdown[0]*15 + match.AutoBackBalls*6 + match.AutoHighBalls*4 + match.AutoLowBalls*2
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
		} else {
			breakdown[4] = 0
		}
		breakdown[5] = match.LowFuel + match.HighFuel + match.BackFuel
		//total auto points
		breakdown[6] += match.LowFuel*1 + match.HighFuel*2 + match.BackFuel*3
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
			breakdown[0] = 2
		} else if match.Climbed == "platform" {
			breakdown[0] = 1
		} else {
			breakdown[0] = 0
		}
		breakdown[1] += match.ClimbTime
		if match.Balanced {
			breakdown[2] = 1
		}
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
			breakdown[3] = 2
		} else if match.Card == "yellow" {
			breakdown[3] = 1
		} else {
			breakdown[3] = 0
		}
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

/*Match Census Functions determine the weight of contradictary data on the same match and return a score useable for the system*/
//Below are differing census methods. They may or may not be used.

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
func StatisticalOutliers(data []float64) []float64 {
	outliers := make([]float64, 0, len(data))
	mean := mean(data)
	iqr := findIQR(data)
	for _, num := range data {
		if num > mean+iqr*1.5 || num < mean-iqr*1.5 {
			outliers = append(outliers, num)
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
