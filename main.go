// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

// DO NOT MODIFY THIS FILE DIRECTLY

package main

import (
	"os"
	"time"

	"github.com/drone-plugins/drone-downstream/plugin"
	"github.com/drone-plugins/drone-plugin-lib/errors"
	"github.com/drone-plugins/drone-plugin-lib/urfave"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

var version = "unknown"

func main() {
	settings := &plugin.Settings{}

	if _, err := os.Stat("/run/drone/env"); err == nil {
		_ = godotenv.Overload("/run/drone/env")
	}

	app := &cli.App{
		Name:                      "drone-downstream",
		Usage:                     "trigger a downstream drone build",
		Version:                   version,
		Flags:                     append(settingsFlags(settings), urfave.Flags()...),
		Action:                    run(settings),
		DisableSliceFlagSeparator: true,
	}

	if err := app.Run(os.Args); err != nil {
		errors.HandleExit(err)
	}
}

func run(settings *plugin.Settings) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		urfave.LoggingFromContext(ctx)

		plugin := plugin.New(
			*settings,
			urfave.PipelineFromContext(ctx),
			urfave.NetworkFromContext(ctx),
		)

		if err := plugin.Validate(); err != nil {
			if e, ok := err.(errors.ExitCoder); ok {
				return e
			}

			return errors.ExitMessagef("validation failed: %w", err)
		}

		if err := plugin.Execute(); err != nil {
			if e, ok := err.(errors.ExitCoder); ok {
				return e
			}

			return errors.ExitMessagef("execution failed: %w", err)
		}

		return nil
	}
}

// settingsFlags has the cli.Flags for the plugin.Settings.
func settingsFlags(settings *plugin.Settings) []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:        "repositories",
			Usage:       "List of repositories to trigger",
			EnvVars:     []string{"PLUGIN_REPOSITORIES"},
			Destination: &settings.Repos,
		},
		&cli.StringFlag{
			Name:        "server",
			Usage:       "Trigger a drone build on a custom server",
			EnvVars:     []string{"PLUGIN_SERVER", "DOWNSTREAM_SERVER"},
			Destination: &settings.Server,
		},
		&cli.StringFlag{
			Name:        "token",
			Usage:       "Drone API token from your user settings",
			EnvVars:     []string{"PLUGIN_TOKEN", "DRONE_TOKEN", "DOWNSTREAM_TOKEN"},
			Destination: &settings.Token,
		},
		&cli.BoolFlag{
			Name:        "wait",
			Usage:       "Wait for any currently running builds to finish",
			EnvVars:     []string{"PLUGIN_WAIT"},
			Destination: &settings.Wait,
		},
		&cli.DurationFlag{
			Name:        "timeout",
			Value:       time.Duration(60) * time.Second,
			Usage:       "How long to wait on any currently running builds",
			EnvVars:     []string{"PLUGIN_WAIT_TIMEOUT"},
			Destination: &settings.Timeout,
		},
		&cli.BoolFlag{
			Name:        "last-successful",
			Usage:       "Trigger last successful build",
			EnvVars:     []string{"PLUGIN_LAST_SUCCESSFUL"},
			Destination: &settings.LastSuccessful,
		},
		&cli.StringSliceFlag{
			Name:        "params",
			Usage:       "List of params (key=value or file paths of params) to pass to triggered builds",
			EnvVars:     []string{"PLUGIN_PARAMS"},
			Destination: &settings.Params,
		},
		&cli.StringSliceFlag{
			Name:        "params-from-env",
			Usage:       "List of environment variables to pass to triggered builds",
			EnvVars:     []string{"PLUGIN_PARAMS_FROM_ENV"},
			Destination: &settings.ParamsEnv,
		},
		&cli.StringFlag{
			Name:        "deploy",
			Usage:       "Environment to trigger deploy for the respective build",
			EnvVars:     []string{"PLUGIN_DEPLOY"},
			Destination: &settings.Deploy,
		},
	}
}
