# Smarthome-cli

[![Go Build](https://github.com/smarthome-go/cli/actions/workflows/go.yml/badge.svg)](https://github.com/smarthome-go/cli/actions/workflows/go.yml)
[![Typo Check](https://github.com/smarthome-go/cli/actions/workflows/typos.yml/badge.svg)](https://github.com/smarthome-go/cli/actions/workflows/typos.yml)

## REPL
Smarthome-cli's (formerly Homescript-cli) original purpose is to represent a powerful command-line-interface for the Smarthome server.
A sample output of a *REPL* session can be found below
```
Welcome to Homescript interactive v2.14.3. CLI commands and comments start with #
Server: v0.0.47:go1.18.3 on http://smarthome.box
admin@homescript>
```

## Help Output
```
homescript-cli v2.14.3 : A command line interface for the smarthome server using homescript
A working and set-up Smarthome server instance is required.
For more information and usage documentation visit:

  [1;32mThe Homescript Programming Language:[1;0m
  - https://github.com/smarthome-go/homescript

  [1;33mThe CLI Interface For Homescript:[1;0m
  - https://github.com/smarthome-go/cli

  [1;34mThe Smarthome Server:[1;0m
  - https://github.com/smarthome-go/smarthome

Usage:
  homescript [flags]
  homescript [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  config      CLI configuration
  debug       Server Debug Info
  help        Help about any command
  pipe        Run Code via Stdin
  power       Change switch power
  power       Power Summary
  run         Run a Homescript file
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

## Installation of v2.14.3 (for Linux/AMD64)

```
cd /tmp && wget https://github.com/smarthome-go/cli/releases/download/v2.14.3/homescript_linux_amd64.tar.gz && tar -xvf homescript_linux_amd64.tar.gz && sudo mv homescript /usr/bin && rm -rf homescript_linux_amd64.tar.gz
```

## Installation on Arch Linux

```bash
  yay -S smarthome-cli
# paru -S smarthome-cli
```
