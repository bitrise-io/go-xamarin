package builder

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/bitrise-tools/go-xamarin/project"
	"github.com/bitrise-tools/go-xamarin/solution"
	"github.com/bitrise-tools/go-xamarin/utility"
)

// Model ...
type Model struct {
	solution solution.Model
}

// OutputMap ...
type OutputMap map[constants.ProjectType]map[constants.OutputType]string

// ProjectIterator ...
type ProjectIterator func(project project.Model) error

// ProjectWithConfigIterator ...
type ProjectWithConfigIterator func(project project.Model, projectConfig project.ConfigurationPlatformModel) error

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

// IterateOnAllProjects ...
func (builder Model) IterateOnAllProjects(projectTypeWhiteList []constants.ProjectType, iterator ProjectIterator) error {
	for _, project := range builder.solution.ProjectMap {
		if !isProjectTypeAllowed(project.ProjectType, projectTypeWhiteList...) {
			log.Detail("project whitelist omitts project: %s", project.Name)
			continue
		}

		if err := iterator(project); err != nil {
			return err
		}
	}
	return nil
}

// IterateOnBuildableProjects ...
func (builder Model) IterateOnBuildableProjects(configuration, platform string, projectTypeWhiteList []constants.ProjectType, iterator ProjectWithConfigIterator) error {
	if err := validateSolutionConfig(builder.solution, configuration, platform); err != nil {
		return err
	}

	solutionConfig := utility.ToConfig(configuration, platform)

	for _, project := range builder.solution.ProjectMap {
		if !isProjectTypeAllowed(project.ProjectType, projectTypeWhiteList...) {
			log.Detail("project whitelist omitts project: %s", project.Name)
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
			return fmt.Errorf("project contains mapping for solution config (%s -> %s), but does not have project config for it", solutionConfig, projectConfig)
		}

		if (project.ProjectType == constants.ProjectTypeIos ||
			project.ProjectType == constants.ProjectTypeMac ||
			project.ProjectType == constants.ProjectTypeTVOs) &&
			project.OutputType != "exe" {
			log.Warn("project (%s) does not archivable based on output type (%s), skipping...", project.Name, project.OutputType)
			continue
		}
		if project.ProjectType == constants.ProjectTypeAndroid &&
			!project.AndroidApplication {
			log.Warn("(%s) is not an android application project, skipping...", project.Name)
			continue
		}

		if project.ProjectType != constants.ProjectTypeUnknown {
			if err := iterator(project, config); err != nil {
				return err
			}
		}
	}

	return nil
}

// CleanAll ...
func (builder Model) CleanAll(projectTypeFilter ...constants.ProjectType) error {
	log.Info("Cleaning project output dirs (bin, obj) ...")

	iterator := func(project project.Model) error {
		projectDir := filepath.Dir(project.Pth)
		binPth := filepath.Join(projectDir, "bin")
		objPth := filepath.Join(projectDir, "obj")

		if exist, err := pathutil.IsDirExists(binPth); err != nil {
			return err
		} else if exist {
			log.Detail("remove: %s", binPth)
			if err := os.RemoveAll(binPth); err != nil {
				return err
			}
		}

		if exist, err := pathutil.IsDirExists(objPth); err != nil {
			return err
		} else if exist {
			log.Detail("remove: %s", binPth)
			if err := os.RemoveAll(objPth); err != nil {
				return err
			}
		}

		return nil
	}

	return builder.IterateOnAllProjects(projectTypeFilter, iterator)
}

