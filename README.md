# drone-downstream

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-downstream/status.svg)](http://beta.drone.io/drone-plugins/drone-downstream)
[![Join the chat at https://gitter.im/drone/drone](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/drone/drone)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-downstream?status.svg)](http://godoc.org/github.com/drone-plugins/drone-downstream)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-downstream)](https://goreportcard.com/report/github.com/drone-plugins/drone-downstream)
[![](https://images.microbadger.com/badges/image/plugins/downstream.svg)](https://microbadger.com/images/plugins/downstream "Get your own image badge on microbadger.com")

Drone plugin to trigger downstream repository builds. For the usage information and a listing of the available options please take a look at [the docs](DOCS.md).

## Build

Build the binary with the following commands:

```
go build
```

## Docker

Build the Docker image with the following commands:

```
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -tags netgo -o release/linux/amd64/drone-downstream
docker build --rm -t plugins/downstream .
```

## Usage

Execute from the working directory:

```sh
docker run --rm \
  -e PLUGIN_REPOSITORIES=octocat/Hello-World \
  -e PLUGIN_TOKEN=eyJhbFciHiJISzI1EiIsUnR5cCW6IkpXQCJ9.ezH0ZXh0LjoidGJvZXJnZXIiLCJ0eXBlIjoidXNlciJ9.1m_3QFA6eA7h4wrBby2aIRFAEhQWPrlj4dsO_Gfchtc \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/downstream
```
