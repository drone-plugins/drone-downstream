package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/drone/drone-go/drone"
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

	client := drone.NewClientToken(p.Server, p.Token)

	for _, entry := range p.Repos {

		// parses the repository name in owner/name@branch format
		owner, name, branch := parseRepoBranch(entry)
		if len(owner) == 0 || len(name) == 0 {
			return fmt.Errorf("Error: unable to parse repository name %s.\n", entry)
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
				} else if p.LastSuccessful && build.Status != drone.StatusSuccess {
					builds, err := client.BuildList(owner, name)
					if err != nil {
						return fmt.Errorf("Error: unable to get build list for %s.\n", entry)
					}

					build = nil
					for _, b := range builds {
						if b.Branch == branch && b.Status == drone.StatusSuccess {
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
						_, err = client.BuildFork(owner, name, build.Number)
						if err != nil {
							if waiting {
								continue
							}
							return fmt.Errorf("Error: unable to trigger a new build for %s.\n", entry)
						}
						fmt.Printf("Starting new build %d for %s.\n", build.Number, entry)
						break I
					} else {
						// rebuild the latest build
						_, err = client.BuildStart(owner, name, build.Number)
						if err != nil {
							if waiting {
								continue
							}
							return fmt.Errorf("Error: unable to trigger build for %s.\n", entry)
						}
						fmt.Printf("Restarting build %d for %s\n", build.Number, entry)
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
