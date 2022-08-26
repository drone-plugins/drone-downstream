# drone-downstream

[![Build Status](http://cloud.drone.io/api/badges/drone-plugins/drone-downstream/status.svg)](http://cloud.drone.io/drone-plugins/drone-downstream)
[![Gitter chat](https://badges.gitter.im/drone/drone.png)](https://gitter.im/drone/drone)
[![Join the discussion at https://community.harness.io/](https://img.shields.io/badge/discourse-forum-orange.svg)](https://community.harness.io/)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![](https://images.microbadger.com/badges/image/plugins/downstream.svg)](https://microbadger.com/images/plugins/downstream "Get your own image badge on microbadger.com")
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-downstream?status.svg)](http://godoc.org/github.com/drone-plugins/drone-downstream)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-downstream)](https://goreportcard.com/report/github.com/drone-plugins/drone-downstream)

Drone plugin to trigger downstream repository builds. For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-downstream/).

## Build

Build the binary with the following command:

```console
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o release/linux/amd64/drone-downstream
```

## Docker

Build the Docker image with the following command:

```console
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --file docker/Dockerfile.linux.amd64 --tag plugins/downstream .
```

## Usage

```console
docker run --rm \
  -e PLUGIN_REPOSITORIES=octocat/Hello-World \
  -e PLUGIN_TOKEN=eyJhbFciHiJISzI1EiIsUnR5cCW6IkpXQCJ9.ezH0ZXh0LjoidGJvZXJnZXIiLCJ0eXBlIjoidXNlciJ9.1m_3QFA6eA7h4wrBby2aIRFAEhQWPrlj4dsO_Gfchtc \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/downstream
```
