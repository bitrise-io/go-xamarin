package builder

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/bitrise-tools/go-xamarin/project"
	"github.com/bitrise-tools/go-xamarin/solution"
	"github.com/bitrise-tools/go-xamarin/utility"
)

// Model ...
type Model struct {
	solution solution.Model

	projectTypeWhitelist []constants.ProjectType
	forceMDTool          bool
}

// OutputMap ...
type OutputMap map[constants.ProjectType]map[constants.OutputType]string

// BuildSolutionCommandCallback ...
type BuildSolutionCommandCallback func(command BuildCommand)

// BuildCommandCallback ...
type BuildCommandCallback func(project project.Model, command BuildCommand)

// ClearCommandCallback ...
type ClearCommandCallback func(project project.Model, dir string)

// New ...
func New(solutionPth string, projectTypeWhitelist []constants.ProjectType, forceMDTool bool) (Model, error) {
	if err := validateSolutionPth(solutionPth); err != nil {
		return Model{}, err
	}

	solution, err := solution.New(solutionPth, true)
	if err != nil {
		return Model{}, err
	}

	if projectTypeWhitelist == nil {
		projectTypeWhitelist = []constants.ProjectType{}
	}

	return Model{
		solution: solution,

		projectTypeWhitelist: projectTypeWhitelist,
		forceMDTool:          forceMDTool,
	}, nil
}

func (builder Model) filteredProjects() []project.Model {
	projects := []project.Model{}

	for _, proj := range builder.solution.ProjectMap {
		if !isProjectTypeAllowed(proj.ProjectType, builder.projectTypeWhitelist...) {
			continue
		}

		if proj.ProjectType != constants.ProjectTypeUnknown {
			projects = append(projects, proj)
		}
	}

	return projects
}

func (builder Model) buildableProjects(configuration, platform string) ([]project.Model, error) {
	projects := []project.Model{}

	solutionConfig := utility.ToConfig(configuration, platform)
	filteredProjects := builder.filteredProjects()

	for _, proj := range filteredProjects {
		//
		// Solution config - project config mapping
		_, ok := proj.ConfigMap[solutionConfig]
		if !ok {
			// fmt.Sprintf("project (%s) do not have config for solution config (%s), skipping...", proj.Name, solutionConfig)
			continue
		}

		if (proj.ProjectType == constants.ProjectTypeIos ||
			proj.ProjectType == constants.ProjectTypeMac ||
			proj.ProjectType == constants.ProjectTypeTVOs) &&
			proj.OutputType != "exe" {
			// fmt.Sprintf("project (%s) does not archivable based on output type (%s), skipping...", project.Name, project.OutputType)
			continue
		}
		if proj.ProjectType == constants.ProjectTypeAndroid &&
			!proj.AndroidApplication {
			// fmt.Sprintf("(%s) is not an android application project, skipping...", proj.Name)
			continue
		}

		if proj.ProjectType != constants.ProjectTypeUnknown {
			projects = append(projects, proj)
		}
	}

	return projects, nil
}

