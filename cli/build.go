package cli

import (
	"fmt"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-xamarin/builder"
	"github.com/bitrise-tools/go-xamarin/constants"
	"github.com/urfave/cli"
)

func buildCmd(c *cli.Context) error {
	solutionPth := c.String(solutionFilePathKey)
	solutionConfiguration := c.String(solutionConfigurationKey)
	solutionPlatform := c.String(solutionPlatformKey)
	forceMdtool := c.Bool(forceMDToolKey)

	fmt.Println()
	log.Info("Config:")
	log.Detail("- solution: %s", solutionPth)
	log.Detail("- configuration: %s", solutionConfiguration)
	log.Detail("- platform: %s", solutionPlatform)
	log.Detail("- force-mdtool: %v", forceMdtool)

	if solutionPth == "" {
		return fmt.Errorf("missing required input: %s", solutionFilePathKey)
	}
	if solutionConfiguration == "" {
		return fmt.Errorf("missing required input: %s", solutionConfigurationKey)
	}
	if solutionPlatform == "" {
		return fmt.Errorf("missing required input: %s", solutionPlatformKey)
	}

	buildHandler, err := builder.New(solutionPth, nil, forceMdtool)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println()
	log.Info("Building all projects in solution: %s", solutionPth)

	callback := func(solutionName string, projectName string, sdk constants.SDK, testFramwork constants.TestFramework, commandStr string, alreadyPerformed bool) {
		if projectName != "" {
			fmt.Println()
			log.Info("Building project: %s", projectName)
			log.Done("$ %s", commandStr)
			if alreadyPerformed {
				log.Warn("build command already performed, skipping...")
			}
			fmt.Println()
		}
	}

	warnings, err := buildHandler.BuildAllProjects(solutionConfiguration, solutionPlatform, nil, callback)
	for _, warning := range warnings {
		log.Warn(warning)
	}
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println()
	log.Info("Collecting generated outputs")

	outputMap, err := buildHandler.CollectProjectOutputs(solutionConfiguration, solutionPlatform)
	if err != nil {
		return err
	}

	for projectName, projectOutput := range outputMap {
		fmt.Println()
		log.Info("%s outputs:", projectName)

		for _, output := range projectOutput.Outputs {
			log.Done("%s: %s", output.OutputType, output.Pth)
		}
	}

	return nil
}
