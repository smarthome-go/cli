module github.com/smarthome-go/cli

go 1.18

require github.com/smarthome-go/sdk v0.14.0

require (
	github.com/briandowns/spinner v1.18.1
	github.com/chzyer/readline v1.5.0
	github.com/fatih/color v1.13.0
	github.com/howeyc/gopass v0.0.0-20210920133722-c8aef6fb66ef
	github.com/pelletier/go-toml v1.9.5
	github.com/rodaine/table v1.0.1
	github.com/sergi/go-diff v1.2.0
	github.com/spf13/cobra v1.4.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

require (
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.9 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.0.0-20220517005047-85d78b3ac167 // indirect
	golang.org/x/sys v0.0.0-20220310020820-b874c991c1a5 // indirect
	golang.org/x/term v0.0.0-20201126162022-7de9c90e9dd1 // indirect
)

replace github.com/smarthome-go/sdk => ../sdk
