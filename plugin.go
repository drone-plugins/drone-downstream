package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/drone/drone-go/drone"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

// Plugin defines the Downstream plugin parameters.
type Plugin struct {
	Repos          []string
	Server         string
	Token          string
	Fork           bool
	Wait           bool
	Timeout        time.Duration
	LastSuccessful bool
	Params         []string
	ParamsEnv      []string
	Deploy         string
}

// Exec runs the plugin
func (p *Plugin) Exec() error {
	if len(p.Token) == 0 {
		return fmt.Errorf("Error: you must provide your Drone access token.")
	}

	if len(p.Server) == 0 {
		return fmt.Errorf("Error: you must provide your Drone server.")
	}

	if p.Wait && p.LastSuccessful {
		return fmt.Errorf("Error: only one of wait and last_successful can be true; choose one")
	}

	if !p.Fork {
		fmt.Fprintln(
			os.Stderr,
			"Warning: \"fork: false\" will be deprecated in future\n"+
				"         set \"fork: true\" to disable this warning",
		)
	}

	params, err := parseParams(p.Params)
	if err != nil {
		return fmt.Errorf("Error: unable to parse params: %s.\n", err)
	}

	for _, k := range p.ParamsEnv {
		v, exists := os.LookupEnv(k)
		if !exists {
			return fmt.Errorf("Error: param_from_env %s is not set.\n", k)
		}

		params[k] = v
	}

	config := new(oauth2.Config)

	auther := config.Client(
		oauth2.NoContext,
		&oauth2.Token{
			AccessToken: p.Token,
		},
	)

	client := drone.NewClient(p.Server, auther)

	for _, entry := range p.Repos {

		// parses the repository name in owner/name@branch format
		owner, name, branch := parseRepoBranch(entry)
		if len(owner) == 0 || len(name) == 0 {
			return fmt.Errorf("Error: unable to parse repository name %s.\n", entry)
		}

		// check for mandatory build no during deploy trigger
		if len(p.Deploy) != 0 {
			if branch == "" {
				return fmt.Errorf("Error: build no or branch must be mentioned for deploy, format repository@build/branch")
			}
			if _, err := strconv.Atoi(branch); err != nil && !p.LastSuccessful {
				return fmt.Errorf("Error: for deploy build no must be numeric only " +
					" or for branch deploy last_successful should be true," +
					" format repository@build/branch")
			}
		}

		waiting := false

		timeout := time.After(p.Timeout)
		tick := time.Tick(1 * time.Second)

		// Keep trying until we're timed out, successful or got an error
		// Tagged with "I" due to break nested in select
	I:
		for {
			select {
			// Got a timeout! fail with a timeout error
			case <-timeout:
				return fmt.Errorf("Error: timed out waiting on a build for %s.\n", entry)
			// Got a tick, we should check on the build status
			case <-tick:
				// first handle the deploy trigger
				if len(p.Deploy) != 0 {
					var build *drone.Build
					if p.LastSuccessful {
						// Get the last successful build of branch
						builds, err := client.BuildList(owner, name, drone.ListOptions{})
						if err != nil {
							return fmt.Errorf("Error: unable to get build list for %s", entry)
						}

						for _, b := range builds {
							if b.Source == branch && b.Status == drone.StatusPassing {
								build = b
								break
							}
						}
						if build == nil {
							return fmt.Errorf("Error: unable to get last successful build for %s", entry)
						}
					} else {
						// Get build by number
						buildNumber, _ := strconv.Atoi(branch)
						build, err = client.Build(owner, name, buildNumber)
						if err != nil {
							return fmt.Errorf("Error: unable to get requested build %v for deploy for %s", buildNumber, entry)
						}
					}
					if p.Wait && !waiting && (build.Status == drone.StatusRunning || build.Status == drone.StatusPending) {
						fmt.Printf("BuildLast for repository: %s, returned build number: %v with a status of %s. Will retry for %v.\n", entry, build.Number, build.Status, p.Timeout)
						waiting = true
						continue
					}
					if (build.Status != drone.StatusRunning && build.Status != drone.StatusPending) || !p.Wait {
						// start a new deploy
						_, err = client.Promote(owner, name, int(build.Number), p.Deploy, params)
						if err != nil {
							if waiting {
								continue
							}
							return fmt.Errorf("Error: unable to trigger deploy for %s - err %v", entry, err)
						}
						fmt.Printf("Starting deploy for %s/%s env - %s build - %d.\n", owner, name, p.Deploy, build.Number)
						logParams(params, p.ParamsEnv)
						break I
					}
				}

				// get the latest build for the specified repository
				build, err := client.BuildLast(owner, name, branch)
				if err != nil {
					if waiting {
						continue
					}
					return fmt.Errorf("Error: unable to get latest build for %s.\n", entry)
				}
				if p.Wait && !waiting && (build.Status == drone.StatusRunning || build.Status == drone.StatusPending) {
					fmt.Printf("BuildLast for repository: %s, returned build number: %v with a status of %s. Will retry for %v.\n", entry, build.Number, build.Status, p.Timeout)
					waiting = true
					continue
				} else if p.LastSuccessful && build.Status != drone.StatusPassing {
					builds, err := client.BuildList(owner, name, drone.ListOptions{})
					if err != nil {
						return fmt.Errorf("Error: unable to get build list for %s.\n", entry)
					}

					build = nil
					for _, b := range builds {
						if b.Source == branch && b.Status == drone.StatusPassing {
							build = b
							break
						}
					}
					if build == nil {
						return fmt.Errorf("Error: unable to get last successful build for %s.\n", entry)
					}
				}

				if (build.Status != drone.StatusRunning && build.Status != drone.StatusPending) || !p.Wait {
					if p.Fork {
						// start a new  build
						_, err = client.BuildRestart(owner, name, int(build.Number), params)
						if err != nil {
							if waiting {
								continue
							}
							return fmt.Errorf("Error: unable to trigger a new build for %s.\n", entry)
						}
						fmt.Printf("Starting new build %d for %s.\n", build.Number, entry)
						logParams(params, p.ParamsEnv)
						break I
					} else {
						// rebuild the latest build
						_, err = client.BuildRestart(owner, name, int(build.Number), params)
						if err != nil {
							if waiting {
								continue
							}
							return fmt.Errorf("Error: unable to trigger build for %s.\n", entry)
						}
						fmt.Printf("Restarting build %d for %s\n", build.Number, entry)
						logParams(params, p.ParamsEnv)

						break I
					}
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
