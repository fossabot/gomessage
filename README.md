# Go Start!

With [skaffold](https://skaffold.dev/) (distroless 12.6MB):

```
skaffold run -f ci/skaffold.yaml
```

With [buildpacks](https://buildpacks.io/) ([OCI](https://opencontainers.org/) 30.7MB):

```
pushd broadcaster && pack build broadcaster --builder paketobuildpacks/builder:tiny && popd
pushd listener && pack build listener --builder paketobuildpacks/builder:tiny && popd
pushd decoder && pack build decoder --builder paketobuildpacks/builder:tiny && popd
pushd writer && pack build writer --builder paketobuildpacks/builder:tiny && popd
pushd reporter && pack build reporter --builder paketobuildpacks/builder:tiny && popd
docker run broadcaster
docker run listener
docker run decoder
docker run writer
docker run reporter
```

### Local Development

```
go get github.com/azer/yolo
make help
```
