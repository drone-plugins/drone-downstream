// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package main

import (
	"github.com/grafana/drone-downstream/plugin"
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
	}
}
