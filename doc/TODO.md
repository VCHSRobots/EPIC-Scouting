## 2020-01: January

### Backend

 - [ ] Automatic startup / reboot daemon (independent BASH script)
 - [ ] Automatic update and reboot from production git branch
 - [ ] Automatic database backups
 - [ ] Secure a domain name and production server
 - [X] Logging module
   - [X] Leveled logging
   - [X] Thread-safe
   - [X] Module-specific loggers
   - [ ] Integrate Gin log messages under their appropriate logLevels with *lumberjack*
 - [ ] Database interaction module
 - [ ] Authentication middleware
 - [ ] Gin router
 - [ ] Task scheduler daemon
 - [ ] *Low priority:* GZIP served page data to save bandwidth
 - [ ] *Low priority:* Allow reloading of configuration without rebooting by sysadmins
 - [ ] Feature-complete API
   - [ ] User registration
   - [ ] Team registration
   - [ ] Sysadmin page / console
     - [ ] Network throughput
     - [ ] Storage usage
     - [ ] Latest and upcoming backup times
     - [ ] Number of teams / users
     - [ ] Log output
   - [ ] TBA integration
   - [ ] POST match results
   - [ ] GET pre-rendered data (graphs, CSV database dumps, weighted skills)
   - [ ] Allow uploading YAML descriptions of match criteria for describing game objectives
   - [ ] Edit per-team objective weights on the fly
   - [ ] Log dumps for sysadmins

### Frontend

 - [ ] Sysadmin overview page
 - [ ] Team admin overview page
 - [ ] Campaign selection screen for teams
 - [ ] Simple list-based data entry

### Miscellaneous

 - [ ] Continue adding docstrings and comments to code
 - [ ] Complete API documentation
 - [ ] Complete initial `README.md`

## 2020-02: February

### Backend

 - [ ] HTTPS cert and encryption as a default
 - [ ] Lorem Ipsum
 - [ ] Foo bar baz

### Frontend

 - [ ] Offline data collection
   - [ ] Create QR or other visual codes to display end-of-match results by team members if the page disconnects (or is otherwise requested to do so)
   - [ ] Allow team admins to upload images of these QR codes to later transfer to the server
 - [ ] CSS and beautification pass
 - [ ] Map-based data entry
 - [ ] Finalize adaptive design to target mobile devices
 - [ ] *Very low priority:* Create a fancy 1337 h4ck3r console-based client (complete with flashing green text)
 - [ ] "About this software" page, include link to git repository
   - [ ] Donation button or link to team Patreon / whatever

### QA testing

 - [ ] Non-programming team members practice scouting with previous years footage; compare data quality against ground-truth match results
 - [ ] Fine-tune match prediction algorithms against actual results
 - [ ] Experiment with rate-limited and high-latency cellular connections to simulate worst-case scenarios; further minimize required data

### Miscellaneous

 - [ ] Complete LaTeX whitepaper and final documentation
   - [ ] Render to HTML and PDF; host both on production server website

## 2020-03: March

### Miscellaneous

 - [ ] 1.0 release!
 - [ ] Advertise service at scrimmages, Chief Delphi, Reddit, Discord, etc (business team)
 - [ ] Create introductory video showcase / tutorial (business team)

## 2020-03-??: Regionals

 - [ ] Regional 1
 - [ ] Regional 2

## 2020-04-15 - 18: FIRST Championship at Houston

 - [ ] You have lost [The Game](https://en.wikipedia.org/wiki/The_Game_(mind_game))
