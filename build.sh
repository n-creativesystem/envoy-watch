#!/usr/bin/env sh

build() {
    echo $1
    echo $2
    GOOS=$1 \
    GOARCH=$2 \
    CGO_ENABLED=0 \
    go build -a -tags netgo -installsuffix netgo \
    -ldflags="-s -w -X \"github.com/n-creativesystem/envoy-watch/version.Version=${VERSION}\" -X \"github.com/n-creativesystem/envoy-watch/version.Revision=${REVISION}\" -extldflags \"-static\"" \
    -o bin/$arch/$os/${NAME}
}

for os in darwin linux; do
    for arch in amd64 386; do
        build $os $arch
    done
done