// Build ...
func (builder Model) Build(configuration, platform string, forceMDTool bool, projectTypeFilter ...constants.ProjectType) error {
	iterator := func(project project.Model, projectConfig project.ConfigurationPlatformModel) error {
		switch project.ProjectType {
		case constants.ProjectTypeIos, constants.ProjectTypeTVOs:
			if forceMDTool {
				if err := NewMDToolCommand(builder.solution.Pth).SetTarget("build").SetConfiguration(projectConfig.Configuration).SetPlatform(projectConfig.Platform).SetProjectName(project.Name).Run(); err != nil {
					return err
				}

				if isArchitectureArchiveable(projectConfig.MtouchArchs) {
					if err := NewMDToolCommand(builder.solution.Pth).SetTarget("archive").SetConfiguration(projectConfig.Configuration).SetPlatform(projectConfig.Platform).SetProjectName(project.Name).Run(); err != nil {
						return err
					}
				}
			} else {
				command := NewXbuildCommand(builder.solution.Pth).SetTarget("Build").SetConfiguration(configuration).SetPlatform(platform)

				if projectConfig.BuildIpa {
					command.SetBuildIpa()
				} else if isArchitectureArchiveable(projectConfig.MtouchArchs) {
					command.SetArchiveOnBuild()
				}

				if err := command.Run(); err != nil {
					return err
				}
			}
		case constants.ProjectTypeMac:
			if forceMDTool {
				if err := NewMDToolCommand(builder.solution.Pth).SetTarget("build").SetConfiguration(projectConfig.Configuration).SetPlatform(projectConfig.Platform).SetProjectName(project.Name).Run(); err != nil {
					return err
				}

				if err := NewMDToolCommand(builder.solution.Pth).SetTarget("archive").SetConfiguration(projectConfig.Configuration).SetPlatform(projectConfig.Platform).SetProjectName(project.Name).Run(); err != nil {
					return err
				}
			} else {
				if err := NewXbuildCommand(builder.solution.Pth).SetTarget("Build").SetConfiguration(configuration).SetPlatform(platform).Run(); err != nil {
					return err
				}
			}
		case constants.ProjectTypeAndroid:
			command := NewXbuildCommand(project.Pth).SetConfiguration(projectConfig.Configuration)

			if projectConfig.SignAndroid {
				command.SetTarget("SignAndroidPackage")
			} else {
				command.SetTarget("PackageForAndroid")
			}

			if !isPlatformAnyCPU(projectConfig.Platform) {
				command.SetPlatform(projectConfig.Platform)
			}

			if err := command.Run(); err != nil {
				return err
			}
		}

		return nil
	}

	return builder.IterateOnBuildableProjects(configuration, platform, projectTypeFilter, iterator)
}

// CollectOutput ...
func (builder Model) CollectOutput(configuration, platform string, forceMDTool bool, projectTypeFilter ...constants.ProjectType) (OutputMap, error) {
	outputMap := OutputMap{}

	iterator := func(project project.Model, projectConfig project.ConfigurationPlatformModel) error {
		switch project.ProjectType {
		case constants.ProjectTypeIos, constants.ProjectTypeTVOs:
			projectTypeOutputMap := outputMap[project.ProjectType]

			if forceMDTool {
				if isArchitectureArchiveable(projectConfig.MtouchArchs) {
					if xcarchivePth, err := exportLatestXCArchiveFromXcodeArchives(project.AssemblyName); err != nil {
						return err
					} else if xcarchivePth != "" {
						projectTypeOutputMap[constants.OutputTypeXCArchive] = xcarchivePth
					}

					if projectConfig.BuildIpa {
						if ipaPth, err := exportIpa(projectConfig.OutputDir, project.AssemblyName); err != nil {
							return err
						} else if ipaPth != "" {
							projectTypeOutputMap[constants.OutputTypeIPA] = ipaPth
						}
					}
				}
			} else {
				if projectConfig.BuildIpa {
					if ipaPth, err := exportIpa(projectConfig.OutputDir, project.AssemblyName); err != nil {
						return err
					} else if ipaPth != "" {
						projectTypeOutputMap[constants.OutputTypeIPA] = ipaPth
					}
				} else if isArchitectureArchiveable(projectConfig.MtouchArchs) {
					if xcarchivePth, err := exportLatestXCArchiveFromXcodeArchives(project.AssemblyName); err != nil {
						return err
					} else if xcarchivePth != "" {
						projectTypeOutputMap[constants.OutputTypeXCArchive] = xcarchivePth
					}
				}
			}
		case constants.ProjectTypeMac:
			projectTypeOutputMap := outputMap[project.ProjectType]

			if appPth, err := exportApp(projectConfig.OutputDir, project.AssemblyName); err != nil {
				return err
			} else if appPth != "" {
				projectTypeOutputMap[constants.OutputTypeAPP] = appPth
			}

			if pkgPth, err := exportPkg(projectConfig.OutputDir, project.AssemblyName); err != nil {
				return err
			} else if pkgPth != "" {
				projectTypeOutputMap[constants.OutputTypePKG] = pkgPth
			}
		case constants.ProjectTypeAndroid:
			projectTypeOutputMap := outputMap[project.ProjectType]

			if apkPth, err := exportApk(projectConfig.OutputDir, project.ManifestPth, projectConfig.SignAndroid); err != nil {
				return err
			} else if apkPth != "" {
				projectTypeOutputMap[constants.OutputTypeAPK] = apkPth
			}
		}

		return nil
	}

	if err := builder.IterateOnBuildableProjects(configuration, platform, projectTypeFilter, iterator); err != nil {
		return OutputMap{}, err
	}

	return outputMap, nil
}
