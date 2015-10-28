package main

import (
	"testing"
)

func Test_parseRepoBranch(t *testing.T) {

	var tests = []struct {
		Repo   string
		Owner  string
		Name   string
		Branch string
	}{
		{"octocat/hello-world", "octocat", "hello-world", ""},
		{"octocat/hello-world@master", "octocat", "hello-world", "master"},
	}

	for _, test := range tests {

		owner, name, branch := parseRepoBranch(test.Repo)
		if owner != test.Owner {
			t.Errorf("wanted repository owner %s, got %s", test.Owner, owner)
		}
		if name != test.Name {
			t.Errorf("wanted repository name %s, got %s", test.Name, name)
		}
		if branch != test.Branch {
			t.Errorf("wanted repository branch %s, got %s", test.Branch, branch)
		}
	}
}
