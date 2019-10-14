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

func Test_getServerWithDefaults(t *testing.T) {
	var tests = []struct {
		Server string
		Host   string
		Proto  string
		Result string
	}{
		{"", "drone.example.com", "http", "http://drone.example.com"},
		{"", "drone.example.com:8000", "http", "http://drone.example.com:8000"},
		{"", "drone.example.com", "https", "https://drone.example.com"},
		{"", "drone.example.com:8888", "https", "https://drone.example.com:8888"},
		{"https://drone.example.com", "drone.example.com:8888", "https", "https://drone.example.com"},
	}

	for _, test := range tests {
		server := getServerWithDefaults(test.Server, test.Host, test.Proto)

		if server != test.Result {
			t.Errorf("wanted server url %s, got %s", test.Result, server)
		}
	}
}
