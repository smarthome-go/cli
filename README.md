# homescript-cli

## REPL
Homescript-cli's original purpose is to represent a powerful command-line-interface for the Smarthome server.
Sample output of a *REPL* session can be found below
```
mik@mik-pc ~/g/s/g/M/smarthome (main) » homescript -i "http://cloud.box:8123" -u admin -p admin
Server: v0.0.15-beta:go1.18 on http://cloud.box:8123
Welcome to Homescript interactive v0.4.0-beta. CLI commands and comments start with #
admin@homescript> switch('s2', off) 
```

## Help Output
```
homescript-cli v2.0.0-beta-rc.1 : A command line interface for the smarthome server using homescript
A working and set-up Smarthome server instance is required.
For more information and usage documentation visit:

  The Homescript Programming Language:
  - https://github.com/MikMuellerDev/homescript

  The CLI Interface For Homescript:
  - https://github.com/MikMuellerDev/homescript-cli

  The Smarthome Server:
  - https://github.com/MikMuellerDev/smarthome

Usage:
  homescript [flags]
  homescript [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      REPL configuration
  debug       Smarthome Server Debug Info
  help        Help about any command
  pipe        Run Code via Stdin
  power       Power Summary
  run         Run a homescript file
  switches    List switches

Flags:
  -h, --help              help for homescript
  -i, --ip string         URL used for connecting to Smarthome (default "http://localhost")
  -p, --password string   the user's password used for connection
  -u, --username string   smarthome user used for connection
  -v, --verbose           verbose output
      --version           version for homescript

Use "homescript [command] --help" for more information about a command.
```

## Installation of v2.0.0-beta (form AMD64)
```
cd /tmp && wget https://github.com/MikMuellerDev/homescript-cli/releases/download/v2.0.0-beta/homescript_linux_amd64.tar.gz && tar -xvf homescript_linux_amd64.tar.gz && sudo mv homescript /usr/bin && rm -rf homescript_linux_amd64.tar.gz
```