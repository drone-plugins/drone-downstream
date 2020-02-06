package main

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	version = "unknown"
)

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
			EnvVar: "DOWNSTREAM_SERVER,PLUGIN_SERVER",
		},
		cli.StringFlag{
			Name:   "system.host",
			Usage:  "Host for default value of server flag",
			EnvVar: "DRONE_SYSTEM_HOST,PLUGIN_HOST",
		},
		cli.StringFlag{
			Name:   "system.proto",
			Usage:  "Protocol for default value of server flag",
			EnvVar: "DRONE_SYSTEM_PROTO,PLUGIN_PROTO",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "Drone API token from your user settings",
			EnvVar: "DRONE_TOKEN,DOWNSTREAM_TOKEN,PLUGIN_TOKEN",
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
		cli.BoolFlag{
			Name:   "block",
			Usage:  "block on completion of downstream build",
			EnvVar: "PLUGIN_BLOCK",
		},
		cli.DurationFlag{
			Name:   "block-timeout",
			Value:  time.Duration(60) * time.Second,
			Usage:  "How long to wait on blocking downstream build",
			EnvVar: "PLUGIN_BLOCK_TIMEOUT",
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
		Host:           c.String("system.host"),
		Proto:          c.String("system.proto"),
		Token:          c.String("token"),
		Wait:           c.Bool("wait"),
		Timeout:        c.Duration("timeout"),
		LastSuccessful: c.Bool("last-successful"),
		Block:          c.Bool("block"),
		BlockTimeout:   c.Duration("block-timeout"),
		Params:         c.StringSlice("params"),
		ParamsEnv:      c.StringSlice("params-from-env"),
		Deploy:         c.String("deploy"),
	}

	return plugin.Exec()
}
