package main

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	version = "0.0.0"
	build   = "0"
)

func main() {
	app := cli.NewApp()
	app.Name = "downstream plugin"
	app.Usage = "downstream plugin"
	app.Version = fmt.Sprintf("%s+%s", version, build)
	app.Action = run
	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:   "repositories",
			Usage:  "List of repositories to trigger",
			EnvVar: "PLUGIN_REPOSITORIES",
		},
		cli.StringFlag{
			Name:   "server",
			Usage:  "Trigger a drone build on a custom server",
			EnvVar: "DOWNSTREAM_SERVER,PLUGIN_SERVER",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "Drone API token from your user settings",
			EnvVar: "DRONE_TOKEN,DOWNSTREAM_TOKEN,PLUGIN_TOKEN",
		},
		cli.BoolFlag{
			Name:   "fork",
			Usage:  "Trigger a new build for a repository",
			EnvVar: "PLUGIN_FORK",
		},
		cli.BoolFlag{
			Name:   "wait",
			Usage:  "Wait for any currently running builds to finish",
			EnvVar: "PLUGIN_WAIT",
		},
		cli.DurationFlag{
			Name:   "timeout",
			Value:  time.Duration(60) * time.Second,
			Usage:  "How long to wait on any currently running builds",
			EnvVar: "PLUGIN_WAIT_TIMEOUT",
		},
		cli.BoolFlag{
			Name:   "last-successful",
			Usage:  "Trigger last successful build",
			EnvVar: "PLUGIN_LAST_SUCCESSFUL",
		},
		cli.StringSliceFlag{
			Name:   "params",
			Usage:  "List of params (key=value or file paths of params) to pass to triggered builds",
			EnvVar: "PLUGIN_PARAMS",
		},
		cli.StringSliceFlag{
			Name:   "params-from-env",
			Usage:  "List of environment variables to pass to triggered builds",
			EnvVar: "PLUGIN_PARAMS_FROM_ENV",
		},
		cli.StringFlag{
			Name:   "deploy",
			Usage:  "Environment to trigger deploy for the respective build",
			EnvVar: "PLUGIN_DEPLOY",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := Plugin{
		Repos:          c.StringSlice("repositories"),
		Server:         c.String("server"),
		Token:          c.String("token"),
		Fork:           c.Bool("fork"),
		Wait:           c.Bool("wait"),
		Timeout:        c.Duration("timeout"),
		LastSuccessful: c.Bool("last-successful"),
		Params:         c.StringSlice("params"),
		ParamsEnv:      c.StringSlice("params-from-env"),
		Deploy:         c.String("deploy"),
	}

	return plugin.Exec()
}
