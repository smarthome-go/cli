# homescript-cli

## REPL
Homescript-cli's original purpose is to represent a powerful command-line-interface for the Smarthome server.
Sample output of a *REPL* session can be found below
```
Welcome to Homescript interactive v2.2.0-beta. CLI commands and comments start with #
Server: v0.0.26:go1.18.1 on http://cloud.box:8123
admin@homescript> 
```

## Help Output
```
homescript-cli v2.2.0-beta : A command line interface for the smarthome server using homescript
A working and set-up Smarthome server instance is required.
For more information and usage documentation visit:

  The Homescript Programming Language:
  - https://github.com/smarthome-go/homescript

  The CLI Interface For Homescript:
  - https://github.com/smarthome-go/cli

  The Smarthome Server:
  - https://github.com/smarthome-go/smarthome

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
  ws          workspace

Flags:
  -h, --help              help for homescript
  -i, --ip string         URL used for connecting to Smarthome (default "http://localhost")
  -p, --password string   the user's password used for connection
  -u, --username string   smarthome user used for connection
  -v, --verbose           verbose output
      --version           version for homescript

Use "homescript [command] --help" for more information about a command.
```

## Installation of v2.2.0-beta (for Linux/AMD64)
```
cd /tmp && wget https://github.com/smarthome-go/cli/releases/download/v2.2.0-beta/homescript_linux_amd64.tar.gz && tar -xvf homescript_linux_amd64.tar.gz && sudo mv homescript /usr/bin && rm -rf homescript_linux_amd64.tar.gz
```
