package cli

import (
	"fmt"

	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-tools/go-xamarin/builder"
	"github.com/urfave/cli"
)

func clean(c *cli.Context) error {
	solutionPth := c.String(solutionFilePathKey)

	fmt.Println("")
	log.Info("Config:")
	log.Detail("- solution: %s", solutionPth)
	fmt.Println("")

	if solutionPth == "" {
		return fmt.Errorf("missing required input: %s", solutionFilePathKey)
	}

	builder, err := builder.New(solutionPth)
	if err != nil {
		return err
	}

	if err := builder.CleanAll(); err != nil {
		return err
	}

	return nil
}
