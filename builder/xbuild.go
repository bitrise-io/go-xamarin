package builder

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-tools/go-xamarin/constants"
)

// BuildCommand ...
type BuildCommand interface {
	PrintableCommand() string
}

// XbuildCommandModel ...
type XbuildCommandModel struct {
	buildTool string

	projectPth    string // can be solution or project path
	configuration string
	platform      string
	target        string

	buildIpa       bool
	archiveOnBuild bool

	exportOutput func() (OutputMap, error)
}

// NewXbuildCommand ...
func NewXbuildCommand(projectPth string) *XbuildCommandModel {
	return &XbuildCommandModel{
		projectPth: projectPth,
		buildTool:  constants.XbuildPath,
	}
}

// SetTarget ...
func (xbuild *XbuildCommandModel) SetTarget(target string) *XbuildCommandModel {
	xbuild.target = target
	return xbuild
}

// SetConfiguration ...
func (xbuild *XbuildCommandModel) SetConfiguration(configuration string) *XbuildCommandModel {
	xbuild.configuration = configuration
	return xbuild
}

// SetPlatform ...
func (xbuild *XbuildCommandModel) SetPlatform(platform string) *XbuildCommandModel {
	xbuild.platform = platform
	return xbuild
}

// SetBuildIpa ...
func (xbuild *XbuildCommandModel) SetBuildIpa() *XbuildCommandModel {
	xbuild.buildIpa = true
	return xbuild
}

// SetArchiveOnBuild ...
func (xbuild *XbuildCommandModel) SetArchiveOnBuild() *XbuildCommandModel {
	xbuild.archiveOnBuild = true
	return xbuild
}

func (xbuild XbuildCommandModel) buildCommandSlice() []string {
	cmdSlice := []string{xbuild.buildTool}

	if xbuild.projectPth != "" {
		cmdSlice = append(cmdSlice, xbuild.projectPth)
	}

	if xbuild.target != "" {
		cmdSlice = append(cmdSlice, fmt.Sprintf("/target:%s", xbuild.target))
	}

	if xbuild.configuration != "" {
		cmdSlice = append(cmdSlice, fmt.Sprintf("/p:Configuration=%s", xbuild.configuration))
	}

	if xbuild.platform != "" {
		cmdSlice = append(cmdSlice, fmt.Sprintf("/p:Platform=%s", xbuild.platform))
	}

	if xbuild.archiveOnBuild {
		cmdSlice = append(cmdSlice, "/p:ArchiveOnBuild=true")
	}

	if xbuild.buildIpa {
		cmdSlice = append(cmdSlice, "/p:BuildIpa=true")
	}

	return append(cmdSlice, "/verbosity:minimal", "/nologo")
}

// PrintableCommand ...
func (xbuild XbuildCommandModel) PrintableCommand() string {
	cmdSlice := xbuild.buildCommandSlice()

	return cmdex.PrintableCommandArgs(false, cmdSlice)
}

// Run ...
func (xbuild XbuildCommandModel) Run() error {
	cmdSlice := xbuild.buildCommandSlice()

	command, err := cmdex.NewCommandFromSlice(cmdSlice)
	if err != nil {
		return err
	}

	command.SetStdout(os.Stdout)
	command.SetStderr(os.Stderr)

	return command.Run()
}
