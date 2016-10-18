package nunit

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xamarin/constants"
)

// Model ...
type Model struct {
	nunitConsolePth string
	projectPth      string
	config          string

	customOptions []string
}

// New ...
func New(nunitConsolePth string) (*Model, error) {
	absNunitConsolePth, err := pathutil.AbsPath(nunitConsolePth)
	if err != nil {
		return nil, fmt.Errorf("Failed to expand path (%s), error: %s", nunitConsolePth, err)
	}

	return &Model{nunitConsolePth: absNunitConsolePth}, nil
}

// SetProjectPth ...
func (nunitConsole *Model) SetProjectPth(projectPth string) *Model {
	nunitConsole.projectPth = projectPth
	return nunitConsole
}

// SetConfig ...
func (nunitConsole *Model) SetConfig(config string) *Model {
	nunitConsole.config = config
	return nunitConsole
}

// SetCustomOptions ...
func (nunitConsole *Model) SetCustomOptions(options ...string) {
	nunitConsole.customOptions = options
}

func (nunitConsole *Model) commandSlice() []string {
	cmdSlice := []string{constants.MonoPath}
	cmdSlice = append(cmdSlice, nunitConsole.nunitConsolePth)
	cmdSlice = append(cmdSlice, nunitConsole.projectPth)
	cmdSlice = append(cmdSlice, fmt.Sprintf("/config:%s", nunitConsole.config))
	cmdSlice = append(cmdSlice, nunitConsole.customOptions...)
	return cmdSlice
}

// PrintableCommand ...
func (nunitConsole Model) PrintableCommand() string {
	cmdSlice := nunitConsole.commandSlice()

	return cmdex.PrintableCommandArgs(true, cmdSlice)
}

// Run ...
func (nunitConsole Model) Run() error {
	cmdSlice := nunitConsole.commandSlice()

	command, err := cmdex.NewCommandFromSlice(cmdSlice)
	if err != nil {
		return err
	}

	command.SetStdout(os.Stdout)
	command.SetStderr(os.Stderr)

	return command.Run()
}
