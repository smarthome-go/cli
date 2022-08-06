## Changelog

### Bugfixes
- Refactored code in order to fix typos and remove `ioutil` deprecation
- Fixed broken `MDIcon` syncing when using `smarthome-cli ws push / pull`
- Fixed conflicting `power` command

### Power Subcommand
- `smarthome-cli power` is now a subcommand which removes the preexisting conflict

#### Help Output
```
Power subcommand for interacting with switches and viewing statistics

Usage:
  smarthome-cli power [flags]
  smarthome-cli power [command]

Available Commands:
  draw        Power Draw & States
  off         Deactivate Switch
  on          Activate Switch
  toggle      Toggle Switch Power

Flags:
  -h, --help   help for power

Global Flags:
  -i, --ip string         URL used for connecting to Smarthome (default "http://localhost")
  -p, --password string   the user's password used for connection
  -u, --username string   smarthome user used for connection
  -v, --verbose           verbose output

Use "smarthome-cli power [command] --help" for more information about a command.
```
