# drone-curseforge

[![Build Status](http://cloud.drone.io/api/badges/drone-plugins/drone-curseforge/status.svg)](http://cloud.drone.io/drone-plugins/drone-curseforge)
[![Gitter chat](https://badges.gitter.im/drone/drone.png)](https://gitter.im/drone/drone)
[![Join the discussion at https://discourse.drone.io](https://img.shields.io/badge/discourse-forum-orange.svg)](https://discourse.drone.io)
[![Drone questions at https://stackoverflow.com](https://img.shields.io/badge/drone-stackoverflow-orange.svg)](https://stackoverflow.com/questions/tagged/drone.io)
[![](https://images.microbadger.com/badges/image/plugins/curseforge.svg)](https://microbadger.com/images/plugins/curseforge "Get your own image badge on microbadger.com")
[![Go Doc](https://godoc.org/github.com/drone-plugins/drone-curseforge?status.svg)](http://godoc.org/github.com/drone-plugins/drone-curseforge)
[![Go Report](https://goreportcard.com/badge/github.com/drone-plugins/drone-curseforge)](https://goreportcard.com/report/github.com/drone-plugins/drone-curseforge)

Drone plugin to publish releases to CurseForge. For the usage information and a listing of the available options please take a look at [the docs](http://plugins.drone.io/drone-plugins/drone-curseforge/).

## Build

Build the binary with the following command:

```console
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0
export GO111MODULE=on

go build -v -a -tags netgo -o release/linux/amd64/drone-curseforge ./cmd/drone-curseforge
```

## Docker

Build the Docker image with the following command:

```console
docker build \
  --label org.label-schema.build-date=$(date -u +"%Y-%m-%dT%H:%M:%SZ") \
  --label org.label-schema.vcs-ref=$(git rev-parse --short HEAD) \
  --file docker/Dockerfile.linux.amd64 --tag plugins/curseforge .
```

## Usage

```
docker run --rm \
  -e PLUGIN_API_KEY=574d1718-12e7-4ea8-8954-a08cd05c98e1 \
  -e PLUGIN_PROJECT=octopack \
  -e PLUGIN_FILE=dist/octopack-1.0.0.zip \
  -e PLUGIN_TITLE=octoPack \
  -e PLUGIN_NOTE=dist/CHANGELOG.md \
  -e PLUGIN_GAMES=6756,6757,6758 \
  -v $(pwd):$(pwd) \
  -w $(pwd) \
  plugins/curseforge
```
