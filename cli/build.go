package cli

import (
	"fmt"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/xamarin-builder/builder"
	"github.com/urfave/cli"
)

func build(c *cli.Context) error {
	solutionPth := c.String(solutionFilePathKey)
	if solutionPth == "" {
		return fmt.Errorf("missing required input: %s", solutionFilePathKey)
	}

	solutionConfiguration := c.String(solutionConfigurationKey)
	if solutionConfiguration == "" {
		return fmt.Errorf("missing required input: %s", solutionConfigurationKey)
	}

	solutionPlatform := c.String(solutionPlatformKey)
	if solutionPlatform == "" {
		return fmt.Errorf("missing required input: %s", solutionPlatformKey)
	}

	forceMdtool := c.Bool(forceMDToolKey)

	log.Info("Config:")
	log.Detail("solution: %s", solutionPth)
	log.Detail("configuration: %s", solutionConfiguration)
	log.Detail("platform: %s", solutionPlatform)
	log.Detail("force-mdtool: %v", forceMdtool)
	fmt.Println("")

	builder, err := builder.New(solutionPth)
	if err != nil {
		return err
	}

	output, err := builder.Build(solutionConfiguration, solutionPlatform, forceMdtool)
	if err != nil {
		return err
	}

	for outputType, pth := range output {
		log.Done("%s: %s", outputType, pth)
	}

	return nil
}
