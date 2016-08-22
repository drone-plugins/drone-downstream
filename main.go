package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/joho/godotenv"
)

var version string // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "downstream plugin"
	app.Usage = "downstream plugin"
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
		cli.StringFlag{
			Name:  "env-file",
			Usage: "source env file",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.String("env-file") != "" {
		_ = godotenv.Load(c.String("env-file"))
	}

	plugin := Plugin{
		Repos:  c.StringSlice("repositories"),
		Server: c.String("server"),
		Token:  c.String("token"),
		Fork:   c.Bool("fork"),
	}

	return plugin.Exec()
}
