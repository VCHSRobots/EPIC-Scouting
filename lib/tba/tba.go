/*
Package tba manages communication with the The Blue Alliance's API.
See https://www.thebluealliance.com/apidocs for more information.
*/
package tba

import (
	"fmt"
	"net/http"
)

// Putting authentication key variable here until it can get connected to the config file.
// All request functions return a Response struct. Docs can be found at https://golang.org/src/net/http/response.go.
const tba string = "https://www.thebluealliance.com/api/v3"
const keyfield string = "X-TBA-Auth-Key="
const keystring string = "?X-TBA-Auth-Key=eSMjxo253BoTGFoaeZteq7wF1pGLGnZw24aaHxfvfsvF7VvNaTLOf7ZlvhbbJQxs"

var enabled bool

func httpGet(url, dir, querystring string) *http.Response {
	resp, err := http.Get(fmt.Sprintf("%s%s%s", tba, dir, querystring))
	// Return nil if GET fails.
	if err == nil {
		return resp
	}
	return nil
}

//KeyIsWorking checks if error code status is okay, which indicates the auth key is valid
func KeyIsWorking() bool {
	if GetStatus().Status == "200 OK" {
		return true
	}
	return false
}

//GetStatus - Gets TBA API Status
func GetStatus() *http.Response {
	dir := "/status"
	resp := httpGet(tba, dir, keystring)
	return resp
}

//GetMatch gets match data including participating teams and results
func GetMatch(match string) *http.Response {
	dir := fmt.Sprintf("/match/%s", match)
	resp := httpGet(tba, dir, keystring)
	return resp
}

//GetEventTeams gets which teams are at an event
func GetEventTeams(event string) *http.Response {
	dir := fmt.Sprintf("/event/%s/teams", event)
	resp := httpGet(tba, dir, keystring)
	return resp
}

//GetEventMatches gets a list of matches at an event
func GetEventMatches(event string) *http.Response {
	dir := fmt.Sprintf("/event/%s/matches", event)
	resp := httpGet(tba, dir, keystring)
	return resp
}

//GetTeamMatches gets all matches a team was involved in
func GetTeamMatches(team, year string) *http.Response {
	dir := fmt.Sprintf("/team/%s/matches%s", team, year)
	resp := httpGet(tba, dir, keystring)
	return resp
}

//GetTeamEvents gets all events a team was involved in
func GetTeamEvents(team string) *http.Response {
	dir := fmt.Sprintf("/team/%s/events", team)
	resp := httpGet(tba, dir, keystring)
	return resp
}
