package scouting

//QuerySchedule gets called upon GET request from route/scout. It checks the SCHEDULE database for the next match to be scouted. If it is empty, populate the match's entry with teamids of the participants and records of which scouters each team

//Scheduler gets data for teams being scouted in the next match

//GetTeamInfo gets number of scouters an other information from a specific competitior team

//GetTeamMatch gets which match the team is on, updated by the team admin

//GetScouterMatch gets which match a specific scouter is on

//AssignNewUser assigns a user who is ready to scout to the next match

//ShouldAdvanceTeamMatch

//NextMatch gets data on which match is next to be scouted

//openSchedule deserializes the schedule from the database

//writeSchedule serializes the schedule and writes it to the database

//matchParticpants returns match participants

//pickScoutedTeam picks which team to scout based on priority and what teams are already being scouted
