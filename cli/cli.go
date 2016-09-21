package cli

import (
	"os"

	"github.com/bitrise-io/go-utils/log"
	"github.com/urfave/cli"
)

// Run ...
func Run() {
	app := cli.NewApp()
	app.Name = "xamarin-builder"
	app.Usage = "Build xamarin projects"
	app.Version = "0.9.0"

	app.Commands = commands

	if err := app.Run(os.Args); err != nil {
		log.Error("Finished with error: %s", err)
		os.Exit(1)
	}
}
