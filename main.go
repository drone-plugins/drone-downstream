package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin"
)

// Params stores the git clone parameters used to
// configure and customzie the git clone behavior.
type Params struct {
	Repos            []string `json:"repositories"`
	Server           string   `json:"server"`
	Token            string   `json:"token"`
	Fork             bool     `json:"fork"`
	DisableSslVerify bool     `json:"disable_ssl_verify"`
}

var (
	buildCommit string
)

func main() {
	fmt.Printf("Drone Downstream Plugin built from %s\n", buildCommit)

	v := new(Params)
	s := new(drone.System)
	plugin.Param("system", s)
	plugin.Param("vargs", v)
	plugin.MustParse()

	// if no server url is provided we can default
	// to the hosted Drone service.
	if len(v.Token) == 0 {
		fmt.Println("Error: you must provide your Drone access token.")
		os.Exit(1)
	}

	if v.Server == "" {
		v.Server = s.Link
	}

	// create the drone client
	client := drone.NewClientToken(v.Server, v.Token)
	if v.DisableSslVerify {
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		client = drone.NewClientTokenTLS(v.Server, v.Token, tlsConfig)
	}

	for _, entry := range v.Repos {

		// parses the repository name in owner/name@branch format
		owner, name, branch := parseRepoBranch(entry)
		if len(owner) == 0 || len(name) == 0 {
			fmt.Printf("Error: unable to parse repository name %s.\n", entry)
			os.Exit(1)
		}
		if v.Fork {
			// get the latest build for the specified repository
			build, err := client.BuildLast(owner, name, branch)
			if err != nil {
				fmt.Printf("Error: unable to get latest build for %s.\n", entry)
				os.Exit(1)
			}
			// start a new  build
			_, err = client.BuildFork(owner, name, build.Number)
			if err != nil {
				fmt.Printf("Error: unable to trigger a new build for %s.\n", entry)
				os.Exit(1)
			}

			fmt.Printf("Starting new build %d for %s\n", build.Number, entry)

		} else {
			// get the latest build for the specified repository
			build, err := client.BuildLast(owner, name, branch)
			if err != nil {
				fmt.Printf("Error: unable to get latest build for %s.\n", entry)
				os.Exit(1)
			}

			// rebuild the latest build
			_, err = client.BuildStart(owner, name, build.Number)
			if err != nil {
				fmt.Printf("Error: unable to trigger build for %s.\n", entry)
				os.Exit(1)
			}

			fmt.Printf("Restarting build %d for %s\n", build.Number, entry)
		}
	}
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
