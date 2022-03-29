package cmd

import (
	"fmt"

	"github.com/MikMuellerDev/homescript-cli/cmd/log"
)

func listSwitches() {
	// log.Logn("\u2502Id          \u2502 Label	            \u2502   Room     \u2502 Watts")
	// TODO: add power indicator
	log.Logn("\x1b[90m(Id, Label, Room, Watts)\x1b[0m")
	for _, switchItem := range Switches {
		log.Logn(fmt.Sprintf("\u2502%-10s  \u2502 %-20s \u2502 %10s \u2502 %d", switchItem.Id, switchItem.Name, switchItem.RoomId, switchItem.Watts))
	}
}