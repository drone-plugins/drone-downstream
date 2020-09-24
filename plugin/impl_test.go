// Copyright (c) 2020, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		owner, name, branch, err := parseRepoBranch(test.Repo)
		require.NoError(t, err)

		assert.Equal(t, test.Owner, owner)
		assert.Equal(t, test.Name, name)
		assert.Equal(t, test.Branch, branch)
	}
}

func Test_parseParams_invalid(t *testing.T) {
	out, err := parseParams([]string{"invalid"})
	assert.Error(t, err, fmt.Sprintf("Expected error, got %v", out))
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
		require.NoError(t, err)

		assert.Equal(t, test.Output, out)
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

		assert.Equal(t, test.Result, server)
	}
}
