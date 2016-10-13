package builder

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-xamarin/buildtool"
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

// ProjectOutputModel ...
type ProjectOutputModel struct {
	Pth        string
	OutputType constants.OutputType
}

// TestProjectOutputModel ...
type TestProjectOutputModel struct {
	Pth                  string
	OutputType           constants.OutputType
	ReferredProjectNames []string
}

// ProjectOutputMap ...
type ProjectOutputMap map[string][]ProjectOutputModel // Project Name - ProjectOutputModels

// TestProjectOutputMap ...
type TestProjectOutputMap map[string]TestProjectOutputModel // Test Project Name - TestProjectOutputModel

// PrepareBuildCommandCallback ...
type PrepareBuildCommandCallback func(project project.Model, command *buildtool.EditableCommand)

// BuildCommandCallback ...
type BuildCommandCallback func(project project.Model, command buildtool.PrintableCommand, alreadyPerformed bool)

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

// CleanAll ...
func (builder Model) CleanAll(callback ClearCommandCallback) error {
	whitelistedProjects := builder.whitelistedProjects()

	for _, proj := range whitelistedProjects {

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

// BuildAllProjects ...
func (builder Model) BuildAllProjects(configuration, platform string, prepareCallback PrepareBuildCommandCallback, callback BuildCommandCallback) ([]string, error) {
	warnings := []string{}

	if err := validateSolutionConfig(builder.solution, configuration, platform); err != nil {
		return []string{}, err
	}

	buildableProjects, warns := builder.buildableProjects(configuration, platform)
	if len(buildableProjects) == 0 {
		return warns, nil
	}

	perfomedCommands := []buildtool.RunnableCommand{}

	for _, proj := range buildableProjects {
		buildCommands, warns := builder.buildCommandsForProject(configuration, platform, proj)
		warnings = append(warnings, warns...)

		for _, buildCommand := range buildCommands {
			// Callback to let the caller to modify the command
			if prepareCallback != nil {
				editabeCommand := buildtool.EditableCommand(buildCommand)
				prepareCallback(proj, &editabeCommand)
			}

			// Check if same command was already performed
			alreadyPerformed := false
			if buildtool.BuildCommandSliceContains(perfomedCommands, buildCommand) {
				alreadyPerformed = true
			}

			// Callback to notify the caller about next running command
			if callback != nil {
				callback(proj, buildCommand, alreadyPerformed)
			}

			if !alreadyPerformed {
				if err := buildCommand.Run(); err != nil {
					return warnings, err
				}
				perfomedCommands = append(perfomedCommands, buildCommand)
			}
		}
	}

	return warnings, nil
}

// BuildAllProjectsAndXamarinUITestprojects ...
func (builder Model) BuildAllProjectsAndXamarinUITestprojects(configuration, platform string, prepareCallback PrepareBuildCommandCallback, callback BuildCommandCallback) ([]string, error) {
	warnings := []string{}

	if err := validateSolutionConfig(builder.solution, configuration, platform); err != nil {
		return []string{}, err
	}

	buildableTestProjects, buildableReferredProjects, warns := builder.buildableXamarinUITestProjectsAndReferredProjects(configuration, platform)
	if len(buildableTestProjects) == 0 || len(buildableReferredProjects) == 0 {
		return warns, nil
	}

	perfomedCommands := []buildtool.RunnableCommand{}

	for _, proj := range buildableReferredProjects {
		buildCommands, warns := builder.buildCommandsForProject(configuration, platform, proj)
		warnings = append(warnings, warns...)

		for _, buildCommand := range buildCommands {
			// Callback to let the caller to modify the command
			if prepareCallback != nil {
				editabeCommand := buildtool.EditableCommand(buildCommand)
				prepareCallback(proj, &editabeCommand)
			}

			// Check if same command was already performed
			alreadyPerformed := false
			if buildtool.BuildCommandSliceContains(perfomedCommands, buildCommand) {
				alreadyPerformed = true
			}

			// Callback to notify the caller about next running command
			if callback != nil {
				callback(proj, buildCommand, alreadyPerformed)
			}

			if !alreadyPerformed {
				if err := buildCommand.Run(); err != nil {
					return warnings, err
				}
				perfomedCommands = append(perfomedCommands, buildCommand)
			}
		}
	}

	for _, testProj := range buildableTestProjects {
		buildCommand, warns := builder.buildCommandForTestProject(configuration, platform, testProj)
		warnings = append(warnings, warns...)

		// Callback to let the caller to modify the command
		if prepareCallback != nil {
			editabeCommand := buildtool.EditableCommand(buildCommand)
			prepareCallback(testProj, &editabeCommand)
		}

		// Check if same command was already performed
		alreadyPerformed := false
		if buildtool.BuildCommandSliceContains(perfomedCommands, buildCommand) {
			alreadyPerformed = true
		}

		// Callback to notify the caller about next running command
		if callback != nil {
			callback(testProj, buildCommand, alreadyPerformed)
		}

		if !alreadyPerformed {
			if err := buildCommand.Run(); err != nil {
				return warnings, err
			}
			perfomedCommands = append(perfomedCommands, buildCommand)
		}
	}

	return warnings, nil
}

// CollectProjectOutputs ...
func (builder Model) CollectProjectOutputs(configuration, platform string) (ProjectOutputMap, error) {
	projectOutputMap := ProjectOutputMap{}

	buildableProjects, _ := builder.buildableProjects(configuration, platform)

	solutionConfig := utility.ToConfig(configuration, platform)

	for _, proj := range buildableProjects {
		projectConfigKey, ok := proj.ConfigMap[solutionConfig]
		if !ok {
			continue
		}

		projectConfig, ok := proj.Configs[projectConfigKey]
		if !ok {
			continue
		}

		projectOutputs, ok := projectOutputMap[proj.Name]
		if !ok {
			projectOutputs = []ProjectOutputModel{}
		}

		switch proj.ProjectType {
		case constants.ProjectTypeIOS, constants.ProjectTypeTvOS:
			if xcarchivePth, err := exportLatestXCArchiveFromXcodeArchives(proj.AssemblyName); err != nil {
				return ProjectOutputMap{}, err
			} else if xcarchivePth != "" {
				projectOutputs = append(projectOutputs, ProjectOutputModel{
					Pth:        xcarchivePth,
					OutputType: constants.OutputTypeXCArchive,
				})
			}
			if ipaPth, err := exportLatestIpa(projectConfig.OutputDir, proj.AssemblyName); err != nil {
				return ProjectOutputMap{}, err
			} else if ipaPth != "" {
				projectOutputs = append(projectOutputs, ProjectOutputModel{
					Pth:        ipaPth,
					OutputType: constants.OutputTypeIPA,
				})
			}
			if dsymPth, err := exportAppDSYM(projectConfig.OutputDir, proj.AssemblyName); err != nil {
				return ProjectOutputMap{}, err
			} else if dsymPth != "" {
				projectOutputs = append(projectOutputs, ProjectOutputModel{
					Pth:        dsymPth,
					OutputType: constants.OutputTypeDSYM,
				})
			}
		case constants.ProjectTypeMacOS:
			if builder.forceMDTool {
				if xcarchivePth, err := exportLatestXCArchiveFromXcodeArchives(proj.AssemblyName); err != nil {
					return ProjectOutputMap{}, err
				} else if xcarchivePth != "" {
					projectOutputs = append(projectOutputs, ProjectOutputModel{
						Pth:        xcarchivePth,
						OutputType: constants.OutputTypeXCArchive,
					})
				}
			}
			if appPth, err := exportApp(projectConfig.OutputDir, proj.AssemblyName); err != nil {
				return ProjectOutputMap{}, err
			} else if appPth != "" {
				projectOutputs = append(projectOutputs, ProjectOutputModel{
					Pth:        appPth,
					OutputType: constants.OutputTypeAPP,
				})
			}
			if pkgPth, err := exportPKG(projectConfig.OutputDir, proj.AssemblyName); err != nil {
				return ProjectOutputMap{}, err
			} else if pkgPth != "" {
				projectOutputs = append(projectOutputs, ProjectOutputModel{
					Pth:        pkgPth,
					OutputType: constants.OutputTypePKG,
				})
			}
		case constants.ProjectTypeAndroid:
			packageName, err := androidPackageName(proj.ManifestPth)
			if err != nil {
				return ProjectOutputMap{}, err
			}

			if apkPth, err := exportApk(projectConfig.OutputDir, packageName); err != nil {
				return ProjectOutputMap{}, err
			} else if apkPth != "" {
				projectOutputs = append(projectOutputs, ProjectOutputModel{
					Pth:        apkPth,
					OutputType: constants.OutputTypeAPK,
				})
			}
		}

		if len(projectOutputs) > 0 {
			projectOutputMap[proj.Name] = projectOutputs
		}
	}

	return projectOutputMap, nil
}

// CollectXamarinUITestProjectOutputs ...
func (builder Model) CollectXamarinUITestProjectOutputs(configuration, platform string) (TestProjectOutputMap, error) {
	testProjectOutputMap := TestProjectOutputMap{}

	buildableTestProjects, _, _ := builder.buildableXamarinUITestProjectsAndReferredProjects(configuration, platform)

	solutionConfig := utility.ToConfig(configuration, platform)

	for _, testProj := range buildableTestProjects {
		projectConfigKey, ok := testProj.ConfigMap[solutionConfig]
		if !ok {
			continue
		}

		projectConfig, ok := testProj.Configs[projectConfigKey]
		if !ok {
			continue
		}

		if dllPth, err := exportDLL(projectConfig.OutputDir, testProj.AssemblyName); err != nil {
			return TestProjectOutputMap{}, err
		} else if dllPth != "" {
			referredProjectNames := []string{}
			referredProjectIDs := testProj.ReferredProjectIDs
			for _, referredProjectID := range referredProjectIDs {
				referredProject, ok := builder.solution.ProjectMap[referredProjectID]
				if !ok {
					return TestProjectOutputMap{}, fmt.Errorf("project reference exist with project id: %s, but project not found in solution", referredProjectID)
				}

				referredProjectNames = append(referredProjectNames, referredProject.Name)
			}

			testProjectOutputMap[testProj.Name] = TestProjectOutputModel{
				Pth:                  dllPth,
				OutputType:           constants.OutputTypeDLL,
				ReferredProjectNames: referredProjectNames,
			}
		}
	}

	return testProjectOutputMap, nil
}
