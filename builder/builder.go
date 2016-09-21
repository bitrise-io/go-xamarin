package builder

import (
	"fmt"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/xamarin-builder/constants"
	"github.com/bitrise-tools/xamarin-builder/solution"
	"github.com/bitrise-tools/xamarin-builder/utility"
)

// Model ...
type Model struct {
	solution solution.Model
}

// OutputMap ...
type OutputMap map[string]string

// New ...
func New(solutionPth string) (Model, error) {
	if err := validateSolutionPth(solutionPth); err != nil {
		return Model{}, err
	}

	solution, err := solution.New(solutionPth, true)
	if err != nil {
		return Model{}, err
	}

	return Model{
		solution: solution,
	}, nil
}

// BuildSolution ...
func (builder Model) BuildSolution(configuration, platform string, forceMDTool bool) (OutputMap, error) {
	if err := validateSolutionConfig(builder.solution, configuration, platform); err != nil {
		return OutputMap{}, err
	}

	if forceMDTool {
		if err := NewMDToolCommand(builder.solution.Pth).SetTarget("build").SetConfiguration(configuration).SetPlatform(platform).Run(); err != nil {
			return OutputMap{}, err
		}
	} else {
		if err := NewXbuildCommand(builder.solution.Pth).SetTarget("Build").SetConfiguration(configuration).SetPlatform(platform).Run(); err != nil {
			return OutputMap{}, err
		}
	}

	return OutputMap{}, nil
}

// Build ...
func (builder Model) Build(configuration, platform string, forceMDTool bool, projectTypeFilter ...constants.ProjectType) (OutputMap, error) {
	if err := validateSolutionConfig(builder.solution, configuration, platform); err != nil {
		return OutputMap{}, err
	}

	outputMap := OutputMap{}

	solutionConfig := utility.ToConfig(configuration, platform)

	for _, project := range builder.solution.ProjectMap {
		if isFilterForProjectType(project.ProjectType, projectTypeFilter...) {
			continue
		}

		//
		// Solution config - project config mapping
		projectConfig, ok := project.ConfigMap[solutionConfig]
		if !ok {
			log.Warn("project (%s) do not have config for solution config (%s), skipping...", project.Name, solutionConfig)
			continue
		}

		config, ok := project.Configs[projectConfig]
		if !ok {
			return OutputMap{}, fmt.Errorf("project contains mapping for solution config (%s -> %s), but does not have project config for it", solutionConfig, projectConfig)
		}

		//
		// Building projects
		switch project.ProjectType {

		case constants.XamarinIos, constants.XamarinTVOS:
			if project.OutputType != "exe" {
				log.Warn("project (%s) does not archivable based on output type (%s), skipping...", project.Name, project.OutputType)
				continue
			}

			archiveabel := archiveableArchitecture(config.MtouchArchs)

			if forceMDTool {
				if err := NewMDToolCommand(builder.solution.Pth).SetTarget("build").SetConfiguration(config.Configuration).SetPlatform(config.Platform).SetProjectName(project.Name).Run(); err != nil {
					return OutputMap{}, err
				}

				if archiveabel {
					if err := NewMDToolCommand(builder.solution.Pth).SetTarget("archive").SetConfiguration(config.Configuration).SetPlatform(config.Platform).SetProjectName(project.Name).Run(); err != nil {
						return OutputMap{}, err
					}

					xcarchivePth, err := exportLatestXCArchiveFromXcodeArchives(project.AssemblyName)
					if err != nil {
						return OutputMap{}, err
					}

					outputMap["xcarchive"] = xcarchivePth

					if config.BuildIpa {
						ipaPth, err := exportIpa(config.OutputDir, project.AssemblyName)
						if err != nil {
							return OutputMap{}, err
						}

						outputMap["ipa"] = ipaPth
					}
				}
			} else {
				command := NewXbuildCommand(builder.solution.Pth).SetTarget("Build").SetConfiguration(configuration).SetPlatform(platform)

				if config.BuildIpa {
					command.SetBuildIpa()

					if err := command.Run(); err != nil {
						return OutputMap{}, err
					}

					ipaPth, err := exportIpa(config.OutputDir, project.AssemblyName)
					if err != nil {
						return OutputMap{}, err
					}

					outputMap["ipa"] = ipaPth
				} else {
					if err := command.Run(); err != nil {
						return OutputMap{}, err
					}
				}
			}
		case constants.XamarinMac, constants.MonoMac:
			if project.OutputType != "exe" {
				log.Warn("project (%s) does not archivable based on output type (%s), skipping...", project.Name, project.OutputType)
				continue
			}

			if forceMDTool {
				if err := NewMDToolCommand(builder.solution.Pth).SetTarget("build").SetConfiguration(config.Configuration).SetPlatform(config.Platform).SetProjectName(project.Name).Run(); err != nil {
					return OutputMap{}, err
				}

				if err := NewMDToolCommand(builder.solution.Pth).SetTarget("archive").SetConfiguration(config.Configuration).SetPlatform(config.Platform).SetProjectName(project.Name).Run(); err != nil {
					return OutputMap{}, err
				}
			} else {
				if err := NewXbuildCommand(builder.solution.Pth).SetTarget("Build").SetConfiguration(configuration).SetPlatform(platform).Run(); err != nil {
					return OutputMap{}, err
				}
			}

			appPth, err := exportApp(config.OutputDir, project.AssemblyName)
			if err != nil {
				return OutputMap{}, err
			}
			if appPth != "" {
				outputMap["app"] = appPth
			}

			pkgPth, err := exportPkg(config.OutputDir, project.AssemblyName)
			if err != nil {
				return OutputMap{}, err
			}
			if pkgPth != "" {
				outputMap["pkg"] = pkgPth
			}
		case constants.XamarinAndroid:
			if !project.AndroidApplication {
				log.Warn("(%s) is not an android application project, skipping...", project.Name)
				continue
			}

			command := NewXbuildCommand(project.Pth).SetConfiguration(config.Configuration)

			if config.SignAndroid {
				command.SetTarget("SignAndroidPackage")
			} else {
				command.SetTarget("PackageForAndroid")
			}

			if config.Platform != "Any CPU" {
				command.SetPlatform(config.Platform)
			}

			if err := command.Run(); err != nil {
				return OutputMap{}, err
			}

			apkPth, err := exportApk(config.OutputDir, project.ManifestPth, config.SignAndroid)
			if err != nil {
				return OutputMap{}, err
			}
			outputMap["apk"] = apkPth
		}
	}

	return outputMap, nil
}