// CleanAll ...
func (builder Model) CleanAll(callback ClearCommandCallback) error {
	filteredProjects := builder.filteredProjects()
	for _, proj := range filteredProjects {

		projectDir := filepath.Dir(proj.Pth)

		{
			binPth := filepath.Join(projectDir, "bin")
			if exist, err := pathutil.IsDirExists(binPth); err != nil {
				return err
			} else if exist {
				if callback != nil {
					callback(proj, binPth)
				}

				if err := os.RemoveAll(binPth); err != nil {
					return err
				}
			}
		}

		{
			objPth := filepath.Join(projectDir, "obj")
			if exist, err := pathutil.IsDirExists(objPth); err != nil {
				return err
			} else if exist {
				if callback != nil {
					callback(proj, objPth)
				}

				if err := os.RemoveAll(objPth); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// BuildSolution ...
func (builder Model) BuildSolution(configuration, platform string, callback BuildSolutionCommandCallback) error {
	if err := validateSolutionConfig(builder.solution, configuration, platform); err != nil {
		return err
	}

	if builder.forceMDTool {
		cmd := NewMDToolCommand(builder.solution.Pth).SetTarget("build").SetConfiguration(configuration).SetPlatform(platform)

		if callback != nil {
			callback(cmd)
		}

		if err := cmd.Run(); err != nil {
			return err
		}
	} else {
		cmd := NewXbuildCommand(builder.solution.Pth).SetTarget("Build").SetConfiguration(configuration).SetPlatform(platform)

		if callback != nil {
			callback(cmd)
		}

		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

// BuildAllProjects ...
func (builder Model) BuildAllProjects(configuration, platform string, callback BuildCommandCallback) error {
	buildableProjects, err := builder.buildableProjects(configuration, platform)
	if err != nil {
		return fmt.Errorf("Failed to list buildable project, error: %s", err)
	}

	solutionConfig := utility.ToConfig(configuration, platform)

	for _, proj := range buildableProjects {
		projectConfigKey, ok := proj.ConfigMap[solutionConfig]
		if !ok {
			// fmt.Sprintf("project (%s) do not have config for solution config (%s), skipping...", proj.Name, solutionConfig)
			continue
		}

		projectConfig, ok := proj.Configs[projectConfigKey]
		if !ok {
			// fmt.Sprintf("project (%s) contains mapping for solution config (%s), but does not have project configuration", proj.Name, solutionConfig)
			continue
		}

		switch proj.ProjectType {
		case constants.ProjectTypeIos, constants.ProjectTypeTVOs:
			if builder.forceMDTool {
				cmd := NewMDToolCommand(builder.solution.Pth).SetTarget("build").SetConfiguration(projectConfig.Configuration).SetPlatform(projectConfig.Platform).SetProjectName(proj.Name)

				if callback != nil {
					callback(proj, cmd)
				}

				if err := cmd.Run(); err != nil {
					return err
				}

				if isArchitectureArchiveable(projectConfig.MtouchArchs) {
					cmd := NewMDToolCommand(builder.solution.Pth).SetTarget("archive").SetConfiguration(projectConfig.Configuration).SetPlatform(projectConfig.Platform).SetProjectName(proj.Name)

					if callback != nil {
						callback(proj, cmd)
					}

					if err := cmd.Run(); err != nil {
						return err
					}
				}
			} else {
				cmd := NewXbuildCommand(builder.solution.Pth).SetTarget("Build").SetConfiguration(configuration).SetPlatform(platform)

				if isArchitectureArchiveable(projectConfig.MtouchArchs) {
					cmd.SetBuildIpa()
					cmd.SetArchiveOnBuild()
				}

				if callback != nil {
					callback(proj, cmd)
				}

				if err := cmd.Run(); err != nil {
					return err
				}
			}
		case constants.ProjectTypeMac:
			if builder.forceMDTool {
				cmd := NewMDToolCommand(builder.solution.Pth).SetTarget("build").SetConfiguration(projectConfig.Configuration).SetPlatform(projectConfig.Platform).SetProjectName(proj.Name)

				if callback != nil {
					callback(proj, cmd)
				}

				if err := cmd.Run(); err != nil {
					return err
				}

				cmd = NewMDToolCommand(builder.solution.Pth).SetTarget("archive").SetConfiguration(projectConfig.Configuration).SetPlatform(projectConfig.Platform).SetProjectName(proj.Name)

				if callback != nil {
					callback(proj, cmd)
				}

				if err := cmd.Run(); err != nil {
					return err
				}
			} else {
				cmd := NewXbuildCommand(builder.solution.Pth).SetTarget("Build").SetConfiguration(configuration).SetPlatform(platform)
				cmd.SetArchiveOnBuild()

				if callback != nil {
					callback(proj, cmd)
				}

				if err := cmd.Run(); err != nil {
					return err
				}
			}
		case constants.ProjectTypeAndroid:
			cmd := NewXbuildCommand(proj.Pth).SetConfiguration(projectConfig.Configuration)

			if projectConfig.SignAndroid {
				cmd.SetTarget("SignAndroidPackage")
			} else {
				cmd.SetTarget("PackageForAndroid")
			}

			if !isPlatformAnyCPU(projectConfig.Platform) {
				cmd.SetPlatform(projectConfig.Platform)
			}

			if callback != nil {
				callback(proj, cmd)
			}

			if err := cmd.Run(); err != nil {
				return err
			}
		}
	}

	return nil
}

// CollectOutput ...
func (builder Model) CollectOutput(configuration, platform string) (OutputMap, error) {
	outputMap := OutputMap{}

	buildableProjects, err := builder.buildableProjects(configuration, platform)
	if err != nil {
		return OutputMap{}, fmt.Errorf("Failed to list buildable project, error: %s", err)
	}

	solutionConfig := utility.ToConfig(configuration, platform)

	for _, proj := range buildableProjects {
		projectConfigKey, ok := proj.ConfigMap[solutionConfig]
		if !ok {
			// fmt.Sprintf("project (%s) do not have config for solution config (%s), skipping...", proj.Name, solutionConfig)
			continue
		}

		projectConfig, ok := proj.Configs[projectConfigKey]
		if !ok {
			// fmt.Sprintf("project (%s) contains mapping for solution config (%s), but does not have project configuration", proj.Name, solutionConfig)
			continue
		}

		projectTypeOutputMap, ok := outputMap[proj.ProjectType]
		if !ok {
			projectTypeOutputMap = map[constants.OutputType]string{}
		}

		switch proj.ProjectType {
		case constants.ProjectTypeIos, constants.ProjectTypeTVOs:
			if xcarchivePth, err := exportLatestXCArchiveFromXcodeArchives(proj.AssemblyName); err != nil {
				return OutputMap{}, err
			} else if xcarchivePth != "" {
				projectTypeOutputMap[constants.OutputTypeXCArchive] = xcarchivePth
			}
			if ipaPth, err := exportIpa(projectConfig.OutputDir, proj.AssemblyName); err != nil {
				return OutputMap{}, err
			} else if ipaPth != "" {
				projectTypeOutputMap[constants.OutputTypeIPA] = ipaPth
			}
		case constants.ProjectTypeMac:
			if xcarchivePth, err := exportLatestXCArchiveFromXcodeArchives(proj.AssemblyName); err != nil {
				return OutputMap{}, err
			} else if xcarchivePth != "" {
				projectTypeOutputMap[constants.OutputTypeXCArchive] = xcarchivePth
			}
			if appPth, err := exportApp(projectConfig.OutputDir, proj.AssemblyName); err != nil {
				return OutputMap{}, err
			} else if appPth != "" {
				projectTypeOutputMap[constants.OutputTypeAPP] = appPth
			}

			if pkgPth, err := exportPkg(projectConfig.OutputDir, proj.AssemblyName); err != nil {
				return OutputMap{}, err
			} else if pkgPth != "" {
				projectTypeOutputMap[constants.OutputTypePKG] = pkgPth
			}
		case constants.ProjectTypeAndroid:
			if apkPth, err := exportApk(projectConfig.OutputDir, proj.ManifestPth, projectConfig.SignAndroid); err != nil {
				return OutputMap{}, err
			} else if apkPth != "" {
				projectTypeOutputMap[constants.OutputTypeAPK] = apkPth
			}
		}

		outputMap[proj.ProjectType] = projectTypeOutputMap
	}

	return outputMap, nil
}
