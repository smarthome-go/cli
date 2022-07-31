package cmd

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
	"github.com/fatih/color"
	"github.com/rodaine/table"
	"github.com/smarthome-go/sdk"
)

// Prints the server's debugging information
func printDebugInfo() {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Suffix = " Loading debug information"
	s.Start()

	debugInfo, err := Connection.GetDebugInfo()
	if err != nil {
		switch err {
		case sdk.ErrPermissionDenied:
			s.FinalMSG = "Debug information is not available for your user: you lack the permission 'debug' which is required to obtain this information.\n"
		case sdk.ErrConnFailed:
			s.FinalMSG = "Failed to fetch debug information: network connection to Smarthome was interrupted.\n"
		case sdk.ErrServiceUnavailable:
			s.FinalMSG = "Failed to fetch debug information: Smarthome is currently unavailable.\n"
		default:
			s.FinalMSG = fmt.Sprintf("An unexpected error occurred: %s\n", err.Error())
		}
		s.Stop()
		return
	}
	s.Stop()

	// Generate output
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	// columnFmt := color.New(color.FgWhite).SprintfFunc()

	tbl := table.New("Parameter", "Value")
	tbl.WithHeaderFormatter(headerFmt) //.WithFirstColumnFormatter(columnFmt)

	// Smarthome version information
	tbl.AddRow("Server version", debugInfo.ServerVersion)
	tbl.AddRow("Server GO version", debugInfo.GoVersion)

	// Performance statistics
	tbl.AddRow("CPU cores", debugInfo.CpuCores)
	tbl.AddRow("Used MEM", debugInfo.MemoryUsage)
	tbl.AddRow("Active Goroutines", debugInfo.Goroutines)

	// Power statistics
	tbl.AddRow("Power jobs", debugInfo.PowerJobCount)
	tbl.AddRow("Power jobs (FAILED)", debugInfo.PowerJobWithErrorCount)

	// Database status
	onlineStr := "online"
	if !debugInfo.DatabaseOnline {
		onlineStr = "OFFLINE"
	}
	tbl.AddRow("DB status", onlineStr)
	tbl.AddRow("DB conns (open)", debugInfo.DatabaseStats.OpenConnections)
	tbl.AddRow("DB conns (used)", debugInfo.DatabaseStats.InUse)
	tbl.AddRow("DB conns (idle)", debugInfo.DatabaseStats.Idle)

	// Hardware node information
	tbl.AddRow("HW nodes (total  )", debugInfo.HardwareNodesCount)
	tbl.AddRow("HW nodes (online )", debugInfo.HardwareNodesOnline)
	tbl.AddRow("HW nodes (enabled)", debugInfo.HardwareNodesEnabled)
	tbl.AddRow("HMS jobs", debugInfo.HomescriptJobCount)

	tbl.Print()

	// Also print the Hardware nodes
	fmt.Println()
	printHWnodes(debugInfo)
}

func printHWnodes(debugInfo sdk.DebugInfoData) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	tbl := table.New("URL", "Name", "Enabled", "Online")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

	for _, node := range debugInfo.HardwareNodes {
		enabledStr := "yes *"
		if !node.Enabled {
			enabledStr = "no  ."
		}
		onlineStr := "yes *"
		if !node.Online {
			onlineStr = "no  ."
		}
		tbl.AddRow(node.Url, node.Name, enabledStr, onlineStr)
	}

	tbl.Print()
}
