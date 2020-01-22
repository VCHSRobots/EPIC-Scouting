# Scouting System Configuration

## General information

The configuration file `config.yaml` is written in [YAML](https://yaml.org/). (Obviously.)

It must be placed in the same directory as the scouting system's executable. Its name and location can not be changed without editing the source code.

## Options

### Mandatory

 - `DatabasePath`: The directory where database files are stored.
 - `LogPath`: The directory where log files are stored.

### Optional

 - `PortAPI`: The port for the server. `443` by default.
 - `TBAAuthKey`: A user's authentication key for [The Blue Alliance's](https://www.thebluealliance.com) API. Required to pull data from there.
 - `Verbosity`:
   - `-3`: only record `Fatal` log entries.
   - `-2`: only record `Error` or `Fatal` entries.
   - `-1`: only record `Warn`, `Error`, or `Fatal` entries. 
   - `0`: the default setting. Record `Info`, `Warn`, `Error`, and `Fatal` entries.
   - `1`: enable `Debug` messages and nanosecond timestamps for all entries.
 - `DatabaseBackupPath`: The location for database backups. If this is a web address or IP, the scouting server will attempt to use SFTP to upload the database files. `Null` by default.
 - `DatabaseBackupFrequency`: A positive integer; time expressed as seconds. For example, 86400 would be equivalent to once every 24 hours. Values less than or equal to `0` disable backups. `604800` by default.