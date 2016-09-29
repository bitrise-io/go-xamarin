package cli

import "github.com/urfave/cli"

const (
	solutionFilePathKey      string = "path"
	solutionConfigurationKey string = "configuration"
	solutionPlatformKey      string = "platform"

	forceMDToolKey string = "force-mdtool"
)

var commands = []cli.Command{
	{
		Name:   "build",
		Usage:  "Build xamarin projects",
		Action: buildCmd,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  solutionFilePathKey,
				Usage: "Solution file path",
			},
			cli.StringFlag{
				Name:  solutionConfigurationKey,
				Usage: "Solution configuration",
			},
			cli.StringFlag{
				Name:  solutionPlatformKey,
				Usage: "Solution platform",
			},
			cli.BoolFlag{
				Name:  forceMDToolKey,
				Usage: "Force use mdtool",
			},
		},
	},
	{
		Name:   "clean",
		Usage:  "Clean xamarin projects",
		Action: cleanCmd,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  solutionFilePathKey,
				Usage: "Solution file path",
			},
		},
	},
	{
		Name:   "version",
		Usage:  "Prints version",
		Action: versionCmd,
	},
}
