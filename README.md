# drone-downstream

[![Build Status](http://cloud.drone.io/api/badges/grafana/drone-downstream/status.svg)](http://cloud.drone.io/grafana/drone-downstream)
[![Go.dev](https://pkg.go.dev/badge/github.com/grafana/drone-downstream)](https://pkg.go.dev/github.com/grafana/drone-downstream)
[![Go Report](https://goreportcard.com/badge/github.com/grafana/drone-downstream)](https://goreportcard.com/report/github.com/grafana/drone-downstream)

Drone plugin to trigger downstream repository builds. For the usage information and a listing of the available options,
please take a look at [the docs](https://pkg.go.dev/github.com/grafana/drone-downstream).

## Build

Build the binary with the following command:

```console
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -o release/linux/amd64/drone-downstream ./cmd/drone-downstream
```

## Docker

Build the Docker image with the following command:

```console
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --file docker/Dockerfile.linux.amd64 --tag grafana/drone-downstream .
```

## Usage

```console
docker run --rm \
  -e PLUGIN_REPOSITORIES=octocat/Hello-World \
  -e PLUGIN_TOKEN=eyJhbFciHiJISzI1EiIsUnR5cCW6IkpXQCJ9.ezH0ZXh0LjoidGJvZXJnZXIiLCJ0eXBlIjoidXNlciJ9.1m_3QFA6eA7h4wrBby2aIRFAEhQWPrlj4dsO_Gfchtc \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  grafana/drone-downstream
```
