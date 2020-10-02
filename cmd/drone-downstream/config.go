// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package main

import (
	"time"

	"github.com/drone-plugins/drone-downstream/plugin"
	"github.com/urfave/cli/v2"
)

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
			Name:        "last-successful",
			Usage:       "Deploy last successful build",
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
			Usage:       "Environment to trigger deploy to",
			EnvVars:     []string{"PLUGIN_DEPLOY"},
			Destination: &settings.Deploy,
		},
		&cli.BoolFlag{
			Name:        "block",
			Usage:       "Block until the triggered build is finished, makes this build fail if triggered build fails",
			EnvVars:     []string{"PLUGIN_BLOCK"},
			Destination: &settings.Block,
		},
		&cli.DurationFlag{
			Name:        "timeout",
			Value:       time.Duration(60) * time.Minute,
			Usage:       "How long to block until the triggered build is finished",
			EnvVars:     []string{"PLUGIN_BLOCK_TIMEOUT"},
			Destination: &settings.Timeout,
		},
	}
}
