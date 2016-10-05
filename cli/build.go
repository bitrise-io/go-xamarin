package cli

import (
	"fmt"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-xamarin/builder"
	"github.com/bitrise-tools/go-xamarin/buildtool"
	"github.com/bitrise-tools/go-xamarin/project"
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

	callback := func(project project.Model, command buildtool.PrintableCommand, alreadyPerformed bool) {
		fmt.Println()
		log.Info("Building project: %s", project.Name)
		log.Done("$ %s", command.PrintableCommand())
		if alreadyPerformed {
			log.Warn("build command already performed, skipping...")
		}
		fmt.Println()
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

	outputMap, warnings := buildHandler.CollectOutput(solutionConfiguration, solutionPlatform)
	for _, warning := range warnings {
		log.Warn(warning)
	}
	if err != nil {
		return err
	}

	for projectType, output := range outputMap {
		fmt.Println()
		log.Info("%s outputs:", projectType)

		for outputType, pth := range output {
			log.Done("%s: %s", outputType, pth)
		}
	}

	return nil
}
