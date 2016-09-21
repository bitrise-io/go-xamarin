package builder

import (
	"fmt"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/xamarin-builder/constants"
)

// MDToolCommandModel ...
type MDToolCommandModel struct {
	cmdSlice []string

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
		cmdSlice:    []string{constants.MDToolPath},
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

// Run ...
func (mdtool MDToolCommandModel) Run() error {
	if mdtool.target != "" {
		mdtool.cmdSlice = append(mdtool.cmdSlice, mdtool.target)
	}

	if mdtool.solutionPth != "" {
		mdtool.cmdSlice = append(mdtool.cmdSlice, mdtool.solutionPth)
	}

	config := ""
	if mdtool.configuration != "" {
		config = mdtool.configuration
	}

	if mdtool.platform != "" && mdtool.platform != "Any CPU" && mdtool.platform != "AnyCPU" {
		config += "|" + mdtool.platform
	}

	if config != "" {
		mdtool.cmdSlice = append(mdtool.cmdSlice, fmt.Sprintf("-c:%s", config))
	}

	if mdtool.projectName != "" {
		mdtool.cmdSlice = append(mdtool.cmdSlice, fmt.Sprintf("-p:%s", mdtool.projectName))
	}

	log.Info("=> %s", cmdex.PrintableCommandArgs(false, mdtool.cmdSlice))

	command, err := cmdex.NewCommandFromSlice(mdtool.cmdSlice)
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
