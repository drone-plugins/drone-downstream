# drone-downstream

[![Build Status](http://beta.drone.io/api/badges/drone-plugins/drone-downstream/status.svg)](http://beta.drone.io/drone-plugins/drone-downstream)
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-downstream?status.svg)](http://godoc.org/github.com/drone-plugins/drone-downstream)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-downstream)](https://goreportcard.com/report/github.com/drone-plugins/drone-downstream)
[![Join the chat at https://gitter.im/drone/drone](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/drone/drone)

Drone plugin to trigger downstream repository builds. For the usage information
and a listing of the available options please take a look at
[the docs](DOCS.md).

## Build

Build the binary with the following commands:

```
go build
go test
```

## Docker

Build the docker image with the following commands:

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo
docker build --rm=true -t plugins/downstream .
```

Please note incorrectly building the image for the correct x64 linux and with
GCO disabled will result in an error when running the Docker image:

```
docker: Error response from daemon: Container command
'/bin/drone-downstream' not found or does not exist..
```

## Usage

Execute from the working directory:

```sh
docker run --rm \
  -e PLUGIN_REPOSITORIES=octocat/Hello-World \
  -e PLUGIN_TOKEN=eyJhbFciHiJISzI1EiIsUnR5cCW6IkpXQCJ9.ezH0ZXh0LjoidGJvZXJnZXIiLCJ0eXBlIjoidXNlciJ9.1m_3QFA6eA7h4wrBby2aIRFAEhQWPrlj4dsO_Gfchtc \
  -v $(pwd)/$(pwd) \
  -w $(pwd) \
  plugins/downstream
```
