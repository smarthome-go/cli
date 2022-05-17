package workspace

import (
	"errors"
	"fmt"
	"time"

	"github.com/smarthome-go/homescript/homescript/interpreter"
)

type Executor struct {
	ScriptName string
	Username   string
	Output     string
}

// Emulates printing to the console
// Instead, appends the provided message to the output of the executor
// Exists in order to return the script's output to the user
func (self *Executor) Print(args ...string) {
	var output string
	for _, arg := range args {
		self.Output += arg
		output += arg
	}
}

// Returns a boolean if the requested switch is on or off
// Returns an error if the provided switch does not exist
func (self *Executor) SwitchOn(switchId string) (bool, error) {
	return false, nil
}

// Changes the power state of an arbitrary switch
// Checks if the switch exists, if the user is allowed to interact with switches and if the user has the matching switch-permission
// If a check fails, an error is returned
func (self *Executor) Switch(switchId string, powerOn bool) error {
	return nil
}

// Sends a mode request to a given radiGo server
// TODO: implement this feature
func (self *Executor) Play(server string, mode string) error {
	return errors.New("The feature 'radiGo' is not yet implemented")
}

// Sends a notification to the user who issues this command
func (self *Executor) Notify(
	title string,
	description string,
	level interpreter.LogLevel,
) error {
	return nil
}

// Adds a new user to the system
// If the user already exists, an error is returned
func (self *Executor) AddUser(username string, password string, forename string, surname string) error {
	return nil
}

// Deletes a given user: checks whether it is okay to delete this user
func (self *Executor) DelUser(username string) error {
	return nil
}

// Adds an arbitrary permission to a given user
func (self *Executor) AddPerm(username string, permission string) error {
	return nil
}

// Removes an arbitrary permission from a given user
func (self *Executor) DelPerm(username string, permission string) error {
	return nil
}

// Adds a log entry to the internal logging system
func (self *Executor) Log(
	title string,
	description string,
	level interpreter.LogLevel,
) error {
	if level < 0 || level > 5 {
		return fmt.Errorf("Failed to add log event: invalid logging level <%d>: valid logging levels are 1, 2, 3, 4, or 5", level)
	}
	return nil
}

// Executes another Homescript based on its Id
func (self Executor) Exec(homescriptId string) (string, error) {
	// TODO: refer to string below
	return "", fmt.Errorf("NOT IMPLEMENTED: add dynamic linking via fetch")
}

// Returns the name of the user who is currently running the script
func (self *Executor) GetUser() string {
	return self.Username
}

// TODO: Will later be implemented, should return the weather as a human-readable string
func (self *Executor) GetWeather() (string, error) {
	return "rainy", nil
}

// TODO: Will later be implemented, should return the temperature in Celsius
func (self *Executor) GetTemperature() (int, error) {
	return 42, nil
}

// Returns the current time variables
func (self *Executor) GetDate() (int, int, int, int, int, int) {
	now := time.Now()
	return now.Year(), int(now.Month()), now.Day(), now.Hour(), now.Minute(), now.Second()
}
