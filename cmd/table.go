package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/rodaine/table"

	"github.com/MikMuellerDev/smarthome_sdk"
)

func listSwitches() {
	s := spinner.New(spinner.CharSets[11], 150*time.Millisecond)
	s.Suffix = "Loading switches..."
	s.Start()
	switches, err := Connection.GetPersonalSwitches()
	if err != nil {
		switch err {
		case smarthome_sdk.ErrConnFailed:
			s.FinalMSG = "Failed to fetch switches: network connection to Smarthome was interrupted.\n"
		case smarthome_sdk.ErrServiceUnavailable:
			s.FinalMSG = "Failed to fetch switches: Smarthome is currently unavailable.\n"
		default:
			s.FinalMSG = fmt.Sprintf("An unexpected error occurred: %s\n", err.Error())
		}
		s.Stop()
		return
	}
	s.Stop()
	// Update switches for autosuggestion
	Switches = switches

	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("ID", "Name", "Room", "Power", "Watts")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	// Fill the table
	for _, switchItem := range Switches {
		powerIndicator := "off"
		if switchItem.PowerOn {
			powerIndicator = "on "
		}
		tbl.AddRow(switchItem.Id, switchItem.Name, switchItem.RoomId, powerIndicator, switchItem.Watts)
	}
	tbl.Print()
}

func printDebugInfo() {
	s := spinner.New(spinner.CharSets[11], 150*time.Millisecond)
	s.Suffix = "Loading switches..."
	s.Start()

	debugInfo, err := Connection.GetDebugInfo()
	if err != nil {
		switch err {
		case smarthome_sdk.ErrPermissionDenied:
			s.FinalMSG = "Debug information is not available for your user: you do not have the permission 'debug' which is required to view this information.\n"
		case smarthome_sdk.ErrConnFailed:
			s.FinalMSG = "Failed to fetch debug information: network connection to Smarthome was interrupted.\n"
		case smarthome_sdk.ErrServiceUnavailable:
			s.FinalMSG = "Failed to fetch debug information: Smarthome is currently unavailable.\n"
		default:
			s.FinalMSG = fmt.Sprintf("An unexpected error occurred: %s\n", err.Error())
		}
		s.Stop()
		return
	}
	s.Stop()

	// Generate output
	var output string
	output += color.New(color.FgGreen, color.Underline).Sprintf("Parameter%sValue%s\n", strings.Repeat(" ", 24), strings.Repeat(" ", 10))
	output += fmt.Sprintf("Smarthome Server Version: %s   v%s\n", strings.Repeat(" ", 30-len("Smarthome Server Version: ")), debugInfo.ServerVersion)
	var databaseOnlineString = "\x1b[1;31mNO\x1b[1;0m"
	if debugInfo.DatabaseOnline {
		databaseOnlineString = "\x1b[1;32mYES\x1b[1;0m"
	}
	output += fmt.Sprintf("Database Online: %s   %- 10s\n", strings.Repeat(" ", 30-len("Database Online: ")), databaseOnlineString)
	output += fmt.Sprintf("Compiled with: %s   %- 10s\n", strings.Repeat(" ", 30-len("Compiled with: ")), debugInfo.GoVersion)
	output += fmt.Sprintf("CPU Cores: %s   %d\n", strings.Repeat(" ", 30-len("CPU Cores: ")), debugInfo.CpuCores)
	output += fmt.Sprintf("Current Goroutines: %s   %d\n", strings.Repeat(" ", 30-len("Current Goroutines: ")), debugInfo.Goroutines)
	output += fmt.Sprintf("Current Memory Usage: %s   %d\n", strings.Repeat(" ", 30-len("Current Memory Usage: ")), debugInfo.MemoryUsage)
	output += fmt.Sprintf("Current Power Jobs: %s   %d\n", strings.Repeat(" ", 30-len("Current Power Jobs: ")), debugInfo.PowerJobCount)
	output += fmt.Sprintf("Last Power Job Error Count: %s   %d", strings.Repeat(" ", 30-len("Last Power Job Error Count: ")), debugInfo.PowerJobWithErrorCount)
	fmt.Println(output)
}
