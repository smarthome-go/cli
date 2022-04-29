package cmd

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/rodaine/table"

	"github.com/MikMuellerDev/homescript-cli/cmd/homescript"
	"github.com/MikMuellerDev/homescript-cli/cmd/log"
)

func listSwitches() {
	ch := make(chan bool)
	go homescript.StartSpinner("Loading Switches", &ch)
	getPersonalSwitches()
	ch <- true

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "Name", "Room", "Power", "Watts")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, switchItem := range Switches {
		powerIndicator := "off"
		if switchItem.PowerOn {
			powerIndicator = "on "
		}
		tbl.AddRow(switchItem.Id, switchItem.Name, switchItem.RoomId, powerIndicator, switchItem.Watts)
	}
	tbl.Print()
}

func debugInfo() {
	ch := make(chan bool)
	go homescript.StartSpinner("Loading Debug Information", &ch)
	if err := GetDebugInfo(); err != nil {
		log.Loge("Failed to fetch debug info: ", err.Error())
		return
	}
	ch <- true

	var output string

	output += color.New(color.FgGreen, color.Underline).Sprintf("Indicator%sValue%s\n", strings.Repeat(" ", 24), strings.Repeat(" ", 10))
	output += fmt.Sprintf("Smarthome Server Version: %s   v%s\n", strings.Repeat(" ", 30-len("Smarthome Server Version: ")), DebugInfo.ServerVersion)
	var databaseOnlineString = "\x1b[1;31mNO\x1b[1;0m"
	if DebugInfo.DatabaseOnline {
		databaseOnlineString = "\x1b[1;32mYES\x1b[1;0m"
	}
	output += fmt.Sprintf("Database Online: %s   %- 10s\n", strings.Repeat(" ", 30-len("Database Online: ")), databaseOnlineString)
	output += fmt.Sprintf("Compiled with: %s   %- 10s\n", strings.Repeat(" ", 30-len("Compiled with: ")), DebugInfo.GoVersion)
	output += fmt.Sprintf("CPU Cores: %s   %d\n", strings.Repeat(" ", 30-len("CPU Cores: ")), DebugInfo.CpuCores)
	output += fmt.Sprintf("Current Goroutines: %s   %d\n", strings.Repeat(" ", 30-len("Current Goroutines: ")), DebugInfo.Goroutines)
	output += fmt.Sprintf("Current Memory Usage: %s   %d\n", strings.Repeat(" ", 30-len("Current Memory Usage: ")), DebugInfo.MemoryUsage)
	output += fmt.Sprintf("Current Power Jobs: %s   %d\n", strings.Repeat(" ", 30-len("Current Power Jobs: ")), DebugInfo.PowerJobCount)
	output += fmt.Sprintf("Last Power Job Error Count: %s   %d", strings.Repeat(" ", 30-len("Last Power Job Error Count: ")), DebugInfo.PowerJobWithErrorCount)
	fmt.Println(output)
}
