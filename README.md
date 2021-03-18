# Go Start!

With [skaffold](https://skaffold.dev/) (distroless 12.6MB):

```
skaffold run -f ci/skaffold.yaml
kubectl port-forward service/gomessage 8080:8080
```

With [buildpacks](https://buildpacks.io/) ([OCI](https://opencontainers.org/) 30.7MB):

```
pack build rmeharg/gomessage --builder paketobuildpacks/builder:tiny
docker run rmeharg/gomessage
```

### Local Development

```
go get github.com/azer/yolo
make help
```

