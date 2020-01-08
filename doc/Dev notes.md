# EPIC Scouting development notes

## Deployment

 - Install Go
 - Install SQLite

### Development workflow

Development branch source > Built application > Development server
Production branch source > Built application > Production server

Once we complete tests for the development branch, push it to the production branch. Changes should be reflected automatically.

### Backend

 - Go + Gin webserver
   - Automated pulling from the production repository with webhooks and restarting
   - Automated starting of the server on host bootup via systemd service.
   - Easy entry of scouting data
   - Serving stored data
   - Encrypt all communication between client and server, and ensure identity of clients.
 - [SQLite](https://www.sqlite.org/index.html) database and [Golang drivers](https://github.com/mattn/go-sqlite3)
 - Compare data with other teams using our service. Coopertition!

## Frontend

 - HTML, CSS, JS
 - API for applications to interface with the backend (GET and POST).

### Users

 - Supervisors
   - View live match data + Blue Alliance sourced data
     - Raw database table dumps
     - Auto-generated charts and graphs for power rankings, skillsets, etc
   - Export database to CSV
   - Assign users tasks
   - View live user activity
   - Create new "game" setups
     - Lists of questions, etc
   - View server stats (uptime, last updated, version, database storage usage)

## Post-Production

 - Make a nice whitepaper with LaTeX :)

### Databases / API

#### Users

Username:
 - Global permission level: User | SysAdmin
 - Password hash: #########

#### Teams

Team: 4415
  - Users:
    - Alice:
      - Type: Owner | Supervisor | Scout
    - ...
  - Seasons:
    - 2020:
      - Active: True | False # If the campaign is archived or not.
      - Games: # A list of matches and the teams that will be competing in each match, along with their location. Likely imported from TBA.
        - Finals:
          - Active: True | False # If the game is active or not.
          - Location: 5555 Example St, Somewhere, USA
          - Time: 2077-01-01 00:00
          - Schedule: # The expected schedule of the game. Can be overridden by user in the UI if things change.
            - MatchNumber, MatchTime, TeamNumber, AllianceColor, FieldLocation
          - Results:
            - MatchNumber:
              - TeamNumber:
                - ObjectiveOne: Foo
                - InfractionOne: Bar
      - Teams:
        - TeamNumber: 0254
          - Details...

### Interfaces

 - Terminal
 - API / Web app
 - Platform-native applications
   - Internet connection
   - Scanning QR codes to an offline server

## IMPORTANT

DEV PLAN:
 - Build API for POST and GET-ting data, and data crunching. I am still uncertain of best input methods. Scanning QR codes from web app or native app? Live communication over cell network?
   - Consider [OpenAPI](https://swagger.io/tools/swagger-codegen/) for client-side API interaction
 - Generally concerned about connectivity issues at game.

### TBA API:
 - https://www.thebluealliance.com/apidocs/webhooks
 - https://www.thebluealliance.com/apidocs/v3