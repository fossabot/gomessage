# Go Message!

[![Build Status](https://travis-ci.com/rmeharg/gomessage.svg)](https://travis-ci.com/rmeharg/gomessage)
[![codecov](https://codecov.io/gh/rmeharg/gomessage/branch/master/graph/badge.svg)](https://codecov.io/gh/rmeharg/gomessage)
[![Go Report Card](https://goreportcard.com/badge/github.com/rmeharg/gomessage)](https://goreportcard.com/report/github.com/rmeharg/gomessage)
[![GoDoc](https://godoc.org/github.com/rmeharg/gomessage?status.svg)](https://godoc.org/github.com/rmeharg/gomessage)
[![Release](https://img.shields.io/github/release/rmeharg/gomessage.svg?style=flat-square)](https://github.com/rmeharg/gomessage/releases)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/gin-gonic/gin)](https://www.tickgit.com/browse?repo=github.com/gin-gonic/gin)

With [skaffold](https://skaffold.dev/) (distroless 12.6MB):

```
skaffold run -f .ci/skaffold.yaml
```

With [buildpacks](https://buildpacks.io/) ([OCI](https://opencontainers.org/) 30.7MB):

```
pushd broadcaster && pack build broadcaster --builder paketobuildpacks/builder:tiny && popd
pushd listener && pack build listener --builder paketobuildpacks/builder:tiny && popd
pushd decoder && pack build decoder --builder paketobuildpacks/builder:tiny && popd
pushd writer && pack build writer --builder paketobuildpacks/builder:tiny && popd
pushd reporter && pack build reporter --builder paketobuildpacks/builder:tiny && popd
docker run broadcaster -multicast 239.0.0.0:9002 -count 100 -duplicate true
docker run listener -multicast 239.0.0.0:9002 -amqp "amqp://user:CHANGEME@rabbitmq.rabbitmq.svc.cluster.local:5672"
docker run decoder -amqp "amqp://user:CHANGEME@rabbitmq.rabbitmq.svc.cluster.local:5672"
docker run writer -amqp "amqp://user:CHANGEME@rabbitmq.rabbitmq.svc.cluster.local:5672" -influxdb "http://influxdb-enterprise-data.influxdb.svc.cluster.local:8086" -influxdb-authtoken "admin:admin"
docker run reporter -influxdb "http://influxdb-enterprise-data.influxdb.svc.cluster.local:8086" -influxdb-authtoken "admin:admin"
```

### Local Development

```
go get github.com/azer/yolo
make help
```

### Architecture

![alt text](./.ci/gomessage-architecture.png "GoMessage Architecture")
