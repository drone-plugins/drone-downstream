package main

import (
	"fmt"
	"strings"

	"github.com/drone/drone-go/drone"
)

// Plugin defines the Downstream plugin parameters.
type Plugin struct {
	Repos  []string
	Server string
	Token  string
	Fork   bool
}

// Exec runs the plugin
func (p *Plugin) Exec() error {

	if len(p.Token) == 0 {
		return fmt.Errorf("Error: you must provide your Drone access token.")
	}

	if len(p.Server) == 0 {
		return fmt.Errorf("Error: you must provide your Drone server.")
	}

	client := drone.NewClientToken(p.Server, p.Token)

	for _, entry := range p.Repos {

		// parses the repository name in owner/name@branch format
		owner, name, branch := parseRepoBranch(entry)
		if len(owner) == 0 || len(name) == 0 {
			return fmt.Errorf("Error: unable to parse repository name %s.\n", entry)
		}
		if p.Fork {
			// get the latest build for the specified repository
			build, err := client.BuildLast(owner, name, branch)
			if err != nil {
				return fmt.Errorf("Error: unable to get latest build for %s.\n", entry)
			}
			// start a new  build
			_, err = client.BuildFork(owner, name, build.Number)
			if err != nil {
				return fmt.Errorf("Error: unable to trigger a new build for %s.\n", entry)
			}

			fmt.Printf("Starting new build %d for %s\n", build.Number, entry)

		} else {
			// get the latest build for the specified repository
			build, err := client.BuildLast(owner, name, branch)
			if err != nil {
				return fmt.Errorf("Error: unable to get latest build for %s.\n", entry)
			}

			// rebuild the latest build
			_, err = client.BuildStart(owner, name, build.Number)
			if err != nil {
				return fmt.Errorf("Error: unable to trigger build for %s.\n", entry)
			}

			fmt.Printf("Restarting build %d for %s\n", build.Number, entry)
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
