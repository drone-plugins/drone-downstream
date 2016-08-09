package main

import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/urfave/cli"
)

var version string // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "Drone downstream plugin"
	app.Usage = "drone downstream plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{

		cli.StringSliceFlag{
			Name:   "repositories",
			Usage:  "List of repositories to trigger",
			EnvVar: "PLUGIN_REPOSITORIES",
		},
		cli.StringFlag{
			Name:   "server",
			Usage:  "Trigger a drone build on a custom server",
			EnvVar: "PLUGIN_SERVER",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "Drone API token from your user settings",
			EnvVar: "DOWNSTREAM_TOKEN,PLUGIN_TOKEN",
		},
		cli.BoolFlag{
			Name:   "fork",
			Usage:  "Trigger a new build for a repository",
			EnvVar: "PLUGIN_FORK",
		},
	}

	app.Run(os.Args)
}

func run(c *cli.Context) error {

	plugin := Plugin{
		Repos:  c.StringSlice("repositories"),
		Server: c.String("server"),
		Token:  c.String("token"),
		Fork:   c.Bool("fork"),
	}

	if err := plugin.Exec(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}
