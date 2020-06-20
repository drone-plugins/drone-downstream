// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/drone/drone-go/drone"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
	"golang.org/x/oauth2"
)

// Settings for the plugin.
type Settings struct {
	Repos          cli.StringSlice
	Server         string
	Token          string
	Wait           bool
	Timeout        time.Duration
	LastSuccessful bool
	Params         cli.StringSlice
	ParamsEnv      cli.StringSlice
	Deploy         string

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

	if p.settings.Wait && p.settings.LastSuccessful {
		return fmt.Errorf("only one of wait and last_successful can be true; choose one")
	}

	var err error
	p.settings.params, err = parseParams(p.settings.Params.Value())
	if err != nil {
		return fmt.Errorf("unable to parse params: %s", err)
	}

	for _, k := range p.settings.ParamsEnv.Value() {
		v, exists := os.LookupEnv(k)
		if !exists {
			return fmt.Errorf("param_from_env %s is not set", k)
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

	for _, entry := range p.settings.Repos.Value() {

		// parses the repository name in owner/name@branch format
		owner, name, branch := parseRepoBranch(entry)
		if len(owner) == 0 || len(name) == 0 {
			return fmt.Errorf("unable to parse repository name %s", entry)
		}

		// check for mandatory build no during deploy trigger
		if len(p.settings.Deploy) != 0 {
			if branch == "" {
				return fmt.Errorf("build no or branch must be mentioned for deploy, format repository@build/branch")
			}
			if _, err := strconv.Atoi(branch); err != nil && !p.settings.LastSuccessful {
				return fmt.Errorf("for deploy build no must be numeric only " +
					" or for branch deploy last_successful should be true," +
					" format repository@build/branch")
			}
		}

		waiting := false

		timeout := time.After(p.settings.Timeout)
		//lint:ignore SA1015 refactor later
		tick := time.Tick(1 * time.Second)

		var err error

		// Keep trying until we're timed out, successful or got an error
		// Tagged with "I" due to break nested in select
	I:
		for {
			select {
			// Got a timeout! fail with a timeout error
			case <-timeout:
				return fmt.Errorf("timed out waiting on a build for %s", entry)
			// Got a tick, we should check on the build status
			case <-tick:
				// first handle the deploy trigger
				if len(p.settings.Deploy) != 0 {
					var build *drone.Build
					if p.settings.LastSuccessful {
						// Get the last successful build of branch
						builds, err := client.BuildList(owner, name, drone.ListOptions{})
						if err != nil {
							return fmt.Errorf("unable to get build list for %s", entry)
						}

						for _, b := range builds {
							if b.Source == branch && b.Status == drone.StatusPassing {
								build = b
								break
							}
						}
						if build == nil {
							return fmt.Errorf("unable to get last successful build for %s", entry)
						}
					} else {
						// Get build by number
						buildNumber, _ := strconv.Atoi(branch)
						build, err = client.Build(owner, name, buildNumber)
						if err != nil {
							return fmt.Errorf("unable to get requested build %v for deploy for %s", buildNumber, entry)
						}
					}
					if p.settings.Wait && !waiting && (build.Status == drone.StatusRunning || build.Status == drone.StatusPending) {
						fmt.Printf("BuildLast for repository: %s, returned build number: %v with a status of %s. Will retry for %v.\n", entry, build.Number, build.Status, p.settings.Timeout)
						waiting = true
						continue
					}
					if (build.Status != drone.StatusRunning && build.Status != drone.StatusPending) || !p.settings.Wait {
						// start a new deploy
						_, err = client.Promote(owner, name, int(build.Number), p.settings.Deploy, p.settings.params)
						if err != nil {
							if waiting {
								continue
							}
							return fmt.Errorf("unable to trigger deploy for %s - err %v", entry, err)
						}
						fmt.Printf("starting deploy for %s/%s env - %s build - %d", owner, name, p.settings.Deploy, build.Number)
						logParams(p.settings.params, p.settings.ParamsEnv.Value())
						break I
					}
				}

				// get the latest build for the specified repository
				build, err := client.BuildLast(owner, name, branch)
				if err != nil {
					if waiting {
						continue
					}
					return fmt.Errorf("unable to get latest build for %s: %s", entry, err)
				}
				if p.settings.Wait && !waiting && (build.Status == drone.StatusRunning || build.Status == drone.StatusPending) {
					fmt.Printf("BuildLast for repository: %s, returned build number: %v with a status of %s. Will retry for %v.\n", entry, build.Number, build.Status, p.settings.Timeout)
					waiting = true
					continue
				} else if p.settings.LastSuccessful && build.Status != drone.StatusPassing {
					builds, err := client.BuildList(owner, name, drone.ListOptions{})
					if err != nil {
						return fmt.Errorf("unable to get build list for %s", entry)
					}

					build = nil
					for _, b := range builds {
						if b.Source == branch && b.Status == drone.StatusPassing {
							build = b
							break
						}
					}
					if build == nil {
						return fmt.Errorf("unable to get last successful build for %s", entry)
					}
				}

				if (build.Status != drone.StatusRunning && build.Status != drone.StatusPending) || !p.settings.Wait {
					// rebuild the latest build
					_, err = client.BuildRestart(owner, name, int(build.Number), p.settings.params)
					if err != nil {
						if waiting {
							continue
						}
						return fmt.Errorf("unable to trigger build for %s", entry)
					}
					fmt.Printf("Restarting build %d for %s\n", build.Number, entry)
					logParams(p.settings.params, p.settings.ParamsEnv.Value())

					break I
				}
			}
		}
	}

	return nil
}

func parseRepoBranch(repo string) (string, string, string) {
	var (
		owner  string
		name   string
		branch string
	)

	parts := strings.Split(repo, "@")
	if len(parts) == 2 {
		branch = parts[1]
		repo = parts[0]
	}

	parts = strings.Split(repo, "/")
	if len(parts) == 2 {
		owner = parts[0]
		name = parts[1]
	}
	return owner, name, branch
}

func parseParams(paramList []string) (map[string]string, error) {
	params := make(map[string]string)
	for _, p := range paramList {
		parts := strings.SplitN(p, "=", 2)
		if len(parts) == 2 {
			params[parts[0]] = parts[1]
		} else if _, err := os.Stat(parts[0]); os.IsNotExist(err) {
			return nil, fmt.Errorf(
				"invalid param '%s'; must be KEY=VALUE or file path",
				parts[0],
			)
		} else {
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
