## 2020-02-18 — Prerelease 0.5 (Minimum viable product)

### Backend

 - [X] Secure a domain name and production server
 - [X] Logging library
   - [X] Leveled logging
   - [X] Thread-safe
   - [X] Library-specific loggers
 - [ ] Database interaction library
 - [ ] "Calc" library
 - [ ] Authentication middleware
 - [X] Gin router
 - [ ] Task scheduler daemon
 - [X] GZIP served page data to save bandwidth

### Frontend

 - [ ] Sysadmin overview page
 - [ ] Login / Register pages
 - [ ] Team admin overview page
 - [ ] Campaign / event / match selection screen for teams
 - [ ] Simple list-based data entry scouting page
 - [ ] User profiles
 - [ ] User dashboard
 - [ ] Team-joining

## 2020-03-03 — Prerelease 0.75

### Backend

 - [ ] HTTPS cert and encryption as a default
 - [ ] Finish up most pressing TODO comments in code.

### Frontend

 - [ ] CSS and beautification pass
 - [ ] Map-based data entry
 - [ ] Finalize adaptive design to target mobile devices
 - [ ] "About this software" page, include link to git repository and documentation

### QA testing

 - [ ] Non-programming team members practice scouting with previous competition's footage; compare data quality against ground-truth match results
 - [ ] Fine-tune match prediction algorithms against actual results
 - [ ] Experiment with rate-limited and high-latency cellular connections to simulate worst-case scenarios; further minimize required data

## 2020-03-09 — Release 1.0

 - Final necessary polishing before competition

## FUTURE — Release 1.5+

### Medium Priority

 - [ ] Finish any remaining TODOs
 - [ ] Split the front-end and back-end API into two separate servers
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
 - [ ] Allow for creating custom campaigns and match conditions via YAML and SVGs
 - [ ] Allow for custom result weighting algorithms
 - [ ] Switch to a database system that allows for nested tables
 - [ ] Allow reloading of configuration without rebooting by sysadmins

### Low Priority

- [ ] Final code review and beatification pass
- [ ] Complete LaTeX whitepaper and documentation
   - [ ] Render to HTML and PDF; host both on production server website
 - [ ] Advertise service at scrimmages, Chief Delphi, Reddit, Discord, etc (business team)
 - [ ] Create introductory video showcase / tutorial (business team)
 - [ ] Create a fancy 1337 h4ck3r console-based client (complete with flashing green text)