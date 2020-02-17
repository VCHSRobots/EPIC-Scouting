package scouting

import "math"

//Scheduler gets data for teams being scouted in the next match

//GetTeamInfo gets number of scouters an other information from a specific competitior team

//GetTeamMatch gets which match the team is on, updated by the team admin

//GetScouterMatch gets which match a specific scouter is on

//AssignNewUser assigns a user who is ready to scout to the next match
func AssignNewUser(userid string) string {
	priorities := getPriorityMap("match")
	totalPriority := 0
	scoutersWanted := make(map[string]int, 0)
	toScout := ""
	for _, priority := range priorities {
		totalPriority += priority
	}
	for team, priority := range priorities {
		scoutersWanted[team] = int(math.Round(float64(priority) / float64(totalPriority)))
	}
	topPriority := 0
	currentScouters := 0
	for team, scouters := range scoutersWanted {
		currentScouters = 0 //TODO get how many scouters each team has
		if (scouters - currentScouters) > topPriority {
			topPriority = scouters
			toScout = team
		}
	}
	//Add userid to scout team in toScout
	return toScout
}

//TODO getPrioriyMap gets a map of the teams in the next match along with their correlating priorities
func getPriorityMap(matchid string) map[string]int {
	pmap := make(map[string]int, 6)
	return pmap
}

//UpdateMatch updates the match to be scouted

//NextMatch gets data on which match is next to be scouted

//openSchedule deserializes the schedule from the database

//writeSchedule serializes the schedule and writes it to the database

//matchParticpants returns match participants

//pickScoutedTeam picks which team to scout based on priority and what teams are already being scouted
