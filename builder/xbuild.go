package builder

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/xamarin-builder/constants"
)

// XbuildCommandModel ...
type XbuildCommandModel struct {
	cmdSlice []string

	projectPth    string // can be solution or project path
	configuration string
	platform      string
	target        string

	buildIpa bool

	exportOutput func() (OutputMap, error)
}

// NewXbuildCommand ...
func NewXbuildCommand(projectPth string) *XbuildCommandModel {
	return &XbuildCommandModel{
		projectPth: projectPth,
		cmdSlice:   []string{constants.XbuildPath},
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

// Run ...
func (xbuild XbuildCommandModel) Run() error {
	if xbuild.projectPth != "" {
		xbuild.cmdSlice = append(xbuild.cmdSlice, xbuild.projectPth)
	}

	if xbuild.target != "" {
		xbuild.cmdSlice = append(xbuild.cmdSlice, fmt.Sprintf("/target:%s", xbuild.target))
	}

	if xbuild.configuration != "" {
		xbuild.cmdSlice = append(xbuild.cmdSlice, fmt.Sprintf("/p:Configuration=%s", xbuild.configuration))
	}

	if xbuild.platform != "" {
		xbuild.cmdSlice = append(xbuild.cmdSlice, fmt.Sprintf("/p:Platform=%s", xbuild.platform))
	}

	if xbuild.buildIpa {
		xbuild.cmdSlice = append(xbuild.cmdSlice, "/p:BuildIpa=true")
	}

	xbuild.cmdSlice = append(xbuild.cmdSlice, "/verbosity:minimal", "/nologo")

	log.Info("=> %s", cmdex.PrintableCommandArgs(false, xbuild.cmdSlice))

	command, err := cmdex.NewCommandFromSlice(xbuild.cmdSlice)
	if err != nil {
		return err
	}

	command.SetStdout(os.Stdout)
	command.SetStderr(os.Stderr)

	return command.Run()
}
