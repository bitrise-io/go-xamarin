package builder

import (
	"fmt"

	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/bitrise-tools/go-xamarin/project"
	"github.com/bitrise-tools/go-xamarin/tools/buildtools"
	"github.com/bitrise-tools/go-xamarin/tools/buildtools/mdtool"
	"github.com/bitrise-tools/go-xamarin/tools/buildtools/xbuild"
	"github.com/bitrise-tools/go-xamarin/utility"
)

func (builder Model) buildSolutionCommand(configuration, platform string) buildtools.RunnableCommand {
	var buildCommand buildtools.RunnableCommand

	if builder.forceMDTool {
		buildCommand := mdtool.New(builder.solution.Pth).SetTarget("build")
		buildCommand.SetConfiguration(configuration)
		buildCommand.SetPlatform(platform)
	} else {
		buildCommand := xbuild.New(builder.solution.Pth).SetTarget("Build")
		buildCommand.SetConfiguration(configuration)
		buildCommand.SetPlatform(platform)
	}

	return buildCommand
}

func (builder Model) buildProjectCommand(configuration, platform string, proj project.Model) ([]buildtools.RunnableCommand, []string) {
	warnings := []string{}

	solutionConfig := utility.ToConfig(configuration, platform)

	projectConfigKey, ok := proj.ConfigMap[solutionConfig]
	if !ok {
		warnings = append(warnings, fmt.Sprintf("project (%s) do not have config for solution config (%s), skipping...", proj.Name, solutionConfig))
	}

	projectConfig, ok := proj.Configs[projectConfigKey]
	if !ok {
		warnings = append(warnings, fmt.Sprintf("project (%s) contains mapping for solution config (%s), but does not have project configuration", proj.Name, solutionConfig))
	}

	// Prepare build commands
	buildCommands := []buildtools.RunnableCommand{}

	switch proj.ProjectType {
	case constants.ProjectTypeIOS, constants.ProjectTypeTvOS:
		if builder.forceMDTool {
			command := mdtool.New(builder.solution.Pth).SetTarget("build")
			command.SetConfiguration(projectConfig.Configuration)
			command.SetPlatform(projectConfig.Platform)
			command.SetProjectName(proj.Name)

			buildCommands = append(buildCommands, command)

			if isArchitectureArchiveable(projectConfig.MtouchArchs...) {
				command := mdtool.New(builder.solution.Pth).SetTarget("archive")
				command.SetConfiguration(projectConfig.Configuration)
				command.SetPlatform(projectConfig.Platform)
				command.SetProjectName(proj.Name)

				buildCommands = append(buildCommands, command)
			}
		} else {
			command := xbuild.New(builder.solution.Pth).SetTarget("Build")
			command.SetConfiguration(configuration)
			command.SetPlatform(platform)

			if isArchitectureArchiveable(projectConfig.MtouchArchs...) {
				command.SetBuildIpa(true)
				command.SetArchiveOnBuild(true)
			}

			buildCommands = append(buildCommands, command)
		}
	case constants.ProjectTypeMacOS:
		if builder.forceMDTool {
			command := mdtool.New(builder.solution.Pth).SetTarget("build")
			command.SetConfiguration(projectConfig.Configuration)
			command.SetPlatform(projectConfig.Platform)
			command.SetProjectName(proj.Name)

			buildCommands = append(buildCommands, command)

			command = mdtool.New(builder.solution.Pth).SetTarget("archive")
			command.SetConfiguration(projectConfig.Configuration)
			command.SetPlatform(projectConfig.Platform)
			command.SetProjectName(proj.Name)

			buildCommands = append(buildCommands, command)
		} else {
			command := xbuild.New(builder.solution.Pth).SetTarget("Build")
			command.SetConfiguration(configuration)
			command.SetPlatform(platform)
			command.SetArchiveOnBuild(true)

			buildCommands = append(buildCommands, command)
		}
	case constants.ProjectTypeAndroid:
		command := xbuild.New(proj.Pth)
		if projectConfig.SignAndroid {
			command.SetTarget("SignAndroidPackage")
		} else {
			command.SetTarget("PackageForAndroid")
		}

		command.SetConfiguration(projectConfig.Configuration)

		if !isPlatformAnyCPU(projectConfig.Platform) {
			command.SetPlatform(projectConfig.Platform)
		}

		buildCommands = append(buildCommands, command)
	}

	return buildCommands, warnings
}

func (builder Model) buildXamarinUITestProjectCommand(configuration, platform string, proj project.Model) (buildtools.RunnableCommand, []string) {
	warnings := []string{}

	solutionConfig := utility.ToConfig(configuration, platform)

	projectConfigKey, ok := proj.ConfigMap[solutionConfig]
	if !ok {
		warnings = append(warnings, fmt.Sprintf("project (%s) do not have config for solution config (%s), skipping...", proj.Name, solutionConfig))
	}

	projectConfig, ok := proj.Configs[projectConfigKey]
	if !ok {
		warnings = append(warnings, fmt.Sprintf("project (%s) contains mapping for solution config (%s), but does not have project configuration", proj.Name, solutionConfig))
	}

	command := mdtool.New(builder.solution.Pth)
	command.SetTarget("build")
	command.SetConfiguration(projectConfig.Configuration)
	command.SetProjectName(proj.Name)

	return command, warnings
}

// func (builder Model) buildNunitTestProjectCommand(configuration, platform string, proj project.Model) (buildtools.RunnableCommand, []string) {
// 	warnings := []string{}

// 	solutionConfig := utility.ToConfig(configuration, platform)

// 	projectConfigKey, ok := proj.ConfigMap[solutionConfig]
// 	if !ok {
// 		warnings = append(warnings, fmt.Sprintf("project (%s) do not have config for solution config (%s), skipping...", proj.Name, solutionConfig))
// 	}

// 	projectConfig, ok := proj.Configs[projectConfigKey]
// 	if !ok {
// 		warnings = append(warnings, fmt.Sprintf("project (%s) contains mapping for solution config (%s), but does not have project configuration", proj.Name, solutionConfig))
// 	}

// }
