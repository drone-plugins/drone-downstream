// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/drone/drone-go/drone"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

// Settings for the plugin.
type Settings struct {
	Repos        cli.StringSlice
	Server       string
	Token        string
	Params       cli.StringSlice
	ParamsEnv    cli.StringSlice
	Block        bool
	BlockTimeout time.Duration

	server string
	params map[string]string
}

// Validate handles the settings validation of the plugin.
func (p *Plugin) Validate() error {
	if len(p.settings.Token) == 0 {
		return fmt.Errorf("you must provide your drone access token")
	}

	p.settings.server = getServerWithDefaults(p.settings.Server, p.pipeline.System.Host, p.pipeline.System.Proto)
	if len(p.settings.server) == 0 {
		return fmt.Errorf("you must provide your drone server")
	}

	var err error
	p.settings.params, err = parseParams(p.settings.Params.Value())
	if err != nil {
		return fmt.Errorf("unable to parse params: %s", err)
	}

	upstreamBuildNumber, ok := os.LookupEnv("DRONE_BUILD_NUMBER")
	if ok {
		p.settings.params["DRONE_UPSTREAM_BUILD_NUMBER"] = upstreamBuildNumber
	}

	for _, k := range p.settings.ParamsEnv.Value() {
		v, exists := os.LookupEnv(k)
		if !exists {
			return fmt.Errorf("param_from_env '%s' is not set", k)
		}

		p.settings.params[k] = v
	}

	return nil
}

// Execute provides the implementation of the plugin.
func (p *Plugin) Execute() error {
	config := new(oauth2.Config)

	auther := config.Client(
		context.WithValue(context.Background(), oauth2.HTTPClient, p.network.Client),
		&oauth2.Token{
			AccessToken: p.settings.Token,
		},
	)

	client := drone.NewClient(p.settings.server, auther)

	// For each configured repo, on the format <owner>/<name>@<branch>
	for _, entry := range p.settings.Repos.Value() {
		// Parse the repository name in owner/name@branch format
		owner, name, branch, err := parseRepoBranch(entry)
		if err != nil {
			return err
		}

		build, err := client.BuildCreate(owner, name, "", branch, p.settings.params)
		if err != nil {
			return fmt.Errorf("failed to create Drone build for %s/%s: %w", owner, name, err)
		}

		fmt.Printf("successfully started Drone build for %s/%s: #%d\n", owner, name, build.ID)
		logParams(p.settings.params, p.settings.ParamsEnv.Value())
	}

	return nil
}

func parseRepoBranch(repo string) (string, string, string, error) {
	parts := strings.Split(repo, "@")
	var branch string
	if len(parts) == 2 {
		branch = parts[1]
		repo = parts[0]
	}

	parts = strings.Split(repo, "/")
	var name string
	var owner string
	if len(parts) == 2 {
		owner = parts[0]
		name = parts[1]
	}
	if owner == "" || name == "" {
		return "", "", "", fmt.Errorf("unable to parse repository name %q", repo)
	}

	return owner, name, branch, nil
}

func parseParams(paramList []string) (map[string]string, error) {
	params := make(map[string]string)
	for _, p := range paramList {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) == 2 {
			params[parts[0]] = parts[1]
		} else {
			if _, err := os.Stat(parts[0]); err != nil {
				return nil, fmt.Errorf(
					"invalid param '%s'; must be KEY=VALUE or file path",
					parts[0],
				)
			}

			fileParams, err := godotenv.Read(parts[0])
			if err != nil {
				return nil, err
			}

			for k, v := range fileParams {
				params[k] = v
			}
		}
	}

	return params, nil
}

func logParams(params map[string]string, paramsEnv []string) {
	if len(params) > 0 {
		fmt.Println("  with params:")
		for k, v := range params {
			fromEnv := false
			for _, e := range paramsEnv {
				if k == e {
					fromEnv = true
					break
				}
			}
			if fromEnv {
				v = "[from-environment]"
			}
			fmt.Printf("  - %s: %s\n", k, v)
		}
	}
}

func getServerWithDefaults(server string, host string, protocol string) string {
	if len(server) != 0 {
		return server
	}

	if len(host) == 0 || len(protocol) == 0 {
		return ""
	}

	return fmt.Sprintf("%s://%s", protocol, host)
}

func blockUntilBuildIsFinished(p *Plugin, client drone.Client, namespace, name string, buildNumber int) error {
	fmt.Printf("\nblocking until triggered build is finished\n")

	timeout := time.After(p.settings.BlockTimeout)

	//lint:ignore SA1015 refactor later
	tick := time.Tick(10 * time.Second)

	// listen for SIGINT and SIGTERM to cancel downstream build when stopping this executable
	// this does not work in drone because drone uses SIGKILL to terminate its containers
	// but when running the plugin locally during development, it's very handy
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	defer close(sigs)

	for {
		select {
		case <-sigs:
			err := client.BuildCancel(namespace, name, int(buildNumber))
			if err != nil {
				return fmt.Errorf("could not cancel downstream job %d", buildNumber)
			}

			fmt.Printf("canceled downstream job %d\n", buildNumber)

			return nil

		// Got a timeout! fail with a timeout error
		case <-timeout:
			return fmt.Errorf("timed out waiting for %d", buildNumber)

		// Got a tick, we should check on the build status
		case <-tick:
			build, err := client.Build(namespace, name, buildNumber)
			if err != nil {
				return err
			}

			switch build.Status {
			case drone.StatusError, drone.StatusKilled, drone.StatusFailing, drone.StatusDeclined, drone.StatusSkipped:
				return fmt.Errorf(
					"build %d did not succeed: %s",
					buildNumber,
					build.Status,
				)
			case drone.StatusPassing:
				return nil
			default:
				fmt.Printf("Waiting for build %d in status %s\n", buildNumber, build.Status)
			}
		}
	}
}
