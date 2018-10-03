package main

import (
	"reflect"
	"testing"
)

func Test_parseRepoBranch(t *testing.T) {

	var tests = []struct {
		Repo   string
		Owner  string
		Name   string
		Branch string
		Env    string
	}{
		{"octocat/hello-world", "octocat", "hello-world", "", ""},
		{"octocat/hello-world@master", "octocat", "hello-world", "master", ""},
		{"octocat/hello-world@master@production", "octocat", "hello-world", "master", "production"},
	}

	for _, test := range tests {

		owner, name, branch, env := parseRepoBranch(test.Repo)
		if owner != test.Owner {
			t.Errorf("wanted repository owner %s, got %s", test.Owner, owner)
		}
		if name != test.Name {
			t.Errorf("wanted repository name %s, got %s", test.Name, name)
		}
		if branch != test.Branch {
			t.Errorf("wanted repository branch %s, got %s", test.Branch, branch)
		}
		if env != test.Env {
			t.Errorf("wanted deployment event %s, got %s", test.Env, env)
		}
	}
}

func Test_parseParams_invalid(t *testing.T) {
	out, err := parseParams([]string{"invalid"})
	if err == nil {
		t.Errorf("expected error, got %v", out)
	}
}

func Test_parseParams(t *testing.T) {
	var tests = []struct {
		Input  []string
		Output map[string]string
	}{
		{[]string{}, map[string]string{}},
		{
			[]string{"where=far", "who=you"},
			map[string]string{"where": "far", "who": "you"},
		},
		{
			[]string{"where=very=far"},
			map[string]string{"where": "very=far"}},
		{
			[]string{"test_params.env"},
			map[string]string{
				"SOME_VAR": "someval",
				"FOO":      "BAR",
				"BAR":      "BAZ",
				"foo":      "bar",
				"bar":      "baz",
			},
		},
		{
			[]string{"test_params.env", "where=far", "who=you"},
			map[string]string{
				"SOME_VAR": "someval",
				"FOO":      "BAR",
				"BAR":      "BAZ",
				"foo":      "bar",
				"bar":      "baz",
				"where":    "far",
				"who":      "you",
			},
		},
	}

	for _, test := range tests {
		out, err := parseParams(test.Input)
		if err != nil {
			t.Errorf("unable to parse params: %s", err)

			break
		}

		if !reflect.DeepEqual(out, test.Output) {
			t.Errorf("wanted params %+v, got %+v", test.Output, out)
		}
	}
}
