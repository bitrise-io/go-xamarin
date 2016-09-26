package cli

import (
	"fmt"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-xamarin/builder"
	"github.com/urfave/cli"
)

func build(c *cli.Context) error {
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

	builder, err := builder.New(solutionPth)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if err := builder.Build(solutionConfiguration, solutionPlatform, forceMdtool); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	output, err := builder.CollectOutput(solutionConfiguration, solutionPlatform, forceMdtool)
	if err != nil {
		return err
	}

	for outputType, pth := range output {
		log.Done("%s: %s", outputType, pth)
	}

	return nil
}
