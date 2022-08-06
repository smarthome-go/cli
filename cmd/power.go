package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func createCmdPower() *cobra.Command {
	// Parent power command
	cmdPower := &cobra.Command{
		Use:   "power",
		Short: "Power Subcommand",
		Long:  "Power subcommand for interacting with switches and viewing statistics",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				panic(err.Error())
			}
		},
	}

	// Power on
	cmdPowerOn := &cobra.Command{
		Use:   "on [switch-id] ",
		Short: "Activate Switch",
		Long:  "Activate an arbitrary switch",
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
			// Initialize Smarthome connection
			InitConn()
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := Connection.SetPower(args[0], true); err != nil {
				fmt.Printf("Could not activate switch.\nError: %s\n", err.Error())
				os.Exit(1)
			}
			fmt.Printf("Successfully turned switch %s on.\n", args[0])
			os.Exit(0)
		},
	}

	// Power off
	cmdPowerOff := &cobra.Command{
		Use:   "off [switch-id] ",
		Short: "Deactivate Switch",
		Long:  "Deactivate an arbitrary switch",
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
			// Initialize Smarthome connection
			InitConn()
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := Connection.SetPower(args[0], false); err != nil {
				fmt.Printf("Could not deactivate switch.\nError: %s\n", err.Error())
				os.Exit(1)
			}
			fmt.Printf("Successfully turned switch %s off.\n", args[0])
			os.Exit(0)
		},
	}

	// Power toggle
	cmdPowerToggle := &cobra.Command{
		Use:   "toggle [switch-id] ",
		Short: "Toggle Switch Power",
		Long:  "Toggle the power-state of an arbitrary switch",
		Args:  cobra.ExactArgs(1),
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
			// Initialize Smarthome connection
			InitConn()
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Get current power state
			if Verbose {
				fmt.Println("Retrieving personal switches in order to get current power state...")
			}
			switches, err := Connection.GetPersonalSwitches()
			if err != nil {
				fmt.Printf("Could not toggle switch: Failed to retrieve current power state.\nError: %s\n", err.Error())
				os.Exit(1)
			}
			currentPower := false
			for _, sw := range switches {
				if sw.Id == args[0] {
					currentPower = sw.PowerOn
				}
			}
			if Verbose {
				fmt.Println("Current power state has been retrieved successfully")
			}

			// Change the power state of the switch
			if err := Connection.SetPower(args[0], !currentPower); err != nil {
				fmt.Printf("Could not toggle switch.\nError: %s\n", err.Error())
				os.Exit(1)
			}
			fmt.Printf("Successfully toggled switch %s.\n", args[0])
			os.Exit(0)
		},
	}

	cmdPowerSummary := &cobra.Command{
		Use:   "draw",
		Short: "Power Draw & States",
		Long:  "A compact overview of estimated power usage and states",
		Args:  cobra.NoArgs,
		PreRun: func(cmd *cobra.Command, args []string) {
			readConfigFile()
		},
		Run: func(cmd *cobra.Command, args []string) {
			InitConn()
			powerStats()
		},
	}

	cmdPower.AddCommand(cmdPowerOn)
	cmdPower.AddCommand(cmdPowerOff)
	cmdPower.AddCommand(cmdPowerToggle)
	cmdPower.AddCommand(cmdPowerSummary)

	return cmdPower
}
