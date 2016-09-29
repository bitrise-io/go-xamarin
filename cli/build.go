package cli

import (
	"fmt"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-xamarin/builder"
	"github.com/bitrise-tools/go-xamarin/project"
	"github.com/urfave/cli"
)

func buildCmd(c *cli.Context) error {
	solutionPth := c.String(solutionFilePathKey)
	solutionConfiguration := c.String(solutionConfigurationKey)
	solutionPlatform := c.String(solutionPlatformKey)
	forceMdtool := c.Bool(forceMDToolKey)

	fmt.Println("")
	log.Info("Config:")
	log.Detail("- solution: %s", solutionPth)
	log.Detail("- configuration: %s", solutionConfiguration)
	log.Detail("- platform: %s", solutionPlatform)
	log.Detail("- force-mdtool: %v", forceMdtool)
	fmt.Println("")

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

	callback := func(project project.Model, command builder.BuildCommand) {
		fmt.Println()
		log.Info("Building project: %s", project.Name)
		log.Detail("-> %s", command.PrintableCommand())
		fmt.Println()
	}

	err = buildHandler.BuildAllProjects(solutionConfiguration, solutionPlatform, callback)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	output, err := buildHandler.CollectOutput(solutionConfiguration, solutionPlatform)
	if err != nil {
		return err
	}

	for outputType, pth := range output {
		log.Done("%s: %s", outputType, pth)
	}

	return nil
}
