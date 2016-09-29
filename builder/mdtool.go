package builder

import (
	"fmt"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-tools/go-xamarin/constants"
)

// MDToolCommandModel ...
type MDToolCommandModel struct {
	buildTool string

	solutionPth   string
	projectName   string
	configuration string
	platform      string
	target        string
}

// NewMDToolCommand ...
func NewMDToolCommand(solutionPth string) *MDToolCommandModel {
	return &MDToolCommandModel{
		solutionPth: solutionPth,
		buildTool:   constants.MDToolPath,
	}
}

// SetTarget ...
func (mdtool *MDToolCommandModel) SetTarget(target string) *MDToolCommandModel {
	mdtool.target = target
	return mdtool
}

// SetConfiguration ...
func (mdtool *MDToolCommandModel) SetConfiguration(configuration string) *MDToolCommandModel {
	mdtool.configuration = configuration
	return mdtool
}

// SetPlatform ...
func (mdtool *MDToolCommandModel) SetPlatform(platform string) *MDToolCommandModel {
	mdtool.platform = platform
	return mdtool
}

// SetProjectName ...
func (mdtool *MDToolCommandModel) SetProjectName(projectName string) *MDToolCommandModel {
	mdtool.projectName = projectName
	return mdtool
}

func (mdtool MDToolCommandModel) buildCommandSlice() []string {
	cmdSlice := []string{mdtool.buildTool}

	if mdtool.target != "" {
		cmdSlice = append(cmdSlice, mdtool.target)
	}

	if mdtool.solutionPth != "" {
		cmdSlice = append(cmdSlice, mdtool.solutionPth)
	}

	config := ""
	if mdtool.configuration != "" {
		config = mdtool.configuration
	}

	if mdtool.platform != "" && mdtool.platform != "Any CPU" && mdtool.platform != "AnyCPU" {
		config += "|" + mdtool.platform
	}

	if config != "" {
		cmdSlice = append(cmdSlice, fmt.Sprintf("-c:%s", config))
	}

	if mdtool.projectName != "" {
		cmdSlice = append(cmdSlice, fmt.Sprintf("-p:%s", mdtool.projectName))
	}

	return cmdSlice
}

// PrintableCommand ...
func (mdtool MDToolCommandModel) PrintableCommand() string {
	cmdSlice := mdtool.buildCommandSlice()

	return cmdex.PrintableCommandArgs(false, cmdSlice)
}

// Run ...
func (mdtool MDToolCommandModel) Run() error {
	cmdSlice := mdtool.buildCommandSlice()

	command, err := cmdex.NewCommandFromSlice(cmdSlice)
	if err != nil {
		return err
	}

	/*
		command.SetStdout(os.Stdout)
		command.SetStderr(os.Stderr)

		return command.Run()
	*/

	return runCommandInDiagnosticMode(*command, "Loading projects", false)
}
