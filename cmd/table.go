package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/rodaine/table"

	"github.com/smarthome-go/sdk"
)

func powerStats() {
	s := spinner.New(spinner.CharSets[11], 150*time.Millisecond)
	s.Suffix = " Loading power states"
	s.Start()
	switches, err := Connection.GetAllSwitches()
	if err != nil {
		switch err {
		case sdk.ErrConnFailed:
			s.FinalMSG = "Failed to fetch power states: network connection to Smarthome was interrupted.\n"
		case sdk.ErrServiceUnavailable:
			s.FinalMSG = "Failed to fetch power states: Smarthome is currently unavailable.\n"
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

	tbl := table.New("ID", "Name", "Power", "Watts")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	// Fill the table
	var total uint16 = 0
	var totalOn uint16 = 0
	var totalNotUsed uint16 = 0

	for _, switchItem := range Switches {
		var powerIndicator string
		if switchItem.PowerOn {
			powerIndicator = "on  *"
			totalOn += switchItem.Watts
		} else {
			powerIndicator = "off ."
			totalNotUsed += switchItem.Watts
		}
		total += switchItem.Watts
		tbl.AddRow(switchItem.Id, switchItem.Name, powerIndicator, switchItem.Watts)
	}

	// Prevent panic
	if total == 0 {
		total = 1
	}

	tbl.AddRow()
	tbl.AddRow("on ", "total (load)", "on   ", fmt.Sprintf("%-4d ~> %3d%s", totalOn, int((float32(totalOn)/float32(total))*100), `%`))
	tbl.AddRow("off", "total (free)", "off ", fmt.Sprintf("%-4d ~> %3d%s", totalNotUsed, int((float32(totalNotUsed)/float32(total))*100), `%`))
	tbl.AddRow("all", "total (all )", "all  ", fmt.Sprintf("%-4d => 100%s", total, `%`))
	tbl.Print()
}

func listSwitches() {
	s := spinner.New(spinner.CharSets[11], 150*time.Millisecond)
	s.Suffix = " Loading switches"
	s.Start()
	switches, err := Connection.GetPersonalSwitches()
	if err != nil {
		switch err {
		case sdk.ErrConnFailed:
			s.FinalMSG = "Failed to fetch switches: network connection to Smarthome was interrupted.\n"
		case sdk.ErrServiceUnavailable:
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
			powerIndicator = "on"
		}
		tbl.AddRow(switchItem.Id, switchItem.Name, switchItem.RoomId, powerIndicator, switchItem.Watts)
	}
	tbl.Print()
}
